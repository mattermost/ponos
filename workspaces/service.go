package workspaces

import (
	cmodel "github.com/mattermost/mattermost-cloud/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ProvisionerRequester the interface which describes Provisioner API client
type ProvisionerRequester interface {
	GetInstallationByDNS(DNS string, request *cmodel.GetInstallationRequest) (*cmodel.InstallationDTO, error)
	DeleteInstallation(id string) error
}

type Service struct {
	provisionerClient ProvisionerRequester
	workspaceClient   *WorkspaceClient
	logger            log.FieldLogger
}

// NewService creates a service to make provisioner requests
func NewService(provisionerClient ProvisionerRequester, workspaceClient *WorkspaceClient, logger log.FieldLogger) *Service {
	return &Service{
		provisionerClient: provisionerClient,
		workspaceClient:   workspaceClient,
		logger:            logger,
	}
}

// DeleteWorkspace deletes a workspace with the DNS name provided
func (s *Service) DeleteWorkspace(dnsName string) error {
	installation, err := s.provisionerClient.GetInstallationByDNS(dnsName, &cmodel.GetInstallationRequest{})
	if err != nil {
		return errors.Wrap(err, "failed to get installation by DNS")
	}
	if installation == nil {
		return errors.New("failed to found an installation with provided DNS name")
	}
	return s.workspaceClient.DeleteWorkspace(installation.ID)
}
