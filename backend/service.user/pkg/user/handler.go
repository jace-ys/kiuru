package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/kru-travel/airdrop-go/pkg/slogger"

	pb "github.com/jace-ys/kru-travel/backend/service.user/api/user"
)

func (u *userService) GetUser(ctx context.Context, r *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var user pb.GetUserResponse
	err := u.db.Transact(func(tx *sqlx.Tx) error {
		row := u.db.QueryRowx("SELECT * FROM users WHERE id=$1", r.Id)
		err := row.StructScan(&user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		slogger.Info().Log("event", "request.get_user_failure", "user_id", r.Id)
		return nil, err
	}
	slogger.Info().Log("event", "request.get_user_success", "user_id", r.Id)
	return &user, nil
}
