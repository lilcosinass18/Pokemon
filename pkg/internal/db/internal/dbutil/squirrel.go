package dbutil

import (
	"github.com/Masterminds/squirrel"
)

func Builder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
