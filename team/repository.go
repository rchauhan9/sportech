package team

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type Repository interface {
	ListTeams(ctx context.Context) ([]Team, error)
	GetTeam(ctx context.Context, id string) (Team, error)
}

func NewRepository() Repository {
	return &repository{}
}

type repository struct{}

var teams = []Team{
	{uuid.New().String(), "Arsenal Football Club"},
	{uuid.New().String(), "Manchester United Football Club"},
	{uuid.New().String(), "Manchester City Football Club"},
	{uuid.New().String(), "Liverpool Football Club"},
	{uuid.New().String(), "Chelsea Football Club"},
	{uuid.New().String(), "Tottenham Hotspur Football Club"},
}

func (r *repository) ListTeams(ctx context.Context) ([]Team, error) {
	return teams, nil
}

func (r *repository) GetTeam(ctx context.Context, id string) (Team, error) {
	teams := funk.Filter(teams, func(team Team) bool {
		return team.ID == id
	}).([]Team)
	if len(teams) < 1 {
		return Team{}, errors.Errorf("error no team with id %s", id)
	}
	return teams[0], nil
}
