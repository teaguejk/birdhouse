package api

import (
	"api/pkg/logging"
)

type ServerEnv struct {
	Logger *logging.Logger
	Config *ServerConfig
}
