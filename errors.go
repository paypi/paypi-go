package paypi

import (
	"errors"
)

var (
	ErrInvalidToken       error = errors.New("paypi: subscriber key is invalid, deny this request")
	ErrUnableToMakeCharge error = errors.New("paypi: unable to make charge, deny this request")
)
