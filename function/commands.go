package function

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/ponos/app"
)

func (a *App) getInstall(creq app.CallRequest) apps.CallResponse {
	return apps.CallResponse{}
}

func (a *App) getInfoCommand() app.Command {
	return app.Command{
		Name:        "info",
		Description: "Ponos App information",

		BaseSubmit: &apps.Call{
			Expand: &apps.Expand{
				ActingUser: apps.ExpandSummary,
				OAuth2App:  apps.ExpandAll,
				OAuth2User: apps.ExpandAll,
			},
		},

		Handler: func(creq app.CallRequest) apps.CallResponse {
			title := "Ponos"
			if BuildDate != "" {
				title += fmt.Sprintf(" built: %s from [%s](https://github.com/mattermost/ponos/commit/%s)",
					BuildDate, BuildHashShort, BuildHash)
			}
			if a.mode != "" {
				title += ", running as " + a.mode
			}
			title += "\n"

			return apps.NewTextResponse(strings.Join([]string{
				title, "\n",
			}, "\n"))
		},
	}
}
