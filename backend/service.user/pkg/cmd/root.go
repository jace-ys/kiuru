package cmd

import (
	"os"

	"github.com/kru-travel/airdrop-go/pkg/slogger"
	"github.com/spf13/cobra"

	"github.com/jace-ys/kru-travel/backend/service.user/pkg/server"
	"github.com/jace-ys/kru-travel/backend/service.user/pkg/user"
)

type serviceConfig struct {
	serverPort int
}

func NewCmd() *cobra.Command {
	var config serviceConfig

	rootCmd := &cobra.Command{
		Use:   "server",
		Short: "Start the server for the service",
		Run: func(cmd *cobra.Command, args []string) {
			userService := user.NewService()
			if err := server.ListenAndServe(userService, config.serverPort); err != nil {
				slogger.Error().Log("msg", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.PersistentFlags().IntVar(&config.serverPort, "port", 8080, "Binding port for the server")

	return rootCmd
}
