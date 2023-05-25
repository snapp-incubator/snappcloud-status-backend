package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/snapp-incubator/snappcloud-status-backend/internal/models"
)

func (handler *Server) services(c *fiber.Ctx) error {
	response := []models.Service{
		models.Service{
			Name: "Paas",
			States: []models.State{
				models.State{
					Region: models.Teh1,
					Status: models.Operational,
				},
				models.State{
					Region: models.Teh2,
					Status: models.Operational,
				},
				models.State{
					Region: models.SnappGroup,
					Status: models.Operational,
				},
			},
		},
	}

	return c.Status(http.StatusOK).JSON(&response)
}
