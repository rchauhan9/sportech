package stadiums

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

func MakeHandler(listStadiumsEndpoint endpoint.Endpoint, getStadiumEndpoint endpoint.Endpoint) http.Handler {
	r := mux.NewRouter()

	listStadiumsHandler := kithttp.NewServer(
		listStadiumsEndpoint,
		decodeListStadiumsRequest,
		encodeListStadiumsResponse,
	)

	getStadiumHandler := kithttp.NewServer(
		getStadiumEndpoint,
		decodeGetStadiumRequest,
		encodeGetStadiumResponse,
	)

	r.Handle("/stadiums/{id}", getStadiumHandler).Methods("GET")
	r.Handle("/stadiums/", listStadiumsHandler).Methods("GET")

	return r
}

func decodeListStadiumsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	_ = mux.Vars(r)
	return listStadiumsRequest{}, nil
}

func encodeListStadiumsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

func decodeGetStadiumRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("bad route")
	}
	return getStadiumRequest{ID: id}, nil
}

func encodeGetStadiumResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	//case erro:
	//	w.WriteHeader(http.StatusNotFound)
	//case ErrInvalidArgument:
	//	w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
