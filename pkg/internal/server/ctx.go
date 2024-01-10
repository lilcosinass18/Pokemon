package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ContextWrapper struct {
	echo.Context
}

func WrapContext(ctx echo.Context) ContextWrapper {
	return ContextWrapper{Context: ctx}
}

func (m ContextWrapper) ParamInt(name string) (int, error) {
	id := m.Param(name)
	num, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid '%s' parameter: %s", name, err.Error()))
	}

	return int(num), nil
}
