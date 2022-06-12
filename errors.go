package microbin

import (
	"errors"
)

var (
	ErrInvalidExpiration = errors.New("invalid expiration value")
	ErrExpiredPaste      = errors.New("paste expired")
	ErrPasteNotFound     = errors.New("not found")
)
