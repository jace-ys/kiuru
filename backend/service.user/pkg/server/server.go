package server

import (
	"fmt"
	"net"

	"github.com/pkg/errors"

	"github.com/kru-travel/airdrop-go/pkg/slogger"
)

type Server interface {
	Init() error
	Serve(lis net.Listener) error
	Shutdown()
}

func ListenAndServe(service Server, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}

	err = service.Init()
	if err != nil {
		return errors.Wrap(err, "initialising service")
	}
	defer service.Shutdown()

	slogger.Info().Log("event", "server.started", "port", port)
	defer slogger.Info().Log("event", "server.stopped")
	if err := service.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to serve")
	}
	return nil
}
