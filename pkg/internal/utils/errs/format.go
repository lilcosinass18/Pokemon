package errs

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

// Format panics if by the time it's called err is nil.
// However, *err can be nil, in which case it is not changed.
//
// Format is mainly meant to be used in defer (see examples)
func Format(err *error, taskPattern string, args ...any) {
	if *err != nil {
		task := fmt.Sprintf(taskPattern, args...)
		*err = fmt.Errorf("failed to %s: %w", task, *err)
	}
}

func Echo(err error) error {
	return echo.NewHTTPError(GetCode(err), err.Error())
}

// FormatEcho panics if by the time it's called err is nil
//
// FormatEcho is mainly meant to be used in defer (see examples)
func FormatEcho(err *error) {
	if *err != nil {
		*err = Echo(*err)
	}
}
