package leagues

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

func MakeHandler(listLeaguesEndpoint endpoint.Endpoint, getLeagueEndpoint endpoint.Endpoint) http.Handler {
	r := mux.NewRouter()

	listLeaguesHandler := kithttp.NewServer(
		listLeaguesEndpoint,
		decodeListLeaguesRequest,
		encodeListLeaguesResponse,
	)

	getLeagueHandler := kithttp.NewServer(
		getLeagueEndpoint,
		decodeGetLeagueRequest,
		encodeGetLeagueResponse,
	)

	r.Handle("/leagues/{id}", getLeagueHandler).Methods("GET")
	r.Handle("/leagues/", listLeaguesHandler).Methods("GET")

	return r
}

func decodeListLeaguesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	_ = mux.Vars(r)
	return listLeaguesRequest{}, nil
}

func encodeListLeaguesResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

func decodeGetLeagueRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("bad route")
	}
	return getLeagueRequest{ID: id}, nil
}

func encodeGetLeagueResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
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
