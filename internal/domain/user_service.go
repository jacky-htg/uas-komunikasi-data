package domain

import (
	"context"
	"database/sql"
	"uas-komdat/internal/pkg/storage"
	"uas-komdat/pb/users"

	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	Db          *sql.DB
	MinioClient *minio.Client
	users.UnimplementedUserServiceServer
}

func (a *UserService) List(ctx context.Context, in *users.EmptyMessage) (*users.Users, error) {
	var list []*users.User
	userRepo := UserRepository{db: a.Db}
	rows, err := userRepo.list(ctx)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user users.User
		err = rows.Scan(&user.Id, &user.Name, &user.Photo)
		if err != nil {
			return nil, err
		}

		user.Photo, err = storage.SignUrl(ctx, a.MinioClient, user.Photo)
		if err != nil {
			return nil, err
		}

		list = append(list, &user)
	}

	if rows.Err() != nil {
		return nil, status.Errorf(codes.Internal, "Row : %v", err)
	}
	return &users.Users{User: list}, nil
}

func (a *UserService) ListStreaming(empty *users.EmptyMessage, stream users.UserService_ListStreamingServer) error {
	ctx := stream.Context()
	userRepo := UserRepository{db: a.Db}
	rows, err := userRepo.list(ctx)
	if err != nil {
		return err
	}

	for rows.Next() {
		var user users.User
		err = rows.Scan(&user.Id, &user.Name, &user.Photo)
		if err != nil {
			return err
		}

		user.Photo, err = storage.SignUrl(ctx, a.MinioClient, user.Photo)
		if err != nil {
			return err
		}

		if err := stream.Send(&user); err != nil {
			return err
		}
	}

	return nil
}
