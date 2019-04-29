package service

import (
	"io"

	"github.com/kfreiman/imageservice/pkg/processor"
	"github.com/kfreiman/imageservice/pkg/repo"
)

// Service describes service abilities.
// This service is able to modify images (the Processor interface) and
// to save and read files (the Repo, which is not limited to images).
type Service interface {
	repo.Repo
	processor.Processor
}

type service struct {
	repo      repo.Repo
	processor processor.Processor
}

// NewService ..
func NewService(
	repo repo.Repo,
	processor processor.Processor,
) Service {
	return &service{repo: repo, processor: processor}
}

// Create file
func (s *service) Create(reader io.Reader) (ID string, size int64, err error) {
	return s.repo.Create(reader)
}

// Open file
func (s *service) Open(ID string) (io.ReadCloser, error) {
	return s.repo.Open(ID)
}

// Modify image
func (s *service) Modify(income io.Reader, output io.Writer, modifiers processor.Modifiers) error {
	return s.processor.Modify(income, output, modifiers)
}
