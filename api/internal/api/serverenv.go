package api

import (
	"api/pkg/logging"
	"api/pkg/oauth"
)

type ServerEnv struct {
	Logger        *logging.Logger
	Config        *ServerConfig
	OAuthVerifier oauth.TokenVerifier
}
