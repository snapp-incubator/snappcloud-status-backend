package querier

import (
	"sort"
	"time"

	"github.com/snapp-incubator/snappcloud-status-backend/internal/models"
)

type Querier interface {
	GetServices() []models.Service
	Start()
	Stop()
}

type querier struct {
	config   *Config
	ticker   *time.Ticker
	done     chan bool
	services []models.Service
}

func New(cfg *Config) Querier {
	return &querier{
		config:   cfg,
		ticker:   time.NewTicker(cfg.RequestInterval),
		done:     make(chan bool),
		services: initializeServices(cfg.Services),
	}
}

func initializeServices(config []ServiceConfig) []models.Service {
	services := make([]models.Service, 0, len(config))

	sort.Slice(config, func(i, j int) bool {
		return config[i].Order < config[j].Order
	})

	for index := 0; index < len(config); index++ {
		services = append(services, models.Service{
			Name:  config[index].Name,
			Query: config[index].Query,
			Status: map[models.Region]models.Status{
				models.Teh1:       models.Unknown,
				models.Teh2:       models.Unknown,
				models.SnappGroup: models.Unknown,
			},
		})
	}

	return services
}

// Start, run this function in a seperate goroutine
func (q *querier) Start() {
	for {
		select {
		case <-q.done:
			return
		case <-q.ticker.C:
			q.Query()
		}
	}
}

func (q *querier) Stop() {
	q.ticker.Stop()
	q.done <- true
}

func (q *querier) GetServices() []models.Service {
	return q.services
}
