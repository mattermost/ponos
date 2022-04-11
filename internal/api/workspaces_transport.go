package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/ponos/workspaces"
)

// AddWorkspacesAPI registers workspaces endpoints on the given router.
func AddWorkspacesAPI(apiRouter *mux.Router, context *Context) {
	addContext := func(handler ContextHandlerFunc) *ContextHandler {
		return NewContextHandler(context, handler)
	}

	migrationsRouter := apiRouter.PathPrefix("/workspaces").Subrouter()
	migrationsRouter.Handle("", addContext(handleDeleteWorkspace)).Methods("DELETE")
}

// handleCreateCustomerMigration responds to Create /api/migrations
func handleDeleteWorkspace(c *Context, w http.ResponseWriter, r *http.Request) {
	request, err := workspaces.NewDeleteWorkspaceRequest(r.Body)
	if err != nil {
		c.Logger.WithError(err).Error("failed to deserialize delete workspace request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := c.WorkspaceService.DeleteWorkspace(request.Name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
