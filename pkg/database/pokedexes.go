package database

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/m50/shinidex/pkg/types"
)

type PokedexesDB struct {
	conn *sqlx.DB
}

func (db Database) Pokedexes() PokedexesDB {
	return PokedexesDB(db)
}

func (db PokedexesDB) FindByOwnerID(ctx context.Context, id string) ([]types.Pokedex, error) {
	dexes := []types.Pokedex{}
	err := db.conn.SelectContext(ctx, &dexes, "SELECT * FROM pokedexes WHERE owner_id = $1 ORDER BY updated;", id)
	return dexes, err
}

func (db PokedexesDB) FindByID(ctx context.Context, id string) (types.Pokedex, error) {
	dex := types.Pokedex{}
	err := db.conn.GetContext(ctx, &dex, "SELECT * FROM pokedexes WHERE id = $1;", id)
	return dex, err
}

func (db PokedexesDB) Insert(ctx context.Context, p types.Pokedex) (string, error) {
	q := `
	INSERT INTO pokedexes (id, name, owner_id, config, created, updated)
	VALUES (:id, :name, :owner_id, :config, :created, :updated);
	`
	p.ID = generateId()
	p.Created = time.Now()
	p.Updated = time.Now()
	_, err := db.conn.NamedExecContext(ctx, q, p)
	return p.ID, err
}

func (db PokedexesDB) Update(ctx context.Context, p types.Pokedex) error {
	q := `
	UPDATE pokedexes
	SET config = :config,
		updated = :updated
	WHERE id = :id;
	`
	p.Updated = time.Now()
	_, err := db.conn.NamedExecContext(ctx, q, p)
	return err
}

func (db PokedexesDB) Delete(ctx context.Context, id string) error {
	q := "DELETE FROM pokedexes WHERE id = $1;"
	_, err := db.conn.ExecContext(ctx, q, id)
	return err
}

type PokedexEntriesDB struct {
	conn *sqlx.DB
}

func (db PokedexesDB) Entries() PokedexEntriesDB {
	return PokedexEntriesDB(db)
}

func (db PokedexEntriesDB) Catch(ctx context.Context, pokedexID, pokemonID, formID string) error {
	q := `INSERT INTO pokedex_entries (pokedex_id, pokemon_id, form_id, created, updated)
	VALUES (:pokedex_id, :pokemon_id, :form_id, :created, :updated)`

	entry := types.PokedexEntry{
		PokedexID: pokedexID,
		PokemonID: pokemonID,
		FormID:    formID,
		Created:   time.Now(),
		Updated:   time.Now(),
	}

	r, err := db.conn.NamedExecContext(ctx, q, entry)
	if err != nil {
		return err
	}
	i, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if i < 1 {
		return errors.New("no rows inserted")
	}
	return nil
}

func (db PokedexEntriesDB) Release(ctx context.Context, pokedexID, pokemonID, formID string) error {
	q := `DELETE FROM pokedex_entries WHERE pokedex_id = $1 AND pokemon_id = $2 AND form_id = $3;`
	_, err := db.conn.ExecContext(ctx, q, pokedexID, pokemonID, formID)
	return err
}

func (db PokedexEntriesDB) List(ctx context.Context, pokedexID string) ([]types.PokedexEntry, error) {
	entries := []types.PokedexEntry{}
	err := db.conn.SelectContext(ctx, &entries, "SELECT * FROM pokedex_entries WHERE pokedex_id = $1", pokedexID)
	return entries, err
}
