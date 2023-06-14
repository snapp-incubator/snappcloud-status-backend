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
	config *Config

	ticker *time.Ticker
	done   chan bool

	states   []state
	services []models.Service
}

type state struct {
	status map[models.Region]models.Status
	config ServiceConfig
}

func New(cfg *Config) Querier {
	instance := &querier{config: cfg}

	instance.ticker = time.NewTicker(cfg.RequestInterval)
	instance.done = make(chan bool)

	instance.initializeState()

	return instance
}

func (q *querier) initializeState() {
	sort.Slice(q.config.Services, func(i, j int) bool {
		return q.config.Services[i].Order < q.config.Services[j].Order
	})

	q.states = make([]state, 0, len(q.config.Services))
	for index := 0; index < len(q.config.Services); index++ {
		q.states = append(q.states, state{
			status: map[models.Region]models.Status{
				models.Teh1:       models.Unknown,
				models.Teh2:       models.Unknown,
				models.SnappGroup: models.Unknown,
			},
			config: q.config.Services[index],
		})
	}
}

// Start, run this function in a seperate goroutine
func (q *querier) Start() {
	for {
		select {
		case <-q.done:
			return
		case <-q.ticker.C:
			q.Query()
			q.generateServices()
		}
	}
}

func (q *querier) Stop() {
	q.ticker.Stop()
	q.done <- true
}

func (q *querier) generateServices() {
	services := make([]models.Service, 0, len(q.states))

	for index := 0; index < len(q.states); index++ {
		services = append(services, models.Service{
			Name:   q.states[index].config.Name,
			Status: q.states[index].status,
		})
	}

	q.services = services
}

func (q *querier) GetServices() []models.Service {
	return q.services
}
