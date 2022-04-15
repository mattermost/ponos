package api

import (
	"net/http"

	"github.com/mattermost/mattermost-cloud/model"
	log "github.com/sirupsen/logrus"
)

// ContextHandlerFunc to add context
type ContextHandlerFunc func(c *Context, w http.ResponseWriter, r *http.Request)

// ContextHandler for HTTP
type ContextHandler struct {
	context *Context
	handler ContextHandlerFunc
}

// ServeHTTP for server to run HTTP
func (h ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := h.context.Clone()
	context.RequestID = model.NewID()
	context.Logger = context.Logger.WithFields(log.Fields{
		"request": context.RequestID,
	})

	h.handler(context, w, r)
}

// NewContextHandler to create a handler for context
func NewContextHandler(context *Context, handler ContextHandlerFunc) *ContextHandler {
	return &ContextHandler{
		context: context,
		handler: handler,
	}
}
