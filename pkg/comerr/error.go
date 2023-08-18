package comerr

import (
	"errors"
	"fmt"
)

func WrapError(err error, msg string) error {
	if err == nil {
		return errors.New(msg)
	} else {
		return fmt.Errorf("%s: %w", msg, err)
	}
}

func UnwrapError(err error) error {
	return errors.Unwrap(err)
}
