package db

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"pokemon-rest-api/pkg/internal/db/internal/pokemons"
)

type txHolder interface {
	getTx() *sqlx.Tx
}

type RO struct {
	tx     *sqlx.Tx
	logger *zap.Logger

	Pokemons *pokemons.RO
}

func newRO(tx *sqlx.Tx, logger *zap.Logger) *RO {
	return &RO{
		tx:     tx,
		logger: logger,

		Pokemons: pokemons.NewRO(tx, logger),
	}
}

func (m *RO) getTx() *sqlx.Tx {
	return m.tx
}

// ro is used to encapsulate the embedded field in RW by making it unexported
type ro = RO

type RW struct {
	*ro

	tx     *sqlx.Tx
	logger *zap.Logger

	Pokemons *pokemons.RW
}

func newRW(tx *sqlx.Tx, logger *zap.Logger) *RW {
	return &RW{
		ro: newRO(tx, logger),

		tx:     tx,
		logger: logger,

		Pokemons: pokemons.NewRW(tx, logger),
	}
}

func (m *RW) getTx() *sqlx.Tx {
	return m.tx
}

func (m *RW) RO() *RO {
	return m.ro
}
