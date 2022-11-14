package managers

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type listManagersRequest struct{}

type listManagersResponse struct {
	Managers []Manager `json:"managers"`
}

type getManagerRequest struct {
	ID string
}

type getManagerResponse struct {
	Manager Manager `json:"manager"`
}

func MakeListManagersEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(listManagersRequest)
		managers, err := svc.ListManagers(ctx)
		return listManagersResponse{Managers: managers}, err
	}
}

func MakeGetManagerEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getManagerRequest)
		manager, err := svc.GetManager(ctx, req.ID)
		return getManagerResponse{Manager: manager}, err
	}
}
