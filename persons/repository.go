package persons

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Repository interface {
	ListPersons(ctx context.Context) ([]Person, error)
	GetPerson(ctx context.Context, id string) (Person, error)
}

func NewRepository(dbPool *pgxpool.Pool) Repository {
	return &repository{dbPool}
}

type repository struct {
	pool *pgxpool.Pool
}

func (r *repository) ListPersons(ctx context.Context) ([]Person, error) {
	query := `
		SELECT
		    id,
		    first_name,
		    last_name,
		    date_of_birth,
		    country_id
	    FROM persons
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching discover category details from database")
	}

	var persons []Person
	for rows.Next() {
		person := Person{}
		if err := rows.Scan(
			&person.ID,
			&person.FirstName,
			&person.LastName,
			&person.DateOfBirth,
			&person.Nationality,
		); err != nil {
			return nil, errors.Wrap(err, "error scanning row from database")
		}
		persons = append(persons, person)
	}
	return persons, nil
}

func (r *repository) GetPerson(ctx context.Context, id string) (Person, error) {
	query := `
	    SELECT
		    id,
		    first_name,
		    last_name,
		    date_of_birth,
		    country_id
	    FROM persons
	`
	row := r.pool.QueryRow(ctx, query, id)
	var person Person
	if err := row.Scan(
		&person.ID,
		&person.FirstName,
		&person.LastName,
		&person.DateOfBirth,
		&person.Nationality,
	); err != nil {
		return person, errors.Wrapf(err, "error getting person with id %s", id)
	}
	return person, nil
}
