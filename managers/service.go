package managers

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rchauhan9/sportech/persons"
	"github.com/samber/lo"
)

type Service interface {
	ListManagers(ctx context.Context) ([]Manager, error)
	GetManager(ctx context.Context, id string) (Manager, error)
}

func NewService(repository Repository, personsService persons.Service) Service {
	return &service{repository: repository, personsService: personsService}
}

type service struct {
	repository     Repository
	personsService persons.Service
}

func (s *service) ListManagers(ctx context.Context) ([]Manager, error) {
	managersDB, err := s.repository.ListManagers(ctx)
	if err != nil {
		return nil, err
	}

	people, err := s.personsService.ListPersons(ctx)

	personsMap := lo.KeyBy[string, persons.Person](people, func(person persons.Person) string {
		return person.ID
	})

	managers := make([]Manager, len(managersDB))
	for i := range managersDB {
		managers[i] = Manager{
			ID:          managersDB[i].ID,
			FirstName:   personsMap[managersDB[i].PersonID].FirstName,
			MiddleNames: personsMap[managersDB[i].PersonID].MiddleNames,
			LastName:    personsMap[managersDB[i].PersonID].LastName,
			DateOfBirth: personsMap[managersDB[i].PersonID].DateOfBirth,
			Nationality: personsMap[managersDB[i].PersonID].Nationality,
			Team:        managersDB[i].TeamID,
			Started:     managersDB[i].Started,
			Ended:       managersDB[i].Ended,
		}
	}

	return managers, nil
}

func (s *service) GetManager(ctx context.Context, id string) (Manager, error) {
	manager, err := s.repository.GetManager(ctx, id)
	if err != nil {
		return Manager{}, err
	}

	person, err := s.personsService.GetPerson(ctx, manager.PersonID)
	if err != nil {
		return Manager{}, errors.Wrapf(err, "error getting manager with id %s", id)
	}

	return Manager{
		ID:          manager.ID,
		FirstName:   person.FirstName,
		MiddleNames: person.MiddleNames,
		LastName:    person.LastName,
		DateOfBirth: person.DateOfBirth,
		Nationality: person.Nationality,
		Team:        manager.TeamID,
		Started:     manager.Started,
		Ended:       manager.Ended,
	}, nil
}
