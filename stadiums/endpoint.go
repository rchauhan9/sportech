package stadiums

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type listStadiumsRequest struct{}

type listStadiumsResponse struct {
	Stadiums []Stadium `json:"stadiums"`
}

type getStadiumRequest struct {
	ID string
}

type getStadiumResponse struct {
	Stadium Stadium `json:"stadium"`
}

func MakeListStadiumsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(listStadiumsRequest)
		Stadiums, err := svc.ListStadiums(ctx)
		return listStadiumsResponse{Stadiums: Stadiums}, err
	}
}

func MakeGetStadiumEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getStadiumRequest)
		Stadium, err := svc.GetStadium(ctx, req.ID)
		return getStadiumResponse{Stadium: Stadium}, err
	}
}
