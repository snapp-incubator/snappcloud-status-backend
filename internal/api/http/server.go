package http

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	app    *fiber.App
}

func New(log *zap.Logger) *Server {
	server := &Server{logger: log}

	server.app = fiber.New(fiber.Config{JSONEncoder: json.Marshal, JSONDecoder: json.Unmarshal})
	server.app.Use(cors.New())

	v1 := server.app.Group("api/v1")
	v1.Get("/services", server.services)

	return server
}

func (server *Server) Serve(port int) error {
	if err := server.app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		server.logger.Error("error resolving server", zap.Error(err))
		return err
	}
	return nil
}
