package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-kit/kit/log/level"
	"github.com/jmoiron/sqlx"
	"github.com/kru-travel/airdrop-go/pkg/gorpc"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"

	pb "github.com/jace-ys/kru-travel/backend/service.user/api/user"
	"github.com/jace-ys/kru-travel/backend/service.user/pkg/permissions"
)

func (s *userService) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	level.Info(s.logger).Log("event", "get_all_users.started")
	defer level.Info(s.logger).Log("event", "get_all_users.finished")

	users, err := s.getAllUsers(ctx)
	if err != nil {
		level.Error(s.logger).Log("event", "get_all_users.failed", "msg", err)
		return nil, gorpc.Error(err)
	}

	level.Info(s.logger).Log("event", "get_all_users.success")
	return &pb.GetAllUsersResponse{
		Users: users,
	}, nil
}

func (s *userService) getAllUsers(ctx context.Context) ([]*pb.User, error) {
	var users []*pb.User
	err := s.db.Transact(ctx, func(tx *sqlx.Tx) error {
		query := `
		SELECT u.id, u.created_at, u.username, u.email, u.name
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
		return nil, gorpc.NewErr(codes.Internal, err)
	}
	return users, nil
}

func (s *userService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	level.Info(s.logger).Log("event", "get_user.started")
	defer level.Info(s.logger).Log("event", "get_user.finished")

	user, err := s.getUser(ctx, req.Id)
	if err != nil {
		level.Error(s.logger).Log("event", "get_user.failed", "msg", err)
		return nil, gorpc.Error(err)
	}

	level.Info(s.logger).Log("event", "get_user.success")
	return &pb.GetUserResponse{
		User: user,
	}, nil
}

func (s *userService) getUser(ctx context.Context, userId string) (*pb.User, error) {
	var user pb.User
	err := s.db.Transact(ctx, func(tx *sqlx.Tx) error {
		query := `
		SELECT u.id, u.created_at, u.username, u.email, u.name
		FROM users as u
		WHERE id=$1
		`
		row := tx.QueryRowxContext(ctx, query, userId)
		err := row.StructScan(&user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, gorpc.NewErr(codes.NotFound, ErrUserNotFound)
		case errors.As(err, &pqErr) && pqErr.Code.Name() == "protocol_violation":
			return nil, gorpc.NewErr(codes.NotFound, ErrUserNotFound)
		default:
			return nil, gorpc.NewErr(codes.Internal, err)
		}
	}
	return &user, nil
}

func (s *userService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	level.Info(s.logger).Log("event", "create_user.started")
	defer level.Info(s.logger).Log("event", "create_user.finished")

	err := s.validateUserPayload(req.User)
	if err != nil {
		level.Error(s.logger).Log("event", "create_user.failed", "msg", err)
		return nil, gorpc.Error(err)
	}

	err = s.hashPassword(req.User)
	if err != nil {
		level.Error(s.logger).Log("event", "create_user.failed", "msg", err)
		return nil, gorpc.Error(err)
	}

	userId, err := s.createUser(ctx, req.User)
	if err != nil {
		level.Error(s.logger).Log("event", "create_user.failed", "msg", err)
		return nil, gorpc.Error(err)
	}

	level.Info(s.logger).Log("event", "create_user.success")
	return &pb.CreateUserResponse{
		Id: userId,
	}, nil
}

func (s *userService) validateUserPayload(user *pb.User) error {
	switch {
	case user.Username == "":
		return gorpc.NewErr(codes.InvalidArgument, fmt.Errorf("%w: %s", ErrInvalidRequest, `missing "username" field`))
	case user.Password == "":
		return gorpc.NewErr(codes.InvalidArgument, fmt.Errorf("%w: %s", ErrInvalidRequest, `missing "password" field`))
	case user.Email == "":
		return gorpc.NewErr(codes.InvalidArgument, fmt.Errorf("%w: %s", ErrInvalidRequest, `missing "email" field`))
	case user.Name == "":
		return gorpc.NewErr(codes.InvalidArgument, fmt.Errorf("%w: %s", ErrInvalidRequest, `missing "name" field`))
	}
	return nil
}

func (s *userService) hashPassword(user *pb.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return gorpc.NewErr(codes.Internal, fmt.Errorf("%w: %w", ErrHashingPssword, err))
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
			return err
		}
		return nil
	})
	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation":
			return "", gorpc.NewErr(codes.AlreadyExists, ErrUserExistsContext(pqErr))
		default:
			return "", gorpc.NewErr(codes.Internal, err)
		}
	}
	return userId, nil
}

func (s *userService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	level.Info(s.logger).Log("event", "delete_user.started")
	defer level.Info(s.logger).Log("event", "delete_user.finished")

	err := s.verifyPermissions(ctx, permissions.UserScope, req.Id)
	if err != nil {
		level.Error(s.logger).Log("event", "delete_user.failed", "msg", err)
		return nil, gorpc.Error(err)
	}

	err = s.deleteUser(ctx, req.Id)
	if err != nil {
		level.Error(s.logger).Log("event", "delete_user.failed", "msg", err)
		return nil, gorpc.Error(err)
	}

	level.Info(s.logger).Log("event", "delete_user.success")
	return &pb.DeleteUserResponse{}, nil
}

func (s *userService) deleteUser(ctx context.Context, userId string) error {
	err := s.db.Transact(ctx, func(tx *sqlx.Tx) error {
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
	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.Is(err, ErrUserNotFound):
			return gorpc.NewErr(codes.NotFound, ErrUserNotFound)
		case errors.As(err, &pqErr) && pqErr.Code.Name() == "protocol_violation":
			return gorpc.NewErr(codes.NotFound, ErrUserNotFound)
		default:
			return gorpc.NewErr(codes.Internal, err)
		}
	}
	return nil
}

func (s *userService) verifyPermissions(ctx context.Context, scopeFunc permissions.ScopeFunc, userId string) error {
	userMD, err := gorpc.GetUserMD(ctx)
	if err != nil {
		return gorpc.NewErr(codes.Internal, err)
	}

	if !scopeFunc(userMD, userId) {
		return gorpc.NewErr(codes.PermissionDenied, ErrPermissionDenied)
	}

	return nil
}
