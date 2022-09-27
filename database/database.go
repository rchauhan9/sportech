package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

func NewDatabasePool(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(ctx, connectionString)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to database %s", connectionString)
	}
	return db, nil
}
