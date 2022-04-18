package app

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/utils"
	"github.com/mattermost/mattermost-plugin-apps/utils/httputils"
)

type HandlerFunc func(CallRequest) apps.CallResponse

func (a *App) HandleCall(p string, h HandlerFunc) {
	a.Router.Path(p).HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		creq := CallRequest{
			GoContext: req.Context(),
		}
		err := json.NewDecoder(req.Body).Decode(&creq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		creq.App = *a
		creq.App.Logger = a.Logger

		cresp := h(creq)
		if cresp.Type == apps.CallResponseTypeError {
			creq.App.Logger.WithError(cresp).Debugw("Call failed.")
		}
		_ = httputils.WriteJSON(w, cresp)
	})
}

func (a *App) HandleCommand(command Command) {
	a.HandleCall(command.Path(), command.Handler)
}

func FormHandler(h func(CallRequest) (apps.Form, error)) HandlerFunc {
	return func(creq CallRequest) apps.CallResponse {
		f, err := h(creq)
		if err != nil {
			creq.App.Logger.WithError(err).Infow("failed to respond with form")
			return apps.NewErrorResponse(err)
		}
		return apps.NewFormResponse(f)
	}
}

func LookupHandler(h func(CallRequest) []apps.SelectOption) HandlerFunc {
	return func(creq CallRequest) apps.CallResponse {
		opts := h(creq)
		return apps.NewLookupResponse(opts)
	}
}

func CallHandler(h func(CallRequest) (string, error)) HandlerFunc {
	return func(creq CallRequest) apps.CallResponse {
		text, err := h(creq)
		if err != nil {
			creq.App.Logger.WithError(err).Infow("failed to process call")
			return apps.NewErrorResponse(err)
		}
		return apps.NewTextResponse(text)
	}
}

func RequireAdmin(h HandlerFunc) HandlerFunc {
	return func(creq CallRequest) apps.CallResponse {
		if creq.Context.ActingUser != nil && !creq.Context.ActingUser.IsSystemAdmin() {
			return apps.NewErrorResponse(
				utils.NewUnauthorizedError("system administrator role is required to invoke " + creq.Path))
		}
		return h(creq)
	}
}

func RequireConnectedUsers(h HandlerFunc) HandlerFunc {
	return func(creq CallRequest) apps.CallResponse {
		cfg, err := creq.GetAppConfig()
		if err != nil {
			return apps.NewErrorResponse(
				utils.NewUnauthorizedError("not able to fetch app config"))
		}

		if cfg == nil {
			return apps.NewErrorResponse(
				utils.NewUnauthorizedError("you need to run /ponos configure first"))
		}
		if err := cfg.Validate(); err != nil {
			return apps.NewErrorResponse(
				utils.NewUnauthorizedError("you need to run /ponos configure first and fill all the fiedls"))
		}
		return h(creq)
	}
}

func (creq CallRequest) AppProxyURL(paths ...string) string {
	p := path.Join(append([]string{creq.Context.AppPath}, paths...)...)
	return creq.Context.MattermostSiteURL + p
}
