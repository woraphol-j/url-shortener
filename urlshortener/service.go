// Package urlshortener provides the use-case of booking a cargo. Used by views
// facing an administrator.
package urlshortener

import (
	"errors"

	"github.com/lytics/base62"
	"github.com/woraphol-j/url-shortener/pkg/mongo"
)

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("invalid argument")

// Service is the interface that provides url shortening methods.
type Service interface {
	// GenerateShortURL generate and save short url for use later
	GenerateShortURL(url string) (string, error)

	// GetOriginalURL returns a read model of a cargo.
	GetOriginalURL(code string) (string, error)
}

type service struct {
	urlDAO *mongo.DAO
}

// GenerateShortURL is
func (s *service) GenerateShortURL(url string) (string, error) {
	code := base62.StdEncoding.EncodeToString([]byte(url))
	s.urlDAO.Save(&mongo.ShortURL{
		Code: code,
		URL:  url,
	})
	return code, nil
}

// GetOriginalURL get the original url based on code
func (s *service) GetOriginalURL(code string) (string, error) {
	shortURL, err := s.urlDAO.Get(code)
	if err != nil {
		return "", err
	}
	return shortURL.URL, nil
}

// NewService creates a url shortening service. It requires DAO object
func NewService(urlDAO *mongo.DAO) Service {
	return &service{urlDAO}
}
