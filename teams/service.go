package teams

import (
	"context"
)

type Service interface {
	ListTeams(ctx context.Context) ([]Team, error)
	GetTeam(ctx context.Context, id string) (Team, error)
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

type service struct {
	repository Repository
}

func (s *service) ListTeams(ctx context.Context) ([]Team, error) {
	return s.repository.ListTeams(ctx)
}

func (s *service) GetTeam(ctx context.Context, id string) (Team, error) {
	return s.repository.GetTeam(ctx, id)
}
