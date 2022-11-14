package players

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type listPlayersRequest struct{}

type listPlayersResponse struct {
	Players []Player `json:"players"`
}

type getPlayerRequest struct {
	ID string
}

type getPlayerResponse struct {
	Player Player `json:"player"`
}

func MakeListPlayersEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(listPlayersRequest)
		players, err := svc.ListPlayers(ctx)
		return listPlayersResponse{Players: players}, err
	}
}

func MakeGetPlayerEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getPlayerRequest)
		player, err := svc.GetPlayer(ctx, req.ID)
		return getPlayerResponse{Player: player}, err
	}
}
