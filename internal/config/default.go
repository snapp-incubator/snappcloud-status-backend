package config

import (
	"github.com/snapp-incubator/snappcloud-status-backend/internal/api/http"
	"github.com/snapp-incubator/snappcloud-status-backend/pkg/logger"
)

func Default() *Config {
	return &Config{
		Logger: &logger.Config{
			Development: true,
			Level:       "debug",
			Encoding:    "console",
		},
		HTTP: &http.Config{
			ListenPort: 8080,
		},
	}
}
