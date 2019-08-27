package service

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/kru-travel/airdrop-go/pkg/slogger"
)

type Server struct {
	*grpc.Server
}

func NewServer() *Server {
	return &Server{grpc.NewServer()}
}

func (s *Server) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	slogger.Info().Log("event", "server.started", "port", port)
	if err := s.Serve(lis); err != nil {
		return err
	}
	return nil
}
