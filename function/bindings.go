package function

import (
	"github.com/mattermost/mattermost-plugin-apps/apps"
	root "github.com/mattermost/ponos"
	"github.com/mattermost/ponos/app"
)

func (a *App) getBindings(creq app.CallRequest) apps.CallResponse {
	bindings := app.AppendBinding(nil, &apps.Binding{
		Location: apps.LocationCommand,
		Bindings: []apps.Binding{
			{
				Label:       "ponos",
				Description: "Automate toil work operations for Cloud",
				Icon:        root.AppManifest.Icon,

				Bindings: append(
					a.commandBindings(creq),
					a.debugCommandBindings(creq)...,
				),
			},
		},
	})
	return apps.NewDataResponse(bindings)
}

func (a *App) commandBindings(creq app.CallRequest) []apps.Binding {
	var bindings []apps.Binding

	// admin commands
	if creq.Context.ActingUser != nil && creq.Context.ActingUser.IsSystemAdmin() {
		bindings = append(bindings,
			configureCommand.Binding(creq),
			addUserCommand.Binding(creq),
			a.getInfoCommand().Binding(creq),
		)
	}

	cfg, err := creq.GetAppConfig()
	if err != nil {
		return bindings
	}
	if cfg == nil {
		return bindings
	}
	ua, err := creq.GetUserAccess()
	if err != nil {
		return bindings
	}
	if len(ua.UserIDS) == 0 || !ua.IsAllowed(creq.Context.ActingUser.Id) {
		return bindings
	}

	bindings = append(bindings, a.createMigrationCommandBinding(creq))
	bindings = append(bindings, a.deleteWorkspaceCommandBinding(creq))
	return bindings
}

func (a *App) debugCommandBindings(creq app.CallRequest) []apps.Binding {
	if !creq.Context.DeveloperMode &&
		(creq.Context.ActingUser == nil || !creq.Context.ActingUser.IsSystemAdmin()) {
		return nil
	}

	return []apps.Binding{}
}
