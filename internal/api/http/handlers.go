package http

import (
	"net/http"

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
	c.JSON(http.StatusOK, &struct {
		Message  string           `json:"message"`
		Services []models.Service `json:"services,omitempty"`
	}{
		Message:  "All services retrieved successfuly.",
		Services: handler.querier.GetServices(),
	})
}
