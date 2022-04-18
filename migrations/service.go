package migrations

import (
	"fmt"
	"os"

	"emperror.dev/errors"
	"github.com/mattermost/ponos/internal/tools/terraform"
	log "github.com/sirupsen/logrus"
)

var templateDirMigrations = "terraform/aws/customer-migrations"

type Service struct {
	logger log.FieldLogger
}

// NewService creates a service to make terraform
func NewService(logger log.FieldLogger) *Service {
	return &Service{
		logger: logger,
	}
}

// CreateMigrationInfra creates the necessary AWS IAM, bucket folder
// so the migration resources are ready to be used.
func (s *Service) CreateMigrationInfra(dto *CustomerMigrationRequest) error {
	tf, err := terraform.New(templateDirMigrations, os.Getenv("PONOS_TERRAFORM_STATE_STORE"), s.logger)
	if err != nil {
		return errors.Wrap(err, "failed to create Terraform commander")
	}

	s.logger.Info("starting migration infra")
	stateObject := fmt.Sprintf("customer-migration-%s", dto.BucketFolder)
	if err = tf.Init(stateObject); err != nil {
		return errors.Wrap(err, "failed to run Terraform init")
	}
	if err := tf.Plan(dto.toTerraformArgs()...); err != nil {
		return errors.Wrap(err, "failed to run Terraform plan")
	}
	s.logger.Info("successfully ran Terraform plan")

	if dto.Apply {
		s.logger.Info("applying migration infra")
		if err := tf.Apply(dto.toTerraformArgs()...); err != nil {
			return errors.Wrap(err, "failed to apply migration infra")
		}
		s.logger.Info("successfully ran Terraform apply")
	}
	return nil
}

// DeleteMigrationInfra creates the necessary AWS IAM, bucket folder
// so the migration resources are total removed from our infrastructure.
func (s *Service) DeleteMigrationInfra(dto *CustomerMigrationRequest) error {
	tf, err := terraform.New(templateDirMigrations, os.Getenv("PONOS_TERRAFORM_STATE_STORE"), s.logger)
	if err != nil {
		return errors.Wrap(err, "failed to create Terraform commander")
	}

	s.logger.Info("destroying migration infra")
	stateObject := fmt.Sprintf("customer-migration-%s", dto.BucketFolder)
	if err = tf.Init(stateObject); err != nil {
		return errors.Wrap(err, "failed to run Terraform init")
	}
	if err := tf.Plan(dto.toTerraformArgs()...); err != nil {
		return errors.Wrap(err, "failed to run Terraform plan")
	}
	s.logger.Info("successfully ran Terraform plan")

	if dto.Apply {
		s.logger.Info("destroying migration infra")
		if err := tf.Destroy(dto.toTerraformArgs()...); err != nil {
			return errors.Wrap(err, "failed to destroy migration infra")
		}
		s.logger.Info("successfully ran Terraform destroy")
	}
	return nil
}
