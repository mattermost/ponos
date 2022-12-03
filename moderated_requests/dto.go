package moderated_requests

type Kind = string

const (
	DeleteWorkspace Kind = "DELETE_WORKSPACE"
)

type State = string

const (
	Pending  State = "PENDING"
	Rejected State = "REJECTED"
	Approved State = "APPROVED"
)

type Payload = map[string]interface{}

type ModeratedRequestData struct {
	Kind Kind    `json:"kind"`
	Data Payload `json:"data"`
}
