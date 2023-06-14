package querier

import (
	"sort"

	"github.com/snapp-incubator/snappcloud-status-backend/internal/models"
)

type Querier interface {
	GetState() map[string]any
}

type querier struct {
	config   *Config
	services []models.Service
}

func New() Querier {
	instance := &querier{}

	return instance
}

func (q *querier) buildInitialState() {
	services := make([]models.Service, 0, len(q.config.services))

	sort.Slice(q.config.services, func(i, j int) bool {
		return q.config.services[i].Order < q.config.services[j].Order
	})

	q.services = services
}

func sortServices() {

}

func (q *querier) GetState() map[string]any {
	return map[string]any{}
}
