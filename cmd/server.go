package cmd

import (
	"os"

	"github.com/snapp-incubator/snappcloud-status-backend/internal/api/http"
	"github.com/snapp-incubator/snappcloud-status-backend/internal/config"
	"github.com/snapp-incubator/snappcloud-status-backend/internal/querier"
	"github.com/snapp-incubator/snappcloud-status-backend/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Server struct{}

func (cmd Server) Command(trap chan os.Signal) *cobra.Command {
	run := func(_ *cobra.Command, _ []string) {
		cmd.main(config.Load(true), trap)
	}

	return &cobra.Command{
		Use:   "server",
		Short: "run PhoneBook server",
		Run:   run,
	}
}

func (cmd *Server) main(cfg *config.Config, trap chan os.Signal) {
	logger := logger.NewZap(cfg.Logger)

	querier := querier.New(cfg.Querier, logger)
	go querier.Start()

	server := http.New(logger, querier)
	go server.Serve(8080)

	// Keep this at the bottom of the main function
	field := zap.String("signal trap", (<-trap).String())
	logger.Info("exiting by receiving a unix signal", field)

	querier.Stop()
}
