package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	cmodel "github.com/mattermost/mattermost-cloud/model"
	"github.com/mattermost/ponos/internal/api"
	"github.com/mattermost/ponos/migrations"
	"github.com/mattermost/ponos/moderated_requests"
	"github.com/mattermost/ponos/workspaces"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	provisionerURL := os.Getenv("PONOS_PROVISIONER_ADDRESS")
	if provisionerURL == "" {
		logger.Error("failed to start, PONOS_PROVISIONER_ADDRESS env variable is required")
		os.Exit(1)
		return
	}
	workspacesURL := os.Getenv("PONOS_WORKSPACES_ADDRESS")
	if workspacesURL == "" {
		logger.Error("failed to start, PONOS_WORKSPACES_ADDRESS env variable is required")
		os.Exit(1)
		return
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})

	// TODO (pantelis.vratsalis): Good for now, but in the long run we want a db migration tool
	db.Table("moderated_requests").AutoMigrate(&moderated_requests.ModeratedRequest{})

	if err != nil {
		logger.Error("failed to start, could not connect to the database")
		os.Exit(1)
		return
	}

	provisionerClient := cmodel.NewClient(provisionerURL)
	workspaceClient := workspaces.NewHTTPClient(workspacesURL, http.DefaultClient)
	workspaceSvc := workspaces.NewService(provisionerClient, workspaceClient, logger)
	migrationsSvc := migrations.NewService(logger)
	moderatedRequestsSvc := moderated_requests.NewService(logger, db)
	router := mux.NewRouter()

	api.Create(router, &api.Context{
		Db:                       db,
		Logger:                   logger,
		WorkspaceService:         workspaceSvc,
		MigrationsService:        migrationsSvc,
		ModeratedRequestsService: moderatedRequestsSvc,
	})

	srv := &http.Server{
		Addr:           ":3000",
		Handler:        router,
		ReadTimeout:    180 * time.Second,
		WriteTimeout:   180 * time.Second,
		IdleTimeout:    time.Second * 180,
		MaxHeaderBytes: 1 << 20,
		ErrorLog:       log.New(&logrusWriter{logger}, "", 0),
	}

	go func() {
		logger.WithField("addr", srv.Addr).Info("API server listening")
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Error("Failed to listen and serve")
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via:
	//  - SIGINT (Ctrl+C)
	//  - SIGTERM (Ctrl+/) (Kubernetes pod rolling termination)
	// SIGKILL and SIGQUIT will not be caught.
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a valid signal.
	sig := <-c
	logger.WithField("shutdown-signal", sig.String()).Info("Shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	srv.Shutdown(ctx) // nolint
}
