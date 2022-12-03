package moderated_requests

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Service struct {
	logger log.FieldLogger
	db     *gorm.DB
}

// NewService creates a service to manage moderated requests
func NewService(logger log.FieldLogger, db *gorm.DB) *Service {
	return &Service{
		logger: logger,
		db:     db,
	}
}

func (s *Service) CreateRequest(req ModeratedRequestData) error {
	data, err := json.Marshal(req.Data)

	if err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return err
	}

	result := s.db.Create(&ModeratedRequest{
		Kind:  req.Kind,
		State: Pending,
		Data:  datatypes.JSON(data),
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *Service) GetRequests() ([]ModeratedRequest, error) {
	var requests []ModeratedRequest
	result := s.db.Find(&requests)

	return requests, result.Error
}
