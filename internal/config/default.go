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
				{
					Order: 1,
					Name:  "PasS",
					Queries: struct {
						Outage     string "koanf:\"outage\""
						Disruption string "koanf:\"disruption\""
					}{
						Disruption: `up{job="node-exporter"}`, Outage: `up{job="node-exporter"}`},
				},
				{
					Order: 2,
					Name:  "IaaS",
					Queries: struct {
						Outage     string "koanf:\"outage\""
						Disruption string "koanf:\"disruption\""
					}{
						Disruption: `up{job="node-exporter"}`, Outage: `up{job="node-exporter"}`},
				},
				{
					Order: 3,
					Name:  "Object Storage (S3)",
					Queries: struct {
						Outage     string "koanf:\"outage\""
						Disruption string "koanf:\"disruption\""
					}{Disruption: `up{job="node-exporter"}`, Outage: `up{job="node-exporter"}`},
				},
				{
					Order: 4,
					Name:  "Container Registry",
					Queries: struct {
						Outage     string "koanf:\"outage\""
						Disruption string "koanf:\"disruption\""
					}{
						Disruption: `up{job="node-exporter"}`,
						Outage:     `up{job="node-exporter"}`,
					},
				},
				{
					Order: 5,
					Name:  "Service LoadBalancer (L4)",
					Queries: struct {
						Outage     string "koanf:\"outage\""
						Disruption string "koanf:\"disruption\""
					}{
						Disruption: `up{job="node-exporter"}`,
						Outage:     `up{job="node-exporter"}`,
					},
				},
				{
					Order: 6,
					Name:  "Ingress (L7)",
					Queries: struct {
						Outage     string "koanf:\"outage\""
						Disruption string "koanf:\"disruption\""
					}{
						Disruption: `up{job="node-exporter"}`,
						Outage:     `up{job="node-exporter"}`,
					},
				},
				{
					Order: 7,
					Name:  "Proxy",
					Queries: struct {
						Outage     string "koanf:\"outage\""
						Disruption string "koanf:\"disruption\""
					}{
						Disruption: `up{job="node-exporter"}`,
						Outage:     `up{job="node-exporter"}`,
					},
				},
				{
					Order: 8,
					Name:  "Monitoring",
					Queries: struct {
						Outage     string "koanf:\"outage\""
						Disruption string "koanf:\"disruption\""
					}{
						Disruption: `up{job="node-exporter"}`,
						Outage:     `up{job="node-exporter"}`,
					},
				},
				{
					Order: 9,
					Name:  "Logging",
					Queries: struct {
						Outage     string "koanf:\"outage\""
						Disruption string "koanf:\"disruption\""
					}{
						Disruption: `up{job="node-exporter"}`,
						Outage:     `up{job="node-exporter"}`,
					},
				},
				{
					Order: 10,
					Name:  "Traffic observability (Hubble)",
					Queries: struct {
						Outage     string "koanf:\"outage\""
						Disruption string "koanf:\"disruption\""
					}{
						Disruption: `up{job="node-exporter"}`,
						Outage:     `up{job="node-exporter"}`,
					},
				},
				{
					Order: 11,
					Name:  "ArgoCD",
					Queries: struct {
						Outage     string "koanf:\"outage\""
						Disruption string "koanf:\"disruption\""
					}{
						Disruption: `up{job="node-exporter"}`,
						Outage:     `up{job="node-exporter"}`,
					},
				},
				{
					Order: 12,
					Name:  "ArgoWF",
					Queries: struct {
						Outage     string "koanf:\"outage\""
						Disruption string "koanf:\"disruption\""
					}{
						Disruption: `up{job="node-exporter"}`,
						Outage:     `up{job="node-exporter"}`,
					},
				},
			},
		},
		Logger: &logger.Config{
			Development: true,
			Level:       "debug",
			Encoding:    "console",
		},
	}
}
