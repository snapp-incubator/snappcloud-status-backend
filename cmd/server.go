package cmd

import (
	"fmt"
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
	loggerObj := logger.NewZap(cfg.Logger)

	querierObj := querier.New(cfg.Querier, loggerObj)
	go querierObj.Start()

	server := http.New(loggerObj, querierObj)
	go func() {
		err := server.Serve(8080)
		if err != nil {
			panic(fmt.Errorf("error on Serve function: %s", err))
		}
	}()

	// Keep this at the bottom of the main function
	field := zap.String("signal trap", (<-trap).String())
	loggerObj.Info("exiting by receiving a unix signal", field)

	querierObj.Stop()
}
