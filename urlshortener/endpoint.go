package urlshortener

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type generateShortURLRequest struct {
	URL string
	// Origin          location.UNLocode
	// Destination     location.UNLocode
	// ArrivalDeadline time.Time
}

type generateShortURLResponse struct {
	ShortURL string `json:"shortUrl,omitempty"`
	Err      error  `json:"error,omitempty"`
}

func (r generateShortURLResponse) error() error { return r.Err }

func makeGenerateShortURLEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(generateShortURLRequest)
		shortURL, err := s.GenerateShortURL(req.URL)
		return generateShortURLResponse{
			ShortURL: shortURL,
			Err:      err,
		}, nil
	}
}

type getOriginalURLRequest struct {
	Code string
}

type getOriginalURLResponse struct {
	OriginalURL string `json:"cargo,omitempty"`
	Err         error  `json:"error,omitempty"`
}

func (r getOriginalURLResponse) error() error { return r.Err }

func makeGetOriginalURLEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getOriginalURLRequest)
		originalURL, err := s.GetOriginalURL(req.Code)
		return getOriginalURLResponse{
			OriginalURL: originalURL,
			Err:         err,
		}, nil
	}
}
