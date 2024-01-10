package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"pokemon-rest-api/pkg/internal/utils/errs"
)

const (
	createPokemons = `
CREATE TABLE if NOT EXISTS pokemons (
    id              INTEGER PRIMARY KEY,
    "name"          VARCHAR(256) NOT NULL,
    height          INTEGER      NOT NULL,
    weight          INTEGER      NOT NULL,
    "order"         INTEGER      NOT NULL,
    base_experience INTEGER      NOT NULL,
    is_default      BOOLEAN      NOT NULL
)
`
	setTransactionSerializable     = "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE"
	setTransactionRepeatableReadRO = "SET TRANSACTION ISOLATION LEVEL REPEATABLE READ, READ ONLY"
)

type Storage interface {
	BeginRW(context.Context) (*RW, error)
	BeginRO(context.Context) (*RO, error)
	CreateSchema(context.Context) error
}

type storage struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewStorage(connString string, logger *zap.Logger) (_ Storage, err error) {
	defer errs.Format(&err, "initialize database")

	cfg, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}

	cfg.Tracer = newPGXTracer(logger)

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	db := sqlx.NewDb(stdlib.OpenDB(*cfg), "pgx")
	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("unable to connect to db: %w", err)
	}

	return &storage{
		db:     db,
		logger: logger,
	}, nil
}

func (m *storage) BeginRW(ctx context.Context) (_ *RW, err error) {
	defer errs.Format(&err, "begin a read-write transaction")

	transaction, err := m.db.Beginx()
	if err != nil {
		return nil, err
	}

	if _, err = transaction.ExecContext(ctx, setTransactionSerializable); err != nil {
		_ = transaction.Rollback()
		return nil, err
	}

	return newRW(transaction, m.logger), nil
}

func (m *storage) BeginRO(ctx context.Context) (_ *RO, err error) {
	defer errs.Format(&err, "begin a read-only transaction")

	transaction, err := m.db.Beginx()
	if err != nil {
		return nil, err
	}

	if _, err = transaction.ExecContext(ctx, setTransactionRepeatableReadRO); err != nil {
		_ = transaction.Rollback()
		return nil, err
	}

	return newRO(transaction, m.logger), nil
}

func finish(tx *sqlx.Tx, err *error) {
	if r := recover(); r != nil {
		*err = errors.Join(*err, fmt.Errorf("panic: %v", r))
	}

	if *err == nil {
		*err = errors.Join(*err, tx.Commit())
	}

	if *err != nil {
		_ = tx.Rollback()
	}
}

func Finish(holder txHolder, err *error) {
	finish(holder.getTx(), err)
}

func (m *storage) CreateSchema(ctx context.Context) (err error) {
	defer errs.Format(&err, "create pokemons db schema")

	transaction, err := m.db.Beginx()
	if err != nil {
		return err
	}
	defer finish(transaction, &err)

	if _, err = transaction.ExecContext(ctx, createPokemons); err != nil {
		return err
	}

	return nil
}
