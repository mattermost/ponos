package function

import (
	"github.com/gorilla/mux"

	"github.com/mattermost/mattermost-plugin-apps/utils"

	root "github.com/mattermost/ponos"
	goapp "github.com/mattermost/ponos/app"
)

var BuildHash string
var BuildHashShort string
var BuildDate string

type App struct {
	goapp.App
	mode string
}

func Init(mode string, r *mux.Router, log utils.Logger) {
	app := App{
		mode: mode,
		App: *goapp.NewApp(r, log).
			WithManifest(root.AppManifest).
			WithStatic(root.Static).
			WithIcon(root.AppManifest.Icon),
	}

	app.HandleCall("/install", app.getInstall)
	app.HandleCall("/bindings", app.getBindings)

	// Command submit handlers.
	app.HandleCommand(app.getInfoCommand())
	app.HandleCommand(configureCommand)
	app.HandleCommand(addUserCommand)

	// Cloud workspace handlers
	app.HandleCall("/workspaces/delete", goapp.RequireConnectedUsers(goapp.CallHandler(app.deleteWorkspaceHandler)))
	app.HandleCall("/form/delete-workspace", goapp.RequireConnectedUsers(goapp.FormHandler(app.deleteWorkspaceFormHandler)))

	// Cloud migration handlers
	app.HandleCall("/migrations/manage", goapp.RequireConnectedUsers(goapp.CallHandler(app.createMigrationHandler)))
	app.HandleCall("/form/delete-migration", goapp.RequireConnectedUsers(goapp.FormHandler(app.createMigrationFormHandler)))

}
