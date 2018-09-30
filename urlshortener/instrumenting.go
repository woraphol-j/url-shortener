package urlshortener

import (
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) GenerateShortURL(url string) (string, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "generateShortURL").Add(1)
		s.requestLatency.With("method", "generateShortURL").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GenerateShortURL(url)
}

func (s *instrumentingService) GetOriginalURL(code string) (string, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "getOriginalURL").Add(1)
		s.requestLatency.With("method", "getOriginalURL").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetOriginalURL(code)
}
