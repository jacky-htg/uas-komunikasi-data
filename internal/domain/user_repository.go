package domain

import (
	"context"
	"database/sql"
	"uas-komdat/pb/users"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRepository struct {
	db *sql.DB
	pb users.User
}

func (a *UserRepository) find(ctx context.Context) error {
	query := `
		SELECT name, photo
	 	FROM users WHERE id = $1
	`
	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement find: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, a.pb.Id).Scan(
		&a.pb.Name,
		&a.pb.Photo,
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Query Row Context find: %v", err)
	}

	return nil
}

func (a *UserRepository) list(ctx context.Context) (*sql.Rows, error) {
	query := `
		SELECT id, name, photo FROM users
	`
	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Prepare statement find: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Query Row Context find: %v", err)
	}

	return rows, nil
}
