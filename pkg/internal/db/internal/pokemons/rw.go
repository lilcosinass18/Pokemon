package pokemons

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"go.uber.org/zap"

	"pokemon-rest-api/listing"
	"pokemon-rest-api/pkg/internal/db/internal/dbnames"
	"pokemon-rest-api/pkg/internal/db/internal/dbutil"
	"pokemon-rest-api/pkg/internal/utils/errs"
	"pokemon-rest-api/pkg/internal/utils/slices"
)

// ro is used to encapsulate the embedded field in RW by making it unexported
type ro = RO

type prepare struct{}

type RW struct {
	*ro
	prepare prepare
}

func NewRW(tx dbutil.DB, logger *zap.Logger) *RW {
	return &RW{ro: NewRO(tx, logger), prepare: struct{}{}}
}

//----------------------------------------------------------------------------------------------------------------------

func (m prepare) insert(pokemons ...*dbPokemon) (_ *dbutil.Select[int], err error) {
	defer errs.Format(&err, "build a query to insert pokemons")

	pokemonsTable := dbnames.Tables.Pokemons

	sq := dbutil.Builder().
		Insert(pokemonsTable.TName).
		Columns(
			pokemonsTable.ID,
			pokemonsTable.Name,
			pokemonsTable.Weight,
			pokemonsTable.Height,
			pokemonsTable.Order,
			pokemonsTable.BaseExperience,
			pokemonsTable.IsDefault,
		)

	for _, pokemon := range pokemons {
		sq = sq.Values(
			pokemon.ID,
			pokemon.Name,
			pokemon.Weight,
			pokemon.Height,
			pokemon.Order,
			pokemon.BaseExperience,
			pokemon.IsDefault,
		)
	}

	sq = sq.Suffix(fmt.Sprintf("RETURNING %s", pokemonsTable.ID))

	return dbutil.NewSelect[int](sq)
}

func (m *RW) Insert(ctx context.Context, pokemons ...*listing.Pokemon) (_ []int, err error) {
	if len(pokemons) == 0 {
		return nil, nil
	}

	defer errs.Format(&err, "fetch pokemons from db")

	req, err := m.prepare.insert(slices.Map(pokemons, newDBPokemon)...)
	if err != nil {
		return nil, err
	}

	result, err := req.Select(ctx, m.tx)
	if err != nil {
		return nil, errs.New(http.StatusConflict, "pokemon with such id exists")
	}

	return result, nil
}

//----------------------------------------------------------------------------------------------------------------------

func (m prepare) delete(id int) (_ *dbutil.Get[dbPokemon], err error) {
	defer errs.Format(&err, "build a query to delete a pokemon")

	pokemons := dbnames.Tables.Pokemons

	return dbutil.NewGet[dbPokemon](
		dbutil.Builder().
			Delete(pokemons.TName).
			Where(squirrel.Eq{pokemons.ID: id}).
			Suffix(fmt.Sprintf("RETURNING %s", pokemons.AllFields())),
	)
}

func (m *RW) Delete(ctx context.Context, id int) (_ *listing.Pokemon, err error) {
	defer errs.Format(&err, "delete a pokemon from db")

	req, err := m.prepare.delete(id)
	if err != nil {
		return nil, err
	}

	pokemon, err := req.Get(ctx, m.tx)
	if err != nil {
		return nil, dbutil.CheckNotFound(err, fmt.Sprintf("pokemon with id %d was not found", id))
	}

	return pokemon.Deserialize(), nil
}
