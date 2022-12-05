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
	moderatedRequestsRouter.Handle("", addContext(getRequests)).Methods("GET")
}

// POST /api/moderated-requests
// Create a request for a change in "pending" state
func createRequest(c *Context, w http.ResponseWriter, r *http.Request) {
	var body moderated_requests.ModeratedRequestData

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		c.Logger.WithError(err).Error("Invalid create moderated request payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createRequestError := c.ModeratedRequestsService.CreateRequest(body)

	if createRequestError != nil {
		c.Logger.WithError(createRequestError).Error("Could not create moderated request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// GET /api/moderated-requests
// Get all moderated requests
func getRequests(c *Context, w http.ResponseWriter, r *http.Request) {
	requests, getRequestsError := c.ModeratedRequestsService.GetRequests()

	if getRequestsError != nil {
		c.Logger.WithError(getRequestsError).Error("Could not retrieve moderated requests")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encodeError := json.NewEncoder(w).Encode(map[string]interface{}{"data": &requests})

	if encodeError != nil {
		c.Logger.WithError(getRequestsError).Error("Could not serialize response")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}
