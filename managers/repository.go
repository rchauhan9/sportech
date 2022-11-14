package managers

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Repository interface {
	ListManagers(ctx context.Context) ([]ManagerDB, error)
	GetManager(ctx context.Context, id string) (ManagerDB, error)
}

func NewRepository(dbPool *pgxpool.Pool) Repository {
	return &repository{dbPool}
}

type repository struct {
	pool *pgxpool.Pool
}

func (r *repository) ListManagers(ctx context.Context) ([]ManagerDB, error) {
	query := `
		SELECT
		    m.id,
		    m.person_id,
		    CONCAT('/teams/', m.team),
		    m.started,
		    m.ended
	    FROM team_managers m
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "errpr fetching discover category details from database")
	}

	var managers []ManagerDB
	for rows.Next() {
		manager := ManagerDB{}
		if err := rows.Scan(
			&manager.ID,
			&manager.Person,
			&manager.Team,
			&manager.Started,
			&manager.Ended,
		); err != nil {
			return nil, errors.Wrap(err, "error scanning row from database")
		}
		managers = append(managers, manager)
	}
	return managers, nil
}

func (r *repository) GetManager(ctx context.Context, id string) (ManagerDB, error) {
	query := `
	    SELECT
		    m.id,
		    m.person_id,
		    CONCAT('/teams/', m.team),
		    m.started,
		    m.ended
	    FROM team_managers m
	    WHERE m.id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)
	var manager ManagerDB
	if err := row.Scan(
		&manager.ID,
		&manager.Person,
		&manager.Team,
		&manager.Started,
		&manager.Ended,
	); err != nil {
		return manager, errors.Wrapf(err, "error getting manager with id %s", id)
	}
	return manager, nil
}
