package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kru-travel/airdrop-go/pkg/gorpc"
	"github.com/kru-travel/airdrop-go/pkg/slogger"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"

	pb "github.com/jace-ys/kru-travel/backend/service.user/api/user"
)

func (s *userService) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	slogger.Info().Log("event", "get_all_users.started")
	defer slogger.Info().Log("event", "get_all_users.finished")

	users, err := s.getAllUsers(ctx)
	if err != nil {
		slogger.Error().Log("event", "get_all_users.failed", "msg", err)
		return nil, gorpc.InternalError()
	}

	slogger.Info().Log("event", "get_all_users.success")
	return &pb.GetAllUsersResponse{
		Users: users,
	}, nil
}

func (s *userService) getAllUsers(ctx context.Context) ([]*pb.User, error) {
	var users []*pb.User
	err := s.db.Transact(ctx, func(tx *sqlx.Tx) error {
		query := `
		SELECT u.id, u.username, u.email, u.name
		FROM users as u
		`
		rows, err := tx.QueryxContext(ctx, query)
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

func (s *userService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	slogger.Info().Log("event", "get_user.started", "user_id", req.Id)
	defer slogger.Info().Log("event", "get_user.finished", "user_id", req.Id)

	user, err := s.getUser(ctx, req.Id)
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

func (s *userService) getUser(ctx context.Context, userId string) (*pb.User, error) {
	var user pb.User
	err := s.db.Transact(ctx, func(tx *sqlx.Tx) error {
		query := `
		SELECT u.id, u.username, u.email, u.name
		FROM users as u
		WHERE id=$1
		`
		row := tx.QueryRowxContext(ctx, query, userId)
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

func (s *userService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	slogger.Info().Log("event", "create_user.started")
	defer slogger.Info().Log("event", "create_user.finished")

	err := s.validateUser(req.User)
	if err != nil {
		slogger.Error().Log("event", "create_user.failed", "msg", err)
		return nil, gorpc.Error(codes.InvalidArgument, err)
	}

	err = s.hashPassword(req.User)
	if err != nil {
		slogger.Error().Log("event", "create_user.failed", "msg", err)
		return nil, gorpc.InternalError()
	}

	userId, err := s.createUser(ctx, req.User)
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

func (s *userService) validateUser(user *pb.User) error {
	switch {
	case user.Username == "":
		return ErrInvalidRequestCtx(`missing "username" field`)
	case user.Password == "":
		return ErrInvalidRequestCtx(`missing "password" field`)
	case user.Email == "":
		return ErrInvalidRequestCtx(`missing "email" field`)
	case user.Name == "":
		return ErrInvalidRequestCtx(`missing "name" field`)
	}
	return nil
}

func (s *userService) hashPassword(user *pb.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}
	user.Password = string(hashedPassword)
	return nil
}

func (s *userService) createUser(ctx context.Context, user *pb.User) (string, error) {
	var userId string
	err := s.db.Transact(ctx, func(tx *sqlx.Tx) error {
		query := `
		INSERT INTO users (username, password, email, name)
		VALUES (:username, :password, :email, :name)
		RETURNING id
		`
		stmt, err := tx.PrepareNamedContext(ctx, query)
		if err != nil {
			return err
		}
		err = stmt.QueryRowxContext(ctx, user).Scan(&userId)
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

func (s *userService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	slogger.Info().Log("event", "delete_user.started", "user_id", req.Id)
	defer slogger.Info().Log("event", "delete_user.finished", "user_id", req.Id)

	err := s.deleteUser(ctx, req.Id)
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

func (s *userService) deleteUser(ctx context.Context, userId string) error {
	return s.db.Transact(ctx, func(tx *sqlx.Tx) error {
		query := `
		DELETE FROM users
		WHERE id=$1
		`
		res, err := tx.ExecContext(ctx, query, userId)
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
