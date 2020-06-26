package errors

import "errors"

func New(text string) error {
	return errors.New(text)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, &target)
}
