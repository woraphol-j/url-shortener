package urlshortener

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// MakeHandler returns a handler for the url shortening service.
func MakeHandler(service Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	generateShortURLHandler := kithttp.NewServer(
		makeGenerateShortURLEndpoint(service),
		decodeGenerateShortURLRequest,
		encodeResponse,
		opts...,
	)

	getOriginalURLHandler := kithttp.NewServer(
		makeGetOriginalURLEndpoint(service),
		decodeGetOriginalURLRequest,
		redirectResponse,
		opts...,
	)

	router := mux.NewRouter()
	router.Handle("/shorturls", generateShortURLHandler).Methods("POST")
	router.Handle("/{code}", getOriginalURLHandler).Methods("GET")
	return router
}

var errBadRoute = errors.New("bad route")

func decodeGenerateShortURLRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return generateShortURLRequest{
		URL: body.URL,
	}, nil
}

func decodeGetOriginalURLRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	code, ok := vars["code"]
	if !ok {
		return nil, errBadRoute
	}
	return getOriginalURLRequest{Code: code}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func redirectResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	result := response.(getOriginalURLResponse)
	if result.OriginalURL == "" {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}
	w.Header().Set("Location", result.OriginalURL)
	w.WriteHeader(http.StatusMovedPermanently)
	return nil
}

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
