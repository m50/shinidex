package database

import (
	"context"
	"fmt"
	"slices"
	"strings"

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

func (db PokemonDB) GetAll(ctx context.Context) (types.PokemonList, error) {
	pokemon := types.PokemonList{}
	err := db.conn.SelectContext(ctx, &pokemon, "SELECT * FROM pokemon ORDER BY national_dex_number")
	return pokemon, err
}

func (db PokemonDB) Get(ctx context.Context, rows, page int) (types.PokemonList, error) {
	pokemon := types.PokemonList{}
	offset := math.Max(page-1, 0) * rows
	err := db.conn.SelectContext(ctx, &pokemon, "SELECT * FROM pokemon ORDER BY national_dex_number LIMIT $1 OFFSET $2", rows, offset)
	return pokemon, err
}

func (db PokemonDB) FindByID(ctx context.Context, id string) (types.Pokemon, error) {
	pokemon := types.Pokemon{}
	err := db.conn.GetContext(ctx, &pokemon, "SELECT * FROM pokemon WHERE id = $1", id)
	return pokemon, err
}

func (db PokemonFormsDB) GetAll(ctx context.Context) ([]types.PokemonForm, error) {
	forms := []types.PokemonForm{}
	err := db.conn.SelectContext(ctx, &forms, "SELECT * FROM pokemon_forms")
	return forms, err
}

func (db PokemonFormsDB) FindByPokemonID(ctx context.Context, pokemonID string) ([]types.PokemonForm, error) {
	forms := []types.PokemonForm{}
	err := db.conn.SelectContext(ctx, &forms, "SELECT * FROM pokemon_forms WHERE pokemon_id = $1", pokemonID)
	return forms, err
}

func (db PokemonFormsDB) FindByID(ctx context.Context, pokemonID string, formID string) (types.PokemonForm, error) {
	form := types.PokemonForm{}
	err := db.conn.GetContext(ctx, &form, "SELECT * FROM pokemon_forms WHERE id = $1 AND pokemon_id = $2", formID, pokemonID)
	return form, err
}

func (db PokemonDB) FindByFullFormID(ctx context.Context, fullFormID string) (types.Pokemon, error) {
	parts := strings.Split(fullFormID, types.IDSeparator)
	pkmn, err := db.FindByID(ctx, parts[0])
	if err != nil || len(parts) == 1 {
		return pkmn, err
	}
	f, err := db.Forms().FindByID(ctx, parts[0], parts[1])
	if err != nil {
		return types.Pokemon{}, err
	}
	return types.Pokemon{
		ID:                pkmn.ID + types.IDSeparator + f.ID,
		Name:              fmt.Sprintf("%s (%s)", pkmn.Name, f.Name),
		NationalDexNumber: pkmn.NationalDexNumber,
		ShinyLocked:       f.ShinyLocked,
		Form:              true,
	}, nil
}

func (db PokemonDB) GetAllAsSeparateForms(ctx context.Context) (types.PokemonList, error) {
	pokemon, err := db.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	forms, err := db.Forms().GetAll(ctx)
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
				ID:                pkmn.ID + types.IDSeparator + f.ID,
				Name:              fmt.Sprintf("%s (%s)", pkmn.Name, f.Name),
				NationalDexNumber: pkmn.NationalDexNumber,
				ShinyLocked:       f.ShinyLocked,
				Form:              true,
			}
		}
	}
	return slices.Clip(result), nil
}
