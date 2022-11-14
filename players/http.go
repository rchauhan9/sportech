package players

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

func MakeHandler(listPlayersEndpoint endpoint.Endpoint, getPlayerEndpoint endpoint.Endpoint) http.Handler {
	r := mux.NewRouter()

	listPlayersHandler := kithttp.NewServer(
		listPlayersEndpoint,
		decodeListPlayersRequest,
		encodeListPlayersResponse,
	)

	getPlayerHandler := kithttp.NewServer(
		getPlayerEndpoint,
		decodeGetPlayerRequest,
		encodeGetPlayerResponse,
	)

	r.Handle("/players/{id}", getPlayerHandler).Methods("GET")
	r.Handle("/players/", listPlayersHandler).Methods("GET")

	return r
}

func decodeListPlayersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	_ = mux.Vars(r)
	return listPlayersRequest{}, nil
}

func encodeListPlayersResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

func decodeGetPlayerRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("bad route")
	}
	return getPlayerRequest{ID: id}, nil
}

func encodeGetPlayerResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
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
