package querier

import "time"

type Config struct {
	ThanosFrontend  string          `koanf:"thanos_frontend"`
	RequestInterval time.Duration   `koanf:"request_interval"`
	RequestTimeout  time.Duration   `koanf:"request_timeout"`
	Services        []ServiceConfig `koanf:"services"`
}

type ServiceConfig struct {
	Name    string `koanf:"name"`
	Order   int    `koanf:"order"`
	Queries struct {
		Outage     string `koanf:"outage"`
		Disruption string `koanf:"disruption"`
	} `koanf:"queries"`
}
