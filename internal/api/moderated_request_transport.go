package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/ponos/moderated_requests"
)

// AddModeratedRequestsAPI registers moderated requests endpoints on the given router.
func AddModeratedRequestsAPI(apiRouter *mux.Router, context *Context) {
	addContext := func(handler ContextHandlerFunc) *ContextHandler {
		return NewContextHandler(context, handler)
	}

	moderatedRequestsRouter := apiRouter.PathPrefix("/moderated-requests").Subrouter()
	moderatedRequestsRouter.Handle("", addContext(createRequest)).Methods("POST")
}

// POST /api/moderated-requests
// Create a request for a change in "pending" state
func createRequest(c *Context, w http.ResponseWriter, r *http.Request) {
	var body moderated_requests.ModeratedRequestData

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request structure."})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c.ModeratedRequestsService.CreateRequest(body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
