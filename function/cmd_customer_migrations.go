package function

import (
	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/utils"
	"github.com/mattermost/ponos/app"
	"github.com/mattermost/ponos/internal/api"
	"github.com/mattermost/ponos/migrations"
	"github.com/pkg/errors"
)

func (a *App) createMigrationHandler(creq app.CallRequest) (string, error) {
	cfg, err := creq.GetAppConfig()
	if err != nil {
		return "", errors.Wrap(err, "failed to get config")
	}

	bucket := creq.GetValue("bucket", "")
	bucketFolder := creq.GetValue("bucket_folder", "")
	destroy, _ := creq.BoolValue("destroy")
	client := api.NewPonosClient("as", cfg.PonosURL)
	// destroy resources
	if destroy {
		if err := client.DeleteMiration(migrations.CustomerMigrationRequest{
			Bucket:       bucket,
			BucketFolder: bucketFolder,
			Apply:        true,
		}); err != nil {
			a.Logger.WithError(err).Errorf("failed to destroy migration for: %s", bucketFolder)
			return "", err
		}
		return "migration resources destroyed succesfully", nil
	}

	// create resources
	if err := client.CreateMiration(migrations.CustomerMigrationRequest{
		Bucket:       bucket,
		BucketFolder: bucketFolder,
		Apply:        true,
	}); err != nil {
		a.Logger.WithError(err).Errorf("failed to create migration for: %s", bucketFolder)
		return "", err
	}
	return "migration resources created succesfully", nil
}

func (a *App) createMigrationFormHandler(creq app.CallRequest) (apps.Form, error) {
	bucket := creq.GetValue("bucket", "")
	if bucket == "" {
		return apps.Form{}, utils.NewInvalidError("field bucket is required ")
	}
	bucketFolder := creq.GetValue("bucket_folder", "")
	if bucketFolder == "" {
		return apps.Form{}, utils.NewInvalidError("field bucket_folder is required ")
	}
	return *a.deleteWorkspaceForm(creq), nil
}

func (a *App) createMigrationCommandBinding(creq app.CallRequest) apps.Binding {
	b := apps.Binding{
		Label:       "migration",
		Description: "Create Cloud migration resources",
		Location:    apps.Location("migration"),
		Icon:        a.App.Icon,
	}
	b.Bindings = append(b.Bindings, apps.Binding{
		Label:    "manage",
		Location: apps.Location("migration-create"),
		Icon:     a.App.Icon,
		Form:     a.createMigrationForm(creq),
	})
	return b
}

func (a *App) createMigrationForm(creq app.CallRequest) *apps.Form {
	return &apps.Form{
		Title: "Create migration resources",
		Icon:  a.Icon,
		Fields: []apps.Field{
			{
				Name:       "bucket",
				ModalLabel: "The existing bucket we store migrations",
				IsRequired: true,
			},
			{
				Name:       "bucket_folder",
				ModalLabel: "Customer's bucket folder",
				Type:       apps.FieldTypeText,
				IsRequired: true,
			},
			{
				Name:       "destroy",
				ModalLabel: "Destroy the resources",
				Type:       apps.FieldTypeBool,
			},
		},
		Submit: apps.NewCall("/migrations/manage").
			WithState(map[string]string{
				"bucket":        creq.GetValue("bucket", ""),
				"bucket_folder": creq.GetValue("bucket_folder", ""),
				"apply":         creq.GetValue("apply", ""),
				"destroy":       creq.GetValue("destroy", ""),
			}).
			WithExpand(apps.Expand{
				ActingUserAccessToken: apps.ExpandAll,
				ActingUser:            apps.ExpandSummary,
			}),
	}
}
