package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/snapp-incubator/snappcloud-status-backend/internal/models"
)

func (handler *Server) services(c *fiber.Ctx) error {
	response := []models.Service{
		{
			Name: "Paas",
			Status: map[models.Region]models.Status{
				models.Teh1:       models.Operational,
				models.Teh2:       models.Operational,
				models.SnappGroup: models.Operational,
			},
		},
	}

	return c.Status(http.StatusOK).JSON(&response)
}
