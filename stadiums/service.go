package stadiums

import (
	"context"
)

type Service interface {
	ListStadiums(ctx context.Context) ([]Stadium, error)
	GetStadium(ctx context.Context, id string) (Stadium, error)
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

type service struct {
	repository Repository
}

func (s *service) ListStadiums(ctx context.Context) ([]Stadium, error) {
	return s.repository.ListStadiums(ctx)
}

func (s *service) GetStadium(ctx context.Context, id string) (Stadium, error) {
	return s.repository.GetStadium(ctx, id)
}
