package api

import "github.com/gorilla/mux"

// Register registers the API endpoints on the given router.
func Create(rootRouter *mux.Router, context *Context) *mux.Router {
	// api handler at /api
	router := rootRouter.PathPrefix("/api").Subrouter()

	AddMigrationsAPI(router, context)
	AddWorkspacesAPI(router, context)
	AddModeratedRequestsAPI(router, context)

	return router
}
