package http

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/snapp-incubator/snappcloud-status-backend/internal/models"
)

func (handler *Server) liveness(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (handler *Server) readiness(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (handler *Server) services(c *gin.Context) {
	templates := []struct {
		name      string
		operation func(map[models.Region]models.Status)
	}{
		// The services goes here with the order you want to display in the UI
		{name: "PasS", operation: pass},
		{name: "IaaS", operation: iaas},
		{name: "Object Storage (S3)", operation: ok},
		{name: "Container Registry", operation: ok},
		{name: "Service LoadBalancer (L4)", operation: ok},
		{name: "Ingress (L7)", operation: ok},
		{name: "Proxy", operation: ok},
		{name: "Monitoring", operation: ok},
		{name: "Logging", operation: ok},
		{name: "Traffic observability (Hubble)", operation: ok},
		{name: "ArgoCD", operation: ok},
		{name: "ArgoWF", operation: ok},
	}

	type response struct {
		Services []models.Service `json:"services,omitempty"`
		Message  string           `json:"message"`
		Error    string           `json:"error,omitempty"`
	}

	services := make([]models.Service, len(templates))
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

			services[index] = models.Service{
				Name:   templates[index].name,
				Status: result,
			}
		}(index)
	}

	wg.Wait()

	// c.Response().Header.Add("Content-Type", "application/json")

	c.JSON(http.StatusOK, &response{
		Message:  "All services retrieved successfuly.",
		Services: services,
	})
}

// TODO: implement
func pass(status map[models.Region]models.Status) {
	status[models.Teh1] = models.Operational
	status[models.Teh2] = models.Warning
	status[models.SnappGroup] = models.Operational
}

// TODO: implement
func iaas(status map[models.Region]models.Status) {
	status[models.Teh1] = models.Operational
	status[models.Teh2] = models.Operational
	status[models.SnappGroup] = models.Outage
}

// TODO: DELETE
func ok(status map[models.Region]models.Status) {
	status[models.Teh1] = models.Operational
	status[models.Teh2] = models.Operational
	status[models.SnappGroup] = models.Operational
}
