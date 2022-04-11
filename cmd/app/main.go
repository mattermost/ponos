package main

import (
	"net"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/utils"
	root "github.com/mattermost/ponos"
	"github.com/mattermost/ponos/config"
	"github.com/mattermost/ponos/function"
	"github.com/oklog/oklog/pkg/group"
	"go.uber.org/zap/zapcore"
)

var logger utils.Logger

func init() {
	logger = utils.MustMakeCommandLogger(zapcore.DebugLevel)
}

func main() {
	cfg, err := config.Load(logger)
	if err != nil {
		logger.WithError(err).Errorf("failed to load config")
		os.Exit(1)
	}
	if cfg.App.Type == apps.DeployHTTP {
		root.AppManifest.Deploy.HTTP.RootURL = "http://192.168.68.104:3000"
	}

	r := mux.NewRouter()
	function.Init(string(cfg.App.Type), r, logger)
	http.Handle("/", r)

	startApp(cfg, r)
}

func startApp(cfg config.Options, r *mux.Router) {
	// http
	if cfg.App.Type == apps.DeployHTTP {
		httpListener, err := net.Listen("tcp", cfg.ListenAddress)
		if err != nil {
			logger.WithError(err).Errorf("failed to listen in %s", cfg.ListenAddress)
			os.Exit(1)
		}
		var g group.Group
		g.Add(func() error {
			logger.With("listen", cfg.ListenAddress).Infof("Server started")
			return http.Serve(httpListener, r)
		}, func(error) {
			httpListener.Close()
		})

		logger.WithError(g.Run()).Errorf("exit")
		return
	}

	// lambda
	lambda.Start(httpadapter.New(r).Proxy)
}
