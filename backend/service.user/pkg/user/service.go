package user

import "google.golang.org/grpc"

type Service struct {
	*grpc.Server
}

func NewService() *Service {
	return &Service{grpc.NewServer()}
}

func (s *Service) Init() error {
	return nil
}

func (s *Service) Shutdown() {
	s.Server.Stop()
}
