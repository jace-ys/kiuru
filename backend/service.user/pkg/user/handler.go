package user

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/kru-travel/airdrop-go/pkg/slogger"
	"github.com/pkg/errors"

	pb "github.com/jace-ys/kru-travel/backend/service.user/api/user"
)

func (u *userService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	slogger.Info().Log("event", "get_user.started", "user_id", req.Id)
	defer slogger.Info().Log("event", "get_user.finished", "user_id", req.Id)

	user, err := u.getUser(ctx, req.Id)
	if err != nil {
		slogger.Error().Log("event", "get_user.failed", "user_id", req.Id, "msg", err)
		return nil, errors.Wrap(err, "failed to get user")
	}

	slogger.Info().Log("event", "get_user.success", "user_id", req.Id)
	return &pb.GetUserResponse{
		User: user,
	}, nil
}

func (u *userService) getUser(ctx context.Context, userId string) (*pb.User, error) {
	var user pb.User
	err := u.db.Transact(ctx, func(tx *sqlx.Tx) error {
		query := `
		SELECT u.id, u.username, u.email, u.name
		FROM users as u
		WHERE id=$1
		`
		row := tx.QueryRowx(query, userId)
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
		return nil, err
	}
	return &user, nil
}

func (u *userService) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	slogger.Info().Log("event", "get_all_users.started")
	defer slogger.Info().Log("event", "get_all_users.finished")

	users, err := u.getAllUsers(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not get all users")
	}

	slogger.Info().Log("event", "get_all_users.success")
	return &pb.GetAllUsersResponse{
		Users: users,
	}, nil
}

func (u *userService) getAllUsers(ctx context.Context) ([]*pb.User, error) {
	var users []*pb.User
	err := u.db.Transact(ctx, func(tx *sqlx.Tx) error {
		query := `
		SELECT u.id, u.username, u.email, u.name
		FROM users as u
		`
		rows, err := tx.Queryx(query)
		if err != nil {
			return err
		}
		for rows.Next() {
			var user pb.User
			if err := rows.StructScan(&user); err != nil {
				return err
			}
			users = append(users, &user)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}
