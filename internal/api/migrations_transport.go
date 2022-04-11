package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/ponos/migrations"
)

// AddMigrationsAPI registers migrations endpoints on the given router.
func AddMigrationsAPI(apiRouter *mux.Router, context *Context) {
	addContext := func(handler ContextHandlerFunc) *ContextHandler {
		return NewContextHandler(context, handler)
	}

	migrationsRouter := apiRouter.PathPrefix("/migrations").Subrouter()
	migrationsRouter.Handle("", addContext(handleCreateCustomerMigration)).Methods("POST")
	migrationsRouter.Handle("", addContext(handleDeleteCustomerMigration)).Methods("DELETE")
}

// handleCreateCustomerMigration responds to Create /api/migrations
func handleCreateCustomerMigration(c *Context, w http.ResponseWriter, r *http.Request) {
	request, err := migrations.NewCustomerMigrationRequestFromReader(r.Body)
	if err != nil {
		c.Logger.WithError(err).Error("failed to deserialize customer migration request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := c.MigrationsService.CreateMigrationInfra(request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

// handleDeleteCustomerMigration responds to DELETE /api/migrations
func handleDeleteCustomerMigration(c *Context, w http.ResponseWriter, r *http.Request) {
	request, err := migrations.NewCustomerMigrationRequestFromReader(r.Body)
	if err != nil {
		c.Logger.WithError(err).Error("failed to deserialize customer migration request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := c.MigrationsService.DeleteMigrationInfra(request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
