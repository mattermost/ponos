package function

import (
	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/utils"
	"github.com/mattermost/ponos/app"
	"github.com/mattermost/ponos/internal/api"
	"github.com/mattermost/ponos/workspaces"
	"github.com/pkg/errors"
)

func (a *App) deleteWorkspaceHandler(creq app.CallRequest) (string, error) {
	dns := creq.GetValue("name", "")
	cfg, err := creq.GetAppConfig()
	if err != nil {
		return "", errors.Wrap(err, "failed to get config")
	}
	client := api.NewPonosClient(cfg.Token, cfg.PonosURL)

	if err := client.DeleteWorkspace(workspaces.DeleteWorkspaceRequest{
		Name: dns,
	}); err != nil {
		a.Logger.WithError(err).Errorf("failed to delete workspace: %s", dns)
		return "", err
	}
	return "deleted succesfully", nil
}

func (a *App) deleteWorkspaceFormHandler(creq app.CallRequest) (apps.Form, error) {
	dns := creq.GetValue("name", "")
	if dns == "" {
		return apps.Form{}, utils.NewInvalidError("field dns is required ")
	}
	return *a.deleteWorkspaceForm(creq), nil
}

func (a *App) deleteWorkspaceCommandBinding(creq app.CallRequest) apps.Binding {
	b := apps.Binding{
		Label:       "workspace",
		Description: "Delete Cloud Workspace",
		Location:    apps.Location("workspace"),
		Icon:        a.App.Icon,
	}
	b.Bindings = append(b.Bindings, apps.Binding{
		Label:    "delete",
		Location: apps.Location("workspace-delete"),
		Icon:     a.App.Icon,
		Form:     a.deleteWorkspaceForm(creq),
	})
	return b
}

func (a *App) deleteWorkspaceForm(creq app.CallRequest) *apps.Form {
	return &apps.Form{
		Title: "Delete Cloud Workspace",
		Icon:  a.Icon,
		Fields: []apps.Field{
			{
				Type:       "text",
				Name:       "name",
				Label:      "name",
				IsRequired: true,
			},
		},
		Submit: apps.NewCall("/workspaces/delete").
			WithState(map[string]string{
				"name": creq.GetValue("name", ""),
			}).
			WithExpand(apps.Expand{
				ActingUserAccessToken: apps.ExpandAll,
				ActingUser:            apps.ExpandSummary,
			}),
	}
}
