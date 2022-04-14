package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"emperror.dev/errors"
	"github.com/mattermost/ponos/migrations"
	"github.com/mattermost/ponos/workspaces"
)

type PonosClient struct {
	client *http.Client
	token  string
	url    string
}

// NewPonosClient factory method to create a Ponos Client
func NewPonosClient(token, url string) *PonosClient {
	return &PonosClient{
		client: http.DefaultClient,
		token:  token,
		url:    url,
	}
}

// DeleteWorkspace will make a request to delete a workspace
func (c *PonosClient) DeleteWorkspace(data workspaces.DeleteWorkspaceRequest) error {
	body, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to parse delete workspaces request payload")
	}
	endpoint := fmt.Sprintf("%s/workspaces", c.url)
	req, err := http.NewRequest("DELETE", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrap(err, "failed to create a request to delete a workspace")
	}
	return c.do(req)
}

// CreateMiration will create the necessary resources for
// migrating a customer in cloud with AWAT. Resources:
// - S3 Bucket
// - IAM resouces & policy
func (c *PonosClient) CreateMiration(data migrations.CustomerMigrationRequest) error {
	body, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to parse create migration request payload")
	}
	endpoint := fmt.Sprintf("%s/migrations", c.url)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrap(err, "failed to create a request to delete a workspace")
	}
	return c.do(req)
}

// DeleteMiration will delete the created resources for
// migrating a customer in cloud with AWAT. Resources:
// - S3 Bucket
// - IAM resouces & policy
func (c *PonosClient) DeleteMiration(data migrations.CustomerMigrationRequest) error {
	body, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to parse create migration request payload")
	}
	endpoint := fmt.Sprintf("%s/migrations", c.url)
	req, err := http.NewRequest("DELETE", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrap(err, "failed to create a request to delete a workspace")
	}
	return c.do(req)
}

func (c *PonosClient) do(req *http.Request) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to delete a workspace")
	}
	if resp.StatusCode > 400 && resp.StatusCode < 500 {
		return errors.New("failed because of a bad request")
	}
	return nil
}
