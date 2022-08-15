package team

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type listTeamsRequest struct{}

type listTeamsResponse struct {
	Teams []Team `json:"teams"`
}

type getTeamRequest struct {
	ID string
}

type getTeamResponse struct {
	Team Team `json:"team"`
}

func MakeListTeamsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(listTeamsRequest)
		teams, err := svc.ListTeams(ctx)
		return listTeamsResponse{Teams: teams}, err
	}
}

func MakeGetTeamEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getTeamRequest)
		team, err := svc.GetTeam(ctx, req.ID)
		return getTeamResponse{Team: team}, err
	}
}
