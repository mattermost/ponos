package moderated_requests

import "github.com/go-playground/validator/v10"

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
	Kind Kind    `json:"kind" validate:"oneof=WORKSPACES_DELETE"`
	Data Payload `json:"data"`
}

func (s *ModeratedRequestData) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
