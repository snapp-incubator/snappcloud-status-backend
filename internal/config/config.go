package config

import (
	"github.com/snapp-incubator/snappcloud-status-backend/internal/querier"
	"github.com/snapp-incubator/snappcloud-status-backend/pkg/logger"
)

type Config struct {
	Querier *querier.Config `koanf:"querier"`
	Logger  *logger.Config  `koanf:"logger"`
}
