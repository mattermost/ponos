package workspaces

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type WorkspaceClient struct {
	URL    string
	client *http.Client
}

// NewHTTPClient creates a workspace HTTP client for managing workspace
func NewHTTPClient(URL string, client *http.Client) *WorkspaceClient {
	return &WorkspaceClient{
		URL:    URL,
		client: client,
	}
}

// DeleteWorkspace makes an http request to delete a workspace
func (c *WorkspaceClient) DeleteWorkspace(installationID string) error {
	if installationID == "" {
		return errors.New("failed to delete workspace, installationID is required")
	}
	endpoint := fmt.Sprintf("%s/installation/%s", c.URL, installationID)
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create a delete HTTP request")
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to delete a workspace")
	}
	if resp.StatusCode >= 400 && resp.StatusCode <= 500 {
		return errors.New("failed to delete because of a bad request")
	}
	return nil
}
