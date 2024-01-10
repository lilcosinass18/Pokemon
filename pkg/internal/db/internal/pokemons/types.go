package pokemons

import (
	"pokemon-rest-api/listing"
)

type dbPokemon struct {
	ID             int    `db:"id"`
	Name           string `db:"name"`
	Height         int    `db:"height"`
	Weight         int    `db:"weight"`
	Order          int    `db:"order"`
	BaseExperience int    `db:"base_experience"`
	IsDefault      bool   `db:"is_default"`
}

func newDBPokemon(pokemon *listing.Pokemon) *dbPokemon {
	return (*dbPokemon)(pokemon)
}

func (m *dbPokemon) Deserialize() *listing.Pokemon {
	return (*listing.Pokemon)(m)
}
