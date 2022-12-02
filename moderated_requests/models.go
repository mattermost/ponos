package moderated_requests

import (
	"time"

	"gorm.io/datatypes"
)

type ModeratedRequest struct {
	ID        uint `gorm:"primaryKey"`
	Kind      string
	State     string
	Data      datatypes.JSON
	CreatedAt time.Time
	UpdatedAt time.Time
}
