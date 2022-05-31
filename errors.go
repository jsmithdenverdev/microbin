package microbin

import (
	"errors"
	"fmt"
)

type errFunc func() string

func (e errFunc) Error() string {
	return e()
}

var (
	ErrInvalidExpiration = errors.New("invalid expiration value")
	ErrExpiredPaste      = func(id int) error {
		return errFunc(func() string {
			return fmt.Sprintf("paste %d is expired", id)
		})
	}
)
