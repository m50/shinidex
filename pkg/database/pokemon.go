package database

import "github.com/jmoiron/sqlx"

type Pokemon struct {
	ID   string
	Name string
}

type PokemonForm struct {
	ID        string
	Name      string
	PokemonID string `db:"pokemon_id"`
}

type PokemonWithForms struct {
	*Pokemon
	Forms []PokemonForm
}

type PokemonDB struct {
	conn *sqlx.DB
}

func (db Database) Pokemon() PokemonDB {
	return PokemonDB(db)
}

func (db PokemonDB) GetAll() ([]Pokemon, error) {
	pokemon := []Pokemon{}
	err := db.conn.Select(&pokemon, "SELECT * FROM pokemon")
	return pokemon, err
}

func (db PokemonDB) FindByID(id string) (Pokemon, error) {
	pokemon := Pokemon{}
	err := db.conn.Get(&pokemon, "SELECT * FROM pokemon WHERE id = $1", id)
	return pokemon, err
}

func (db PokemonDB) FindFormsForPokemon(pokemonID string) ([]PokemonForm, error) {
	forms := []PokemonForm{}
	err := db.conn.Select(&forms, "SELECT * FROM pokemon_forms WHERE pokemon_id = $1", pokemonID)
	return forms, err
}

func (db PokemonDB) FindFormByID(pokemonID string, formID string) (PokemonForm, error) {
	form := PokemonForm{}
	err := db.conn.Get(&form, "SELECT * FROM pokemon_forms WHERE id = $1 AND pokemon_id = $2", formID, pokemonID)
	return form, err
}

func (db PokemonDB) FindWithFormsByID(pokemonID string) (PokemonWithForms, error) {
	pokemon, err := db.FindByID(pokemonID)
	if err != nil {
		return PokemonWithForms{}, err
	}
	forms, err := db.FindFormsForPokemon(pokemonID)
	if err != nil {
		return PokemonWithForms{}, err
	}

	return PokemonWithForms{
		Pokemon: &pokemon,
		Forms: forms,
	}, nil
}
