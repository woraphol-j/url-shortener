package urlshortener

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	// "github.com/woraphol-j/url-shortener/urlshortener"
)

// MakeHandler returns a handler for the url shortening service.
func MakeHandler(service Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	generateShortURLHandler := kithttp.NewServer(
		makeGenerateShortURLEndpoint(service),
		decodegenerateShortURLRequest,
		encodeResponse,
		opts...,
	)
	getOriginalURLHandler := kithttp.NewServer(
		makeGetOriginalURLEndpoint(service),
		decodeGetOriginalURLRequest,
		encodeResponse,
		opts...,
	)

	router := mux.NewRouter()

	// router.Handle("/v1/shorturls", generateShortURLHandler).Methods("POST")
	// router.Handle("/v1/shorturls/{id}", getOriginalURLHandler).Methods("GET")
	router.Handle("/shorturls", generateShortURLHandler).Methods("POST")
	router.Handle("/shorturls/{code}", getOriginalURLHandler).Methods("GET")
	// router.Handle("/shorturls/dd", getOriginalURLHandler).Methods("GET")

	return router
}

var errBadRoute = errors.New("bad route")

func decodegenerateShortURLRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		URL string `json:"url"`
		// Destination     string    `json:"destination"`
		// ArrivalDeadline time.Time `json:"arrivalDeadline"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return generateShortURLRequest{
		URL: body.URL,
		// Origin:          location.UNLocode(body.Origin),
		// Destination:     location.UNLocode(body.Destination),
		// ArrivalDeadline: body.ArrivalDeadline,
	}, nil
}

func decodeGetOriginalURLRequest(_ context.Context, r *http.Request) (interface{}, error) {
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

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	// case cargo.ErrUnknown:
	// 	w.WriteHeader(http.StatusNotFound)
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
