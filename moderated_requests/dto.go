package moderated_requests

type Kind = string

const (
	DeleteWorkspace = "DELETE_WORKSPACE"
)

type State = string

const (
	Pending  = "PENDING"
	Rejected = "REJECTED"
	Approved = "APPROVED"
)

type Payload = map[string]interface{}

type ModeratedRequestData struct {
	Kind Kind    `json:"kind"`
	Data Payload `json:"data"`
}
