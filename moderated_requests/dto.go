package moderated_requests

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type Kind = string

const (
	DeleteWorkspace Kind = "WORKSPACES_DELETE"
)

type State = string

const (
	Pending  State = "PENDING"
	Rejected State = "REJECTED"
	Approved State = "APPROVED"
)

type Payload = map[string]interface{}

type ModeratedRequestData struct {
	Kind Kind    `json:"kind" validate:"kind"`
	Data Payload `json:"data"`
}

func (s *ModeratedRequestData) Validate() error {
	validate := validator.New()
	err := validate.RegisterValidation("kind", kind)

	if err != nil {
		return errors.Wrap(err, "Failed to register custom validator")
	}

	return validate.Struct(s)
}

func kind(f validator.FieldLevel) bool {
	validKinds := []string{
		DeleteWorkspace,
	}

	for _, v := range validKinds {
		if f.Field().String() == v {
			return true
		}
	}

	return false
}
