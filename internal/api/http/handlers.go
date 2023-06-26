package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snapp-incubator/snappcloud-status-backend/internal/models"
)

func (server *Server) liveness(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (server *Server) readiness(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (server *Server) services(c *gin.Context) {
	c.JSON(http.StatusOK, &struct {
		Message  string           `json:"message"`
		Services []models.Service `json:"services,omitempty"`
	}{
		Message:  "All services retrieved successfully.",
		Services: server.querier.GetServices(),
	})
}
