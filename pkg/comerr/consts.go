package comerr

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidData = fmt.Errorf("invalid data")
)

func UnwrapFirst(err error) error {
	for {
		wrappedErr := errors.Unwrap(err)
		if wrappedErr == nil {
			break
		}
		err = wrappedErr
	}
	return err
}
