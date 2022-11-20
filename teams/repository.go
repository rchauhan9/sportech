package teams

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Repository interface {
	ListTeams(ctx context.Context) ([]Team, error)
	GetTeam(ctx context.Context, id string) (Team, error)
}

func NewRepository(dbPool *pgxpool.Pool) Repository {
	return &repository{dbPool}
}

type repository struct {
	pool *pgxpool.Pool
}

func (r *repository) ListTeams(ctx context.Context) ([]Team, error) {
	query := `
		SELECT
		    id,
		    full_name,
		    medium_name,
		    acronym,
		    nickname,
		    year_founded,
		    city,
		    country_id,
		    stadium_id,
		    league_id
	    FROM teams
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching discover category details from database")
	}

	var teams []Team
	for rows.Next() {
		team := Team{}
		if err := rows.Scan(
			&team.ID,
			&team.FullName,
			&team.MediumName,
			&team.Acronym,
			&team.Nickname,
			&team.YearFounded,
			&team.City,
			&team.Country,
			&team.Stadium,
			&team.League,
		); err != nil {
			return nil, errors.Wrap(err, "error scanning row from database")
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func (r *repository) GetTeam(ctx context.Context, id string) (Team, error) {
	query := `
	    SELECT
		    id,
		    full_name,
		    medium_name,
		    acronym,
		    nickname,
		    year_founded,
		    city,
		    country_id,
		    stadium_id,
		    league_id
	    FROM teams
	    WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)
	var team Team
	if err := row.Scan(
		&team.ID,
		&team.FullName,
		&team.MediumName,
		&team.Acronym,
		&team.Nickname,
		&team.YearFounded,
		&team.City,
		&team.Country,
		&team.Stadium,
		&team.League,
	); err != nil {
		return team, errors.Wrapf(err, "error getting team with id %s", id)
	}
	return team, nil
}
