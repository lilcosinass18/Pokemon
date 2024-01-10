package dbutil

import (
	"database/sql"
	"errors"
	"net/http"

	"pokemon-rest-api/pkg/internal/utils/errs"
)

func CheckNotFound(err error, messageNotFound string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return errs.New(http.StatusNotFound, messageNotFound)
	}

	return err
}
