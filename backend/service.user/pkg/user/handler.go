package user

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/kru-travel/airdrop-go/pkg/slogger"

	pb "github.com/jace-ys/kru-travel/backend/service.user/api/user"
)

func (u *userService) GetUser(ctx context.Context, r *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	slogger.Info().Log("event", "get_user.started", "user_id", r.Id)
	defer slogger.Info().Log("event", "get_user.finished", "user_id", r.Id)

	var user pb.GetUserResponse
	err := u.db.Transact(ctx, func(tx *sqlx.Tx) error {
		query := `
		SELECT u.id, u.username, u.email, u.name
		FROM users as u
		WHERE id=$1
		`
		row := tx.QueryRowx(query, r.Id)
		err := row.StructScan(&user)
		switch {
		case err == sql.ErrNoRows:
			return ErrUserNotFound
		case err != nil:
			return err
		}
		return nil
	})
	if err != nil {
		slogger.Info().Log("event", "get_user.failed", "user_id", r.Id, "msg", err)
		return nil, err
	}

	slogger.Info().Log("event", "get_user.success", "user_id", r.Id)
	return &user, nil
}
