package comerr

import (
	"errors"
	"fmt"
	"strings"
)

func WrapMessage(err error, msg string) error {
	if err == nil {
		return errors.New(msg)
	} else {
		return fmt.Errorf("%s: %w", msg, err)
	}
}

func MergeError(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	var errStrs []string
	for _, err := range errs {
		errStrs = append(errStrs, err.Error())
	}
	return errors.New(strings.Join(errStrs, "; "))
}

func UnwrapError(err error) error {
	return errors.Unwrap(err)
}

func ToError(v any) error {
	if v == nil {
		return nil
	}
	err, ok := v.(error)
	if !ok {
		err = fmt.Errorf("error: %v", v)
	}
	return err
}
