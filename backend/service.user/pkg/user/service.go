package user

import (
	"context"

	"github.com/jace-ys/kru-travel/backend/service.user/pkg/crdb"
	"github.com/jace-ys/kru-travel/backend/service.user/pkg/server"
)

type userService struct {
	db *crdb.DB
}

func NewService() *userService {
	return &userService{}
}

func (u *userService) Init(crdbConfig crdb.Config) error {
	db, err := crdb.Connect(crdbConfig)
	if err != nil {
		return err
	}
	u.db = db
	return nil
}

func (u *userService) StartServer(ctx context.Context, s server.Server, port int) error {
	if err := s.Init(ctx, u); err != nil {
		return err
	}
	defer s.Shutdown(ctx)
	return s.Serve(port)
}

func (u *userService) Teardown() {
	u.db.Close()
}
