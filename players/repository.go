package players

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Repository interface {
	ListPlayers(ctx context.Context) ([]PlayerDB, error)
	GetPlayer(ctx context.Context, id string) (PlayerDB, error)
}

func NewRepository(dbPool *pgxpool.Pool) Repository {
	return &repository{dbPool}
}

type repository struct {
	pool *pgxpool.Pool
}

func (r *repository) ListPlayers(ctx context.Context) ([]PlayerDB, error) {
	query := `
		SELECT
		    id,
		    person_id,
		    CONCAT('/teams/', team),
		    squad_number,
		    general_position,
		    specific_position,
		    started,
		    ended
	    FROM team_players
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching discover category details from database")
	}

	var players []PlayerDB
	for rows.Next() {
		player := PlayerDB{}
		if err := rows.Scan(
			&player.ID,
			&player.Team,
			&player.SquadNumber,
			&player.GeneralPosition,
			&player.SpecificPosition,
			&player.Started,
			&player.Ended,
		); err != nil {
			return nil, errors.Wrap(err, "error scanning row from database")
		}
		players = append(players, player)
	}
	return players, nil
}

func (r *repository) GetPlayer(ctx context.Context, id string) (PlayerDB, error) {
	query := `
	    SELECT
		    id,
		    person_id,
		    CONCAT('/teams/', team),
		    squad_number,
		    general_position,
		    specific_position,
		    started,
		    ended
	    FROM team_players
	    WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)
	var player PlayerDB
	if err := row.Scan(
		&player.ID,
		&player.Team,
		&player.SquadNumber,
		&player.GeneralPosition,
		&player.SpecificPosition,
		&player.Started,
		&player.Ended,
	); err != nil {
		return player, errors.Wrapf(err, "error getting player with id %s", id)
	}
	return player, nil
}
