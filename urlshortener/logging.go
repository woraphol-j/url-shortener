package urlshortener

import (
	"time"

	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, service Service) Service {
	return &loggingService{logger, service}
}

func (s *loggingService) GenerateShortURL(url string) (shortURL string, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "generateShortURL",
			"url", url,
			"shortUrl", shortURL,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GenerateShortURL(url)
}

func (s *loggingService) GetOriginalURL(code string) (originalURL string, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "getOriginalURL",
			"code", code,
			"originalUrl", originalURL,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetOriginalURL(code)
}
