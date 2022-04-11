package workspaces

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

// DeleteWorkspaceRequest the data transfer object to delete a
// a workspace with the provided DNS name
type DeleteWorkspaceRequest struct {
	Name string `json:"dns_name"`
}

// Validate validates the values of a delete workspace request
func (r *DeleteWorkspaceRequest) Validate() error {
	if r.Name == "" {
		return errors.New("DNS name cannot be emp")
	}
	return nil
}

// NewCustomerMigrationRequestFromReader decodes the request and returns after validation and setting the defaults.
func NewDeleteWorkspaceRequest(reader io.Reader) (*DeleteWorkspaceRequest, error) {
	var deleteWorkspaceRequest DeleteWorkspaceRequest
	if err := json.NewDecoder(reader).Decode(&deleteWorkspaceRequest); err != nil && err != io.EOF {
		return nil, errors.Wrap(err, "failed to decode provision customer migration request")
	}

	if err := deleteWorkspaceRequest.Validate(); err != nil {
		return nil, errors.Wrap(err, "customer migration request failed validation")
	}

	return &deleteWorkspaceRequest, nil
}
