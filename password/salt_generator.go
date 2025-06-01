package password

import (
	"crypto/rand"
	"errors"
)

type DefaultSaltGenerator struct {
	length int
}

var errSaltGenerationFailed = errors.New("salt generation failed")

func (r *DefaultSaltGenerator) Generate() ([]byte, error) {
	s := make([]byte, r.length)

	i, err := rand.Read(s)
	if err != nil {
		return nil, err
	}

	if i != r.length {
		return nil, errSaltGenerationFailed
	}

	return s, nil
}

func NewDefaultSaltGenerator(length int) *DefaultSaltGenerator {
	return &DefaultSaltGenerator{
		length: length,
	}
}
