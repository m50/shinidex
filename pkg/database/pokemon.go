package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/m50/shinidex/pkg/math"
	"github.com/m50/shinidex/pkg/types"
)

type PokemonDB struct {
	conn *sqlx.DB
}

type PokemonFormsDB struct {
	conn *sqlx.DB
}

func (db Database) Pokemon() PokemonDB {
	return PokemonDB(db)
}

func (db PokemonDB) Forms() PokemonFormsDB {
	return PokemonFormsDB(db)
}

func (db PokemonDB) GetAll() ([]types.Pokemon, error) {
	pokemon := []types.Pokemon{}
	err := db.conn.Select(&pokemon, "SELECT * FROM pokemon ORDER BY national_dex_number")
	return pokemon, err
}

func (db PokemonDB) Get(rows, page int) ([]types.Pokemon, error) {
	pokemon := []types.Pokemon{}
	offset := math.Max(page - 1, 0) * rows
	err := db.conn.Select(&pokemon, "SELECT * FROM pokemon ORDER BY national_dex_number LIMIT $1 OFFSET $2", rows, offset)
	return pokemon, err
}

func (db PokemonDB) FindByID(id string) (types.Pokemon, error) {
	pokemon := types.Pokemon{}
	err := db.conn.Get(&pokemon, "SELECT * FROM pokemon WHERE id = $1", id)
	return pokemon, err
}

func (db PokemonFormsDB) GetAll() ([]types.PokemonForm, error) {
	forms := []types.PokemonForm{}
	err := db.conn.Select(&forms, "SELECT * FROM pokemon_forms")
	return forms, err
}

func (db PokemonFormsDB) FindByPokemonID(pokemonID string) ([]types.PokemonForm, error) {
	forms := []types.PokemonForm{}
	err := db.conn.Select(&forms, "SELECT * FROM pokemon_forms WHERE pokemon_id = $1", pokemonID)
	return forms, err
}

func (db PokemonFormsDB) FindByID(pokemonID string, formID string) (types.PokemonForm, error) {
	form := types.PokemonForm{}
	err := db.conn.Get(&form, "SELECT * FROM pokemon_forms WHERE id = $1 AND pokemon_id = $2", formID, pokemonID)
	return form, err
}

func (db PokemonDB) FindWithFormsByID(pokemonID string) (types.PokemonWithForms, error) {
	pokemon, err := db.FindByID(pokemonID)
	if err != nil {
		return types.PokemonWithForms{}, err
	}
	forms, err := db.Forms().FindByPokemonID(pokemonID)
	if err != nil {
		return types.PokemonWithForms{}, err
	}

	return types.PokemonWithForms{
		Pokemon: pokemon,
		Forms: forms,
	}, nil
}

func (db PokemonDB) GetAllAsSeparatForms() ([]types.Pokemon, error) {
	pokemon := []types.Pokemon{}
	err := db.conn.Select(&pokemon, "SELECT * FROM pokemon ORDER BY national_dex_number")
	return pokemon, err
}
