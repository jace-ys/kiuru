package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/kru-travel/airdrop-go/pkg/gorpc"
	"github.com/kru-travel/airdrop-go/pkg/slogger"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"

	pb "github.com/jace-ys/kru-travel/backend/service.user/api/user"
)

func (u *userService) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	slogger.Info().Log("event", "get_all_users.started")
	defer slogger.Info().Log("event", "get_all_users.finished")

	users, err := u.getAllUsers(ctx)
	if err != nil {
		slogger.Error().Log("event", "get_all_users.failed", "msg", err)
		return nil, gorpc.InternalError()
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
		defer rows.Close()
		for rows.Next() {
			var user pb.User
			if err := rows.StructScan(&user); err != nil {
				return err
			}
			users = append(users, &user)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	slogger.Info().Log("event", "get_user.started", "user_id", req.Id)
	defer slogger.Info().Log("event", "get_user.finished", "user_id", req.Id)

	user, err := u.getUser(ctx, req.Id)
	if err != nil {
		slogger.Error().Log("event", "get_user.failed", "user_id", req.Id, "msg", err)
		switch {
		case errors.Is(err, ErrUserNotFound):
			return nil, gorpc.Error(codes.NotFound, err)
		default:
			return nil, gorpc.InternalError()
		}
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
		case errors.Is(err, sql.ErrNoRows):
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

func (u *userService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	slogger.Info().Log("event", "create_user.started")
	defer slogger.Info().Log("event", "create_user.finished")

	userId, err := u.createUser(ctx, req.User)
	if err != nil {
		slogger.Error().Log("event", "create_user.failed", "msg", err)
		switch {
		case errors.Is(err, ErrUserExists):
			return nil, gorpc.Error(codes.AlreadyExists, err)
		default:
			return nil, gorpc.InternalError()
		}
	}

	slogger.Info().Log("event", "create_user.success", "user_id", userId)
	return &pb.CreateUserResponse{
		Id: userId,
	}, nil
}

func (u *userService) createUser(ctx context.Context, user *pb.User) (string, error) {
	var userId string
	err := u.db.Transact(ctx, func(tx *sqlx.Tx) error {
		query := `
		INSERT INTO users (username, email, name)
		VALUES (:username, :email, :name)
		RETURNING id
		`
		stmt, err := tx.PrepareNamed(query)
		if err != nil {
			return err
		}
		err = stmt.QueryRowx(user).Scan(&userId)
		if err != nil {
			var pqErr *pq.Error
			switch {
			case errors.As(err, &pqErr) && pqErr.Code == "23505":
				return ErrUserExistsCtx(pqErr)
			default:
				return err
			}
		}
		return nil
	})
	if err != nil {
		return userId, err
	}
	return userId, nil
}

func (u *userService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	slogger.Info().Log("event", "delete_user.started", "user_id", req.Id)
	defer slogger.Info().Log("event", "delete_user.finished", "user_id", req.Id)

	err := u.deleteUser(ctx, req.Id)
	if err != nil {
		slogger.Error().Log("event", "delete_user.failed", "user_id", req.Id, "msg", err)
		switch {
		case errors.Is(err, ErrUserNotFound):
			return nil, gorpc.Error(codes.NotFound, err)
		default:
			return nil, gorpc.InternalError()
		}
	}

	slogger.Info().Log("event", "delete_user.success", "user_id", req.Id)
	return &pb.DeleteUserResponse{}, nil
}

func (u *userService) deleteUser(ctx context.Context, userId string) error {
	return u.db.Transact(ctx, func(tx *sqlx.Tx) error {
		query := `
		DELETE FROM users
		WHERE id=$1
		`
		res, err := tx.Exec(query, userId)
		if err != nil {
			return err
		}
		count, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if count == 0 {
			return ErrUserNotFound
		}
		return nil
	})
}
