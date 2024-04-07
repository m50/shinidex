package database

import (
	"fmt"
	"slices"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/m50/shinidex/pkg/math"
	"github.com/m50/shinidex/pkg/types"
)

type PokemonDB struct {
	conn *sqlx.DB
	logger *log.Logger
}

type PokemonFormsDB struct {
	conn *sqlx.DB
	logger *log.Logger
}

func (db Database) Pokemon() PokemonDB {
	return PokemonDB(db)
}

func (db PokemonDB) Forms() PokemonFormsDB {
	return PokemonFormsDB(db)
}

func (db PokemonDB) GetAll() (types.PokemonList, error) {
	pokemon := types.PokemonList{}
	err := db.conn.Select(&pokemon, "SELECT * FROM pokemon ORDER BY national_dex_number")
	return pokemon, err
}

func (db PokemonDB) Get(rows, page int) (types.PokemonList, error) {
	pokemon := types.PokemonList{}
	offset := math.Max(page-1, 0) * rows
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

func (db PokemonDB) GetAllAsSeparatForms() (types.PokemonList, error) {
	pokemon, err := db.GetAll()
	if err != nil {
		return nil, err
	}
	forms, err := db.Forms().GetAll()
	if err != nil {
		return nil, err
	}
	result := make(types.PokemonList, len(pokemon)+len(forms))
	c := -1
	for _, pkmn := range pokemon {
		c++
		result[c] = pkmn
		for _, f := range forms {
			if f.PokemonID != pkmn.ID {
				continue
			}
			c++
			result[c] = types.Pokemon{
				ID:                pkmn.ID + "-" + f.ID,
				Name:              fmt.Sprintf("%s (%s)", pkmn.Name, f.Name),
				NationalDexNumber: pkmn.NationalDexNumber,
				ShinyLocked:       f.ShinyLocked,
			}
		}
	}
	return slices.Clip(result), nil
}
