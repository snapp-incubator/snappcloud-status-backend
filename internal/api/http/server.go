package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/snapp-incubator/snappcloud-status-backend/internal/querier"
	"go.uber.org/zap"
)

type Server struct {
	logger  *zap.Logger
	querier querier.Querier
	engine  *gin.Engine
}

func New(log *zap.Logger, querier querier.Querier) *Server {
	server := &Server{logger: log, querier: querier, engine: gin.Default()}

	server.engine.Use(server.corsMiddleware())

	v1 := server.engine.Group("api/v1")
	v1.GET("/services", server.services)

	healthz := server.engine.Group("healthz")
	healthz.GET("/liveness", server.liveness)
	healthz.GET("/readiness", server.readiness)

	return server
}

func (server *Server) Serve(port int) error {
	if err := server.engine.Run(fmt.Sprintf(":%d", port)); err != nil {
		server.logger.Error("error resolving server", zap.Error(err))
		return err
	}
	return nil
}
