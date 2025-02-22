package repo

import (
	"errors"
)

var ErrNotFound = errors.New("not found in repository")

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
