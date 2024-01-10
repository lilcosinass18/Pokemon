package dbnames

type tablesT struct {
	Pokemons pokemonsT
}

var Tables = &tablesT{
	Pokemons: pokemons,
}
