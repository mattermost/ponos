package config

import (
	"strings"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// App config for Mattermost app
type App struct {
	Type    apps.DeployType
	RootURL string `mapstructure:"root_url"`
	Secret  string
}

// Options config to set to run the app.
type Options struct {
	Debug         bool
	App           App
	ListenAddress string `mapstructure:"address"`
	IsLocal       bool   `mapstructure:"local"`
	Environment   string
}

func (o *Options) Validate() error {
	return nil
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("ponos")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	defaults := map[string]interface{}{
		"debug":       false,
		"environment": "dev",
		"address":     ":3000",

		// application settings if http or lambda
		"app.type":     apps.DeployHTTP,
		"app.root_url": "http://localhost:3000",
		"app.secret":   "secretkey",
	}

	for key, value := range defaults {
		viper.SetDefault(key, value)
	}
}

// Load will load the necessary config
func Load(logger utils.Logger) (Options, error) {
	if err := viper.ReadInConfig(); err != nil {
		logger.Warnf(errors.Wrap(err, "unable to find config.yml. loading config from environment variables").Error())
	}
	var cfg Options
	if err := viper.Unmarshal(&cfg); err != nil {
		return Options{}, errors.Wrap(err, "failed to load")
	}
	if err := cfg.Validate(); err != nil {
		return Options{}, errors.Wrap(err, "failed to validate the config")
	}

	return cfg, nil
}
