package main

import (
	"os"

	"github.com/jace-ys/kru-travel/backend/service.user/pkg/user"
	"github.com/kru-travel/airdrop-go/pkg/config"
	"github.com/kru-travel/airdrop-go/pkg/server"
	"github.com/kru-travel/airdrop-go/pkg/slogger"
)

func init() {
	config.LoadFile("config/config.yaml")
}

func main() {
	port := config.Get("port").Int(8080)
	userService := user.NewService()
	if err := server.ListenAndServe(userService, port); err != nil {
		slogger.Error().Log("msg", err)
		os.Exit(1)
	}
}
