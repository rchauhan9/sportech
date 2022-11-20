package persons

import (
	"context"
)

type Service interface {
	ListPersons(ctx context.Context) ([]Person, error)
	GetPerson(ctx context.Context, id string) (Person, error)
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

type service struct {
	repository Repository
}

func (s *service) ListPersons(ctx context.Context) ([]Person, error) {
	return s.repository.ListPersons(ctx)
}

func (s *service) GetPerson(ctx context.Context, id string) (Person, error) {
	return s.repository.GetPerson(ctx, id)
}
