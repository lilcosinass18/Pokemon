package errs

import (
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	code int
	err  error
}

func New(code int, msg string) error {
	return Error{
		code: code,
		err:  errors.New(msg),
	}
}

func Newf(code int, msg string, args ...any) error {
	return Error{
		code: code,
		err:  fmt.Errorf(msg, args...),
	}
}

func WithCode(err error, code int) error {
	return Error{
		code: code,
		err:  err,
	}
}

func (m Error) Error() string {
	return m.err.Error()
}

func (m Error) Unwrap() error {
	return m.err
}

func (m Error) Is(target error) bool {
	//nolint:errorlint // per errors package documentation, Is method should perform a shallow check
	statusErr, ok := target.(Error)
	if !ok {
		return false
	}

	return m.Error() == statusErr.Error()
}

func GetCode(err error) int {
	var statusErr Error

	if errors.As(err, &statusErr) {
		return statusErr.code
	}

	return http.StatusInternalServerError
}
