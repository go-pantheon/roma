package filewriter

import (
	"os"

	"github.com/pkg/errors"
)

type GenService interface {
	Execute() ([]byte, error)
}

func GenFile(filename string, s GenService) error {
	if _, err := os.Stat(filename); err == nil || !os.IsNotExist(err) {
		return errors.Wrapf(err, "%s already exists", filename)
	}

	b, err := s.Execute()
	if err != nil {
		return err
	}

	return WriteFile(filename, b)
}
