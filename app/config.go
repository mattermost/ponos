package app

import (
	"github.com/pkg/errors"
)

// AppConfig the config of the app
type AppConfig struct {
	PonosURL string `json:"ponos_url"`
	Token    string `json:"token"`
}

func (a AppConfig) Validate() error {
	if a.PonosURL == "" {
		return errors.New("Ponos URL is required")
	}
	if a.Token == "" {
		return errors.New("API token is required")
	}
	return nil
}

// StoreAppConfig store the config of the app in KV store
func (creq CallRequest) StoreAppConfig(cfg *AppConfig) error {
	asBot := creq.AsBot()

	_, err := asBot.KVSet("config", "", cfg)
	if err != nil {
		return errors.Wrap(err, "failed to save the app config")
	}
	return nil
}

// GetAppConfig return store the config of the app from KV store
func (creq CallRequest) GetAppConfig() (*AppConfig, error) {
	asBot := creq.AsBot()

	var cfg *AppConfig
	err := asBot.KVGet("config", "", &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the app config")
	}
	return cfg, nil
}

// DeleteAppConfig deletes the config of the app from KV store
func (creq CallRequest) DeleteAppConfig() error {
	asBot := creq.AsBot()
	err := asBot.KVDelete("p", "config")
	if err != nil {
		return errors.Wrap(err, "failed to get the app config")
	}
	return nil
}

// UsersAccess the users can use the app
type UsersAccess struct {
	UserIDS []string `json:"user_ids"`
}

// IsAllowed if user allowed
func (ua *UsersAccess) IsAllowed(userID string) bool {
	for _, id := range ua.UserIDS {
		if id == userID {
			return true
		}
	}
	return false
}

// StoreUserAccess stores the user access config in KV store
func (creq CallRequest) StoreUserAccess(userID string) error {
	ua, err := creq.GetUserAccess()
	if err != nil {
		return errors.Wrap(err, "failed to get the user access config")
	}
	asBot := creq.AsBot()

	ua.UserIDS = append(ua.UserIDS, userID)
	_, err = asBot.KVSet("p", "user_access", ua)
	if err != nil {
		return errors.Wrap(err, "failed to save the user access config")
	}
	return nil
}

// GetUserAccess returns the user access config from KV store
func (creq CallRequest) GetUserAccess() (*UsersAccess, error) {
	asBot := creq.AsBot()

	var ua *UsersAccess
	err := asBot.KVGet("p", "user_access", &ua)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the user access config")
	}
	if ua == nil {
		return &UsersAccess{}, nil
	}
	return ua, nil
}
