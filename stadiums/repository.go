package stadiums

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Repository interface {
	ListStadiums(ctx context.Context) ([]Stadium, error)
	GetStadium(ctx context.Context, id string) (Stadium, error)
}

func NewRepository(dbPool *pgxpool.Pool) Repository {
	return &repository{dbPool}
}

type repository struct {
	pool *pgxpool.Pool
}

func (r *repository) ListStadiums(ctx context.Context) ([]Stadium, error) {
	query := `
		SELECT
		    id,
		    name,
		    capacity,
		    city,
		    country_id
	    FROM stadiums
	    ORDER BY name ASC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "errpr fetching discover category details from database")
	}

	var stadiums []Stadium
	for rows.Next() {
		stadium := Stadium{}
		if err := rows.Scan(
			&stadium.ID,
			&stadium.Name,
			&stadium.Capacity,
			&stadium.City,
			&stadium.Country,
		); err != nil {
			return nil, errors.Wrap(err, "error scanning row from database")
		}
		stadiums = append(stadiums, stadium)
	}
	return stadiums, nil
}

func (r *repository) GetStadium(ctx context.Context, id string) (Stadium, error) {
	query := `
	    SELECT
		    id,
		    name,
		    capacity,
		    city,
		    country_id
	    FROM stadiums
	    WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)
	var stadium Stadium
	if err := row.Scan(
		&stadium.ID,
		&stadium.Name,
		&stadium.Capacity,
		&stadium.City,
		&stadium.Country,
	); err != nil {
		return stadium, errors.Wrapf(err, "error getting stadium with id %s", id)
	}
	return stadium, nil
}
