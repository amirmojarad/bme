package errorext

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type Error struct {
	Err        error
	Key        ErrorKey
	StatusCode int
}

func (e Error) Error() string {
	return e.Err.Error()
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) StackTrace() string {
	return fmt.Sprintf("%+v", e.Err)
}

func IsNotFound(err error) bool {
	if targetErr := As(err); targetErr != nil {
		return targetErr.Key == ErrNotFound
	}

	return false
}

func Is(err error, key ErrorKey) bool {
	if targetErr := As(err); targetErr != nil {
		return targetErr.Key == key
	}

	return false
}

func As(err error) *Error {
	if err == nil {
		return nil
	}

	errorextErr := &Error{}
	errors.As(err, &errorextErr)

	return errorextErr
}

func newError(err error, key ErrorKey, statusCode int) *Error {
	if err == nil {
		err = errors.New(key.GetMessage().En)
	}

	return &Error{
		Err:        errors.WithStack(err),
		Key:        key,
		StatusCode: statusCode,
	}
}

func New(err error, key ErrorKey) *Error {
	return newError(err, key, http.StatusInternalServerError)
}

func NewNotFound(err error, key ErrorKey) *Error {
	return newError(err, key, http.StatusNotFound)
}

func NewValidation(err error, key ErrorKey) *Error {
	return newError(err, key, http.StatusBadRequest)
}

func NewForbidden(err error, key ErrorKey) *Error {
	return newError(err, key, http.StatusForbidden)
}

func NewAuth(err error, key ErrorKey) *Error {
	return newError(err, key, http.StatusUnauthorized)
}

func NewNotImplemented(err error) *Error {
	return newError(err, ErrNotImplemented, http.StatusNotImplemented)
}

func NewBadRequest(err error, key ErrorKey) *Error {
	return newError(err, key, http.StatusBadRequest)
}
