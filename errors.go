package main

import (
	"fmt"
)

const (
	errorInternalServer = "Internal Server Error"
)

type (
	errorPasteExpired struct {
		ID int
	}

	errorUnrecognizedExpiration struct {
		Expiration string
	}
)

func (e errorPasteExpired) Error() string {
	return fmt.Sprintf("paste %d expired", e.ID)
}

func (e errorUnrecognizedExpiration) Error() string {
	if e.Expiration == "" {

		return ("unrecognized expiration: none provided")
	}

	return fmt.Sprintf("unrecognized expiration: %s", e.Expiration)
}
