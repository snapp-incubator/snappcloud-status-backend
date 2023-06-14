package config

import (
	"time"

	"github.com/snapp-incubator/snappcloud-status-backend/internal/querier"
	"github.com/snapp-incubator/snappcloud-status-backend/pkg/logger"
)

func Default() *Config {
	return &Config{
		Querier: &querier.Config{
			RequestInterval: time.Second * 10,
			ThanosFrontends: struct {
				Teh1       string "koanf:\"teh1\""
				Teh2       string "koanf:\"teh2\""
				SnappGroup string "koanf:\"snappgroup\""
			}{Teh1: "localhost:9090", Teh2: "localhost:9090", SnappGroup: "localhost:9090"},
			Services: []querier.ServiceConfig{
				{Order: 1, Name: "PasS", Query: `up{job="node-exporter"}`},
				{Order: 2, Name: "IaaS", Query: `up{job="node-exporter"}`},
				{Order: 3, Name: "Object Storage (S3)", Query: `up{job="node-exporter"}`},
				{Order: 4, Name: "Container Registry", Query: `up{job="node-exporter"}`},
				{Order: 5, Name: "Service LoadBalancer (L4)", Query: `up{job="node-exporter"}`},
				{Order: 6, Name: "Ingress (L7)", Query: `up{job="node-exporter"}`},
				{Order: 7, Name: "Proxy", Query: `up{job="node-exporter"}`},
				{Order: 8, Name: "Monitoring", Query: `up{job="node-exporter"}`},
				{Order: 9, Name: "Logging", Query: `up{job="node-exporter"}`},
				{Order: 10, Name: "Traffic observability (Hubble)", Query: `up{job="node-exporter"}`},
				{Order: 11, Name: "ArgoCD", Query: `up{job="node-exporter"}`},
				{Order: 12, Name: "ArgoWF", Query: `up{job="node-exporter"}`},
			},
		},
		Logger: &logger.Config{
			Development: true,
			Level:       "debug",
			Encoding:    "console",
		},
	}
}
