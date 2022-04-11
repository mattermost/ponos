package app

import (
	"fmt"

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

func (creq CallRequest) StoreAppConfig(cfg *AppConfig) error {
	asBot := creq.AsBot()

	status, err := asBot.KVSet("config", "", cfg)
	if err != nil {
		return errors.Wrap(err, "failed to save the app config")
	}
	fmt.Print(status)
	return nil
}

func (creq CallRequest) GetAppConfig() (*AppConfig, error) {
	asBot := creq.AsBot()

	var cfg *AppConfig
	err := asBot.KVGet("config", "", &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the app config")
	}
	return cfg, nil
}

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
