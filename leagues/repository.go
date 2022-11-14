package leagues

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Repository interface {
	ListLeagues(ctx context.Context) ([]League, error)
	GetLeague(ctx context.Context, id string) (League, error)
}

func NewRepository(dbPool *pgxpool.Pool) Repository {
	return &repository{dbPool}
}

type repository struct {
	pool *pgxpool.Pool
}

func (r *repository) ListLeagues(ctx context.Context) ([]League, error) {
	query := `
		SELECT
		    id,
		    name,
		    number_of_teams,
		    country_id
	    FROM leagues
		ORDER BY name
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "errpr fetching discover category details from database")
	}

	var leagues []League
	for rows.Next() {
		league := League{}
		if err := rows.Scan(
			&league.ID,
			&league.Name,
			&league.NumberOfTeams,
			&league.Country,
		); err != nil {
			return nil, errors.Wrap(err, "error scanning row from database")
		}
		leagues = append(leagues, league)
	}
	return leagues, nil
}

func (r *repository) GetLeague(ctx context.Context, id string) (League, error) {
	query := `
	    SELECT
		    id,
		    name,
		    number_of_teams,
		    country_id
	    FROM leagues
	    WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)
	var league League
	if err := row.Scan(
		&league.ID,
		&league.Name,
		&league.NumberOfTeams,
		&league.Country,
	); err != nil {
		return league, errors.Wrapf(err, "error getting league with id %s", id)
	}
	return league, nil
}
