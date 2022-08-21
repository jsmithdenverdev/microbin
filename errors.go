package microbin

import (
	"fmt"
)

type ErrorPasteExpired struct {
	ID int
}

func (e ErrorPasteExpired) Error() string {
	return fmt.Sprintf("paste %d expired", e.ID)
}

type ErrorUnrecognizedExpiration struct {
	Expiration Expiration
}

func (e ErrorUnrecognizedExpiration) Error() string {
	if e.Expiration == "" {
		return ("unrecognized expiration: none provided")
	}
	return fmt.Sprintf("unrecognized expiration: %s", e.Expiration)
}
