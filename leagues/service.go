package leagues

import (
	"context"
)

type Service interface {
	ListLeagues(ctx context.Context) ([]League, error)
	GetLeague(ctx context.Context, id string) (League, error)
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

type service struct {
	repository Repository
}

func (s *service) ListLeagues(ctx context.Context) ([]League, error) {
	return s.repository.ListLeagues(ctx)
}

func (s *service) GetLeague(ctx context.Context, id string) (League, error) {
	return s.repository.GetLeague(ctx, id)
}
