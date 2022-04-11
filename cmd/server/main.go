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
	"github.com/mattermost/ponos/workspaces"
)

func main() {
	provisionerClient := cmodel.NewClient(os.Getenv("PONOS_PROVISIONER_ADDRESS"))
	workspaceSvc := workspaces.NewService(provisionerClient, logger)
	migrationsSvc := migrations.NewService(logger)

	router := mux.NewRouter()

	api.Create(router, &api.Context{
		Logger:            logger,
		WorkspaceService:  workspaceSvc,
		MigrationsService: migrationsSvc,
	})

	srv := &http.Server{
		Addr:           ":3001",
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
	srv.Shutdown(ctx)
}
