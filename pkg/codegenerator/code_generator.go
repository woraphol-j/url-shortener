package codegenerator

import (
	"github.com/teris-io/shortid"
)

//TODO: Find a proper way to generate this value
const (
	workerID = 1
	seedID   = 1234
)

//go:generate mockgen -destination=./code_generator_mock.go -package=codegenerator github.com/woraphol-j/url-shortener/pkg/codegenerator CodeGenerator

// CodeGenerator is the short url code generator. It contains
// Generate method used to generate a short code
type CodeGenerator interface {
	Generate() (string, error)
}

type shortIDCodeGenerator struct {
	sid *shortid.Shortid
}

// NewCodeGenerator is a constructor to create a new code generator
func NewCodeGenerator() CodeGenerator {
	sid, err := shortid.New(workerID, shortid.DEFAULT_ABC, seedID)
	if err != nil {
		panic("Cannot build ID generator")
	}

	return &shortIDCodeGenerator{
		sid,
	}
}

// Generate short code randomly
func (sicg shortIDCodeGenerator) Generate() (string, error) {
	return sicg.sid.Generate()
}
