package function

import (
	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/ponos/app"
	"github.com/pkg/errors"
)

var (
	fieldUserAllowed string = "user_allowed"
)

var addUserCommand = app.Command{
	Name:        "allow",
	Hint:        "[ --user]",
	Description: "Allow which users can use the Pono ChatOps app. You can use this command to add more",

	BaseForm: &apps.Form{
		Submit: &apps.Call{
			Expand: &apps.Expand{
				ActingUserAccessToken: apps.ExpandAll,
				ActingUser:            apps.ExpandSummary,
			},
		},
		Title: "Add user access for Ponos App",
		Fields: []apps.Field{
			{
				Name:        fieldUserAllowed,
				ModalLabel:  "User Allowed",
				Type:        apps.FieldTypeUser,
				Description: "The user which is allowed to use the app.",
				IsRequired:  true,
			},
		},
	},

	Handler: app.RequireAdmin(func(creq app.CallRequest) apps.CallResponse {
		userAllowed := creq.GetValue(fieldUserAllowed, "")
		err := creq.StoreUserAccess(userAllowed)
		if err != nil {
			return apps.NewErrorResponse(errors.Wrap(err, "failed to store user access config to Mattermost"))
		}
		return apps.NewTextResponse("updated user access configuration")

	}),
}
