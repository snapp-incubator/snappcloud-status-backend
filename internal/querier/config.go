package querier

import "time"

type Config struct {
	ThanosFrontends struct {
		Teh1       string `koanf:"teh1"`
		Teh2       string `koanf:"teh2"`
		SnappGroup string `koanf:"snappgroup"`
	} `koanf:"thanos_frontends"`
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