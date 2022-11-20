package leagues

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type listLeaguesRequest struct{}

type listLeaguesResponse struct {
	Leagues []League `json:"leagues"`
}

type getLeagueRequest struct {
	ID string
}

type getLeagueResponse struct {
	League League `json:"league"`
}

func MakeListLeaguesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(listLeaguesRequest)
		leagues, err := svc.ListLeagues(ctx)
		return listLeaguesResponse{Leagues: leagues}, err
	}
}

func MakeGetLeagueEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getLeagueRequest)
		league, err := svc.GetLeague(ctx, req.ID)
		return getLeagueResponse{League: league}, err
	}
}
