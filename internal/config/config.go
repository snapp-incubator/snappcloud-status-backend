package config

import (
	"github.com/snapp-incubator/snappcloud-status-backend/internal/api/http"
	"github.com/snapp-incubator/snappcloud-status-backend/pkg/logger"
)

type Config struct {
	Logger *logger.Config `koanf:"logger"`
	HTTP   *http.Config   `koanf:"http"`
}
