package pokemons

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"go.uber.org/zap"

	"pokemon-rest-api/listing"
	"pokemon-rest-api/pkg/internal/db/internal/dbnames"
	"pokemon-rest-api/pkg/internal/db/internal/dbutil"
	"pokemon-rest-api/pkg/internal/utils/errs"
	"pokemon-rest-api/pkg/internal/utils/slices"
)

type prepareRO struct{}

type RO struct {
	tx      dbutil.DB
	prepare prepareRO
	logger  *zap.Logger
}

func NewRO(tx dbutil.DB, logger *zap.Logger) *RO {
	return &RO{tx: tx, prepare: struct{}{}, logger: logger}
}

//----------------------------------------------------------------------------------------------------------------------

func (m prepareRO) fetch() (_ *dbutil.Select[*dbPokemon], err error) {
	defer errs.Format(&err, "build a query to fetch pokemons")

	pokemons := dbnames.Tables.Pokemons

	return dbutil.NewSelect[*dbPokemon](
		dbutil.Builder().
			Select(pokemons.AllFields()).
			From(pokemons.TName),
	)
}

func (m *RO) Fetch(ctx context.Context) (_ []*listing.Pokemon, err error) {
	defer errs.Format(&err, "fetch pokemons from db")

	request, err := m.prepare.fetch()
	if err != nil {
		return nil, err
	}

	pokemons, err := request.Select(ctx, m.tx)
	if err != nil {
		// TODO: not found is OK
		return nil, err
	}

	return slices.Map(pokemons, (*dbPokemon).Deserialize), nil
}

//----------------------------------------------------------------------------------------------------------------------

func (m prepareRO) get(id int) (_ *dbutil.Get[dbPokemon], err error) {
	defer errs.Format(&err, "build a query to get a pokemon")

	pokemons := dbnames.Tables.Pokemons

	return dbutil.NewGet[dbPokemon](
		dbutil.Builder().
			Select(pokemons.AllFields()).
			From(pokemons.TName).
			Where(squirrel.Eq{pokemons.ID: id}),
	)
}

func (m *RO) Get(ctx context.Context, id int) (_ *listing.Pokemon, err error) {
	defer errs.Format(&err, "get a pokemon from db")

	req, err := m.prepare.get(id)
	if err != nil {
		return nil, err
	}

	attribute, err := req.Get(ctx, m.tx)
	if err != nil {
		return nil, dbutil.CheckNotFound(err, fmt.Sprintf("pokemon with id %d was not found", id))
	}

	return attribute.Deserialize(), nil
}
