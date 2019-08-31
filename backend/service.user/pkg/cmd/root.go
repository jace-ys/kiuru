package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/kru-travel/airdrop-go/pkg/slogger"

	"github.com/jace-ys/kru-travel/backend/service.user/pkg/server"
	"github.com/jace-ys/kru-travel/backend/service.user/pkg/user"
)

type serviceConfig struct {
	serverPort  int
	gatewayPort int
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

	rootCmd.PersistentFlags().IntVar(&config.serverPort, "port", 8080, "binding port for the gRPC server")
	rootCmd.PersistentFlags().IntVar(&config.gatewayPort, "gateway-port", 8081, "binding port for the REST gateway proxy")

	return rootCmd
}
