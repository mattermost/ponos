package function

import (
	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/ponos/app"
	"github.com/pkg/errors"
)

var (
	fieldPonosURL string = "ponos_url"
	fieldApiToken string = "api_token"
)

var configureCommand = app.Command{
	Name:        "configure",
	Hint:        "[ --ponos-url --api-token]",
	Description: "Configure Ponos service URL and API token for Ponos Service",

	BaseForm: &apps.Form{
		Submit: &apps.Call{
			Expand: &apps.Expand{
				ActingUserAccessToken: apps.ExpandAll,
				ActingUser:            apps.ExpandSummary,
			},
		},
		Title: "Configure Ponos App",
		Fields: []apps.Field{
			{
				Name:        fieldPonosURL,
				ModalLabel:  "Ponos URL",
				Type:        apps.FieldTypeText,
				Description: "The root URL of your Ponos service.",
				IsRequired:  true,
			},
			{
				Name:        fieldApiToken,
				ModalLabel:  "Ponos API Token",
				Type:        apps.FieldTypeText,
				Description: "The API Token of your Ponos service.",
				IsRequired:  true,
			},
		},
	},

	Handler: app.RequireAdmin(func(creq app.CallRequest) apps.CallResponse {
		apiToken := creq.GetValue(fieldApiToken, "")
		ponosURL := creq.GetValue(fieldPonosURL, "")

		err := creq.StoreAppConfig(&app.AppConfig{
			Token:    apiToken,
			PonosURL: ponosURL,
		})
		if err != nil {
			return apps.NewErrorResponse(errors.Wrap(err, "failed to store App Config to Mattermost"))
		}
		return apps.NewTextResponse(
			"updated app configuration:\n"+
				"  - Ponos URL: `%s`\n",
			ponosURL,
		)

	}),
}
