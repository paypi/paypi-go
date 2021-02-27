package paypi

import (
	"errors"
)

var (
	ErrInvalidToken       error = errors.New("Token is invalid or disallowed")
	ErrUnableToMakeCharge error = errors.New("Unable to make charge")
)
