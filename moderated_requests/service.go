package moderated_requests

import (
	"encoding/json"

	"github.com/pkg/errors"
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
		return errors.Wrap(err, "Failed to read moderated request")
	}

	if err := req.Validate(); err != nil {
		return errors.Wrap(err, "Validation failed")
	}

	result := s.db.Create(&ModeratedRequest{
		Kind:  req.Kind,
		State: Pending,
		Data:  datatypes.JSON(data),
	})

	if result.Error != nil {
		return errors.Wrap(result.Error, "Failed to persist moderated request")
	}

	return nil
}

func (s *Service) GetRequests() ([]ModeratedRequest, error) {
	var requests []ModeratedRequest
	result := s.db.Find(&requests)

	return requests, errors.Wrap(result.Error, "Failed to retrieve moderated requests")
}
