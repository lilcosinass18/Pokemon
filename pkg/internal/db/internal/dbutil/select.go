package dbutil

import (
	"context"

	"github.com/Masterminds/squirrel"
)

type Selecter interface {
	SelectContext(context.Context, any, string, ...any) error
}

type Select[T any] struct {
	Query string
	Args  []any
}

func NewSelect[T any](sq squirrel.Sqlizer) (*Select[T], error) {
	query, args, err := sq.ToSql()
	if err != nil {
		return nil, err
	}

	return &Select[T]{
		Query: query,
		Args:  args,
	}, nil
}

func (m *Select[T]) Select(ctx context.Context, tx Selecter) ([]T, error) {
	dst := make([]T, 0)

	if err := tx.SelectContext(ctx, &dst, m.Query, m.Args...); err != nil {
		return nil, err
	}

	return dst, nil
}
