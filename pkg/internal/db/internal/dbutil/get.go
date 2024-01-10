package dbutil

import (
	"context"

	"github.com/Masterminds/squirrel"
)

type Getter interface {
	GetContext(context.Context, any, string, ...any) error
}

type Get[T any] struct {
	Query string
	Args  []any
}

func NewGet[T any](sq squirrel.Sqlizer) (*Get[T], error) {
	query, args, err := sq.ToSql()
	if err != nil {
		return nil, err
	}

	return &Get[T]{
		Query: query,
		Args:  args,
	}, nil
}

func (m *Get[T]) Get(ctx context.Context, tx Getter) (T, error) {
	var dst, mock T

	if err := tx.GetContext(ctx, &dst, m.Query, m.Args...); err != nil {
		return mock, err
	}

	return dst, nil
}
