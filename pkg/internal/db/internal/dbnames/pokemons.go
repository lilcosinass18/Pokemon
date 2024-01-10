package dbnames

type pokemonsFieldsT struct {
	ID             FieldName
	Name           FieldName
	Height         FieldName
	Weight         FieldName
	Order          FieldName
	BaseExperience FieldName
	IsDefault      FieldName
}

var pokemonsFields = pokemonsFieldsT{
	ID:             "id",
	Name:           "name",
	Height:         "height",
	Weight:         "weight",
	Order:          "\"order\"",
	BaseExperience: "base_experience",
	IsDefault:      "is_default",
}

type pokemonsT struct {
	tableT
	pokemonsFieldsT
}

var pokemons = pokemonsT{
	tableT:          tableT{TName: "pokemons"},
	pokemonsFieldsT: pokemonsFields,
}
