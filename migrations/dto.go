package migrations

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/mattermost/ponos/internal/tools/terraform"
	"github.com/pkg/errors"
)

// CustomerMigrationRequest the data transfer object to create a
// new folder under the migrations S3 buc`ket and AWS account
type CustomerMigrationRequest struct {
	Bucket       string `json:"bucket"`
	BucketFolder string `json:"bucket_folder"`
	Apply        bool   `json:"apply"`
}

// Validate validates the values of a customer migration request
func (r *CustomerMigrationRequest) Validate() error {
	if r.Bucket == "" {
		return errors.New("S3 bucket cannot be empty")
	}
	if r.BucketFolder == "" {
		return errors.New("S3 bucket folder cannot be empty")
	}
	return nil
}

func (d *CustomerMigrationRequest) toTerraformArgs() []string {
	return []string{
		terraform.Arg("var", fmt.Sprintf("region=%s", os.Getenv("AWS_REGION"))),
		terraform.Arg("var", fmt.Sprintf("account_id=%s", os.Getenv("PONOS_ACCOUNT_ID"))),
		terraform.Arg("var", fmt.Sprintf("kms_key=%s", os.Getenv("PONOS_KMS_KEY"))),
		terraform.Arg("var", fmt.Sprintf("bucket=%s", d.Bucket)),
		terraform.Arg("var", fmt.Sprintf("customer_bucket_folder=%s", d.BucketFolder)),
		terraform.Arg("var", fmt.Sprintf("customer_policy_name=%s", d.BucketFolder)),
	}
}

// NewCustomerMigrationRequestFromReader decodes the request and returns after validation and setting the defaults.
func NewCustomerMigrationRequestFromReader(reader io.Reader) (*CustomerMigrationRequest, error) {
	var customerMigrationRequest CustomerMigrationRequest

	if err := json.NewDecoder(reader).Decode(&customerMigrationRequest); err != nil && err != io.EOF {
		return nil, errors.Wrap(err, "failed to decode provision customer migration request")
	}

	if err := customerMigrationRequest.Validate(); err != nil {
		return nil, errors.Wrap(err, "customer migration request failed validation")
	}

	return &customerMigrationRequest, nil
}
