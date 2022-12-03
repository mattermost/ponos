package api

import (
	"github.com/mattermost/ponos/migrations"
	"github.com/mattermost/ponos/moderated_requests"
	"github.com/mattermost/ponos/workspaces"
	log "github.com/sirupsen/logrus"
)

// Context provides the API with all necessary data and interfaces for responding to requests.
//
// It is cloned before each request, allowing per-request changes such as logger annotations.
type Context struct {
	RequestID                string
	WorkspaceService         *workspaces.Service
	MigrationsService        *migrations.Service
	ModeratedRequestsService *moderated_requests.Service
	Logger                   log.FieldLogger
}

// Clone creates a shallow copy of context, allowing clones to apply per-request changes.
func (c *Context) Clone() *Context {
	return &Context{
		MigrationsService:        c.MigrationsService,
		WorkspaceService:         c.WorkspaceService,
		ModeratedRequestsService: c.ModeratedRequestsService,
		Logger:                   c.Logger,
	}
}
