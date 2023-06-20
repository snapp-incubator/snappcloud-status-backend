package querier

import (
	"sort"
	"sync"
	"time"

	"github.com/snapp-incubator/snappcloud-status-backend/internal/models"
	"go.uber.org/zap"
)

type Querier interface {
	GetServices() []models.Service
	Start()
	Stop()
}

type querier struct {
	config  *Config
	loggger *zap.Logger

	ticker *time.Ticker
	done   chan bool

	states   []state
	mutex    sync.RWMutex
	services []models.Service
}

type state struct {
	status models.Status
	config ServiceConfig
}

func New(cfg *Config, lg *zap.Logger) Querier {
	instance := &querier{config: cfg, loggger: lg}

	instance.ticker = time.NewTicker(cfg.RequestInterval)
	instance.done = make(chan bool)

	instance.mutex = sync.RWMutex{}
	instance.initializeState()
	instance.generateServices()

	return instance
}

func (q *querier) initializeState() {
	sort.Slice(q.config.Services, func(i, j int) bool {
		return q.config.Services[i].Order < q.config.Services[j].Order
	})

	q.states = make([]state, 0, len(q.config.Services))
	for index := 0; index < len(q.config.Services); index++ {
		q.states = append(q.states, state{
			status: models.Unknown,
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
		q.mutex.RLock()
		services = append(services, models.Service{
			Name:   q.states[index].config.Name,
			Status: q.states[index].status,
		})
		q.mutex.RUnlock()
	}

	q.services = services
}

func (q *querier) GetServices() []models.Service {
	return q.services
}
