package moderated_requests

import (
	"time"

	"gorm.io/datatypes"
)

type ModeratedRequest struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Kind      string         `json:"kind"`
	State     string         `json:"state"`
	Data      datatypes.JSON `json:"data"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
