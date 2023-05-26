package http

import (
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/snapp-incubator/snappcloud-status-backend/internal/models"
)

func (handler *Server) services(c *fiber.Ctx) error {
	templates := []struct {
		name      string
		operation func(map[models.Region]models.Status)
	}{
		{name: "PasS", operation: pass},
		{name: "IaaS", operation: iaas},
	}

	services := make([]models.Service, 0, len(templates))
	var wg sync.WaitGroup
	wg.Add(len(templates))

	for index := 0; index < len(templates); index++ {
		go func(index int) {
			defer wg.Done()

			result := map[models.Region]models.Status{
				models.Teh1:       models.Unknown,
				models.Teh2:       models.Unknown,
				models.SnappGroup: models.Unknown,
			}
			templates[index].operation(result)

			services = append(services, models.Service{
				Name:   templates[index].name,
				Order:  index + 1,
				Status: result,
			})
		}(index)
	}

	wg.Wait()

	return c.Status(http.StatusOK).JSON(&services)
}

func pass(result map[models.Region]models.Status) {
	// TODO: implement
	result[models.Teh1] = models.Operational
	result[models.Teh2] = models.Warning
	result[models.SnappGroup] = models.Operational
}

func iaas(result map[models.Region]models.Status) {
	// TODO: implement
	result[models.Teh1] = models.Operational
	result[models.Teh2] = models.Operational
	result[models.SnappGroup] = models.Outage
}
