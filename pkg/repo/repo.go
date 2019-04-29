package repo

import (
	"io"
	"os"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
)

// Repo provide basic file operations
type Repo interface {
	Create(reader io.Reader) (ID string, size int64, err error)
	Open(ID string) (io.ReadCloser, error)
}

type repo struct {
	basedir string
}

// NewRepo returns basic Repo's implementation
func NewRepo(basedir string) (Repo, error) {
	if _, err := os.Stat(basedir); err != nil {
		return nil, err
	}
	return &repo{basedir: basedir}, nil
}

// Create method creates local file and puts to it content from provided Reader
func (r *repo) Create(reader io.Reader) (string, int64, error) {
	ID := r.generateID()

	f, err := os.OpenFile(filepath.Join(r.basedir, ID), os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	size, err := io.Copy(f, reader)

	return ID, size, err
}

func (r *repo) generateID() string {
	return uuid.NewV4().String()
}

// Open returns reader of previously created file by ID
func (r *repo) Open(ID string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(r.basedir, ID))
}
