package urlshortener

import (
	"errors"

	cg "github.com/woraphol-j/url-shortener/pkg/codegenerator"
	"github.com/woraphol-j/url-shortener/pkg/repository"
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
	urlDAO        repository.Repository
	codeGenerator cg.CodeGenerator
}

// NewService creates a url shortening service.
// It requires DAO object and the code generator.
func NewService(urlDAO repository.Repository, codeGenerator cg.CodeGenerator) Service {
	return &service{
		urlDAO,
		codeGenerator,
	}
}

// GenerateShortURL is
func (s *service) GenerateShortURL(url string) (string, error) {
	code, err := s.codeGenerator.Generate()
	if err != nil {
		return "", err
	}
	s.urlDAO.Save(&repository.ShortURL{
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
	if shortURL.NotFound {
		return "", nil
	}
	return shortURL.URL, nil
}
