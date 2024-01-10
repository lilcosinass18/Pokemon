package dbutil

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
)

type Execer interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

type Exec struct {
	Query string
	Args  []any
}

func NewExec(sq squirrel.Sqlizer) (*Exec, error) {
	query, args, err := sq.ToSql()
	if err != nil {
		return nil, err
	}

	return &Exec{
		Query: query,
		Args:  args,
	}, nil
}

func (m *Exec) Exec(ctx context.Context, tx Execer) (sql.Result, error) {
	return tx.ExecContext(ctx, m.Query, m.Args...)
}
