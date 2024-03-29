package database

import (
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

func (db PokedexesDB) FindByOwnerID(id string) ([]types.Pokedex, error) {
	dexes := []types.Pokedex{}
	err := db.conn.Select(dexes, "SELECT * FROM pokedexes WHERE owner_id = $1;", id)
	return dexes, err
}

func (db PokedexesDB) FindByID(id string) (types.Pokedex, error) {
	dex := types.Pokedex{}
	err := db.conn.Get(dex, "SELECT * FROM pokedexes WHERE id = $1;", id)
	return dex, err
}

func (db PokedexesDB) Insert(p types.Pokedex) error {
	q := `
	INSERT INTO pokedexes (id, owner_id, config, created, updated)
	VALUES (:id, :owner_id, :config, :created, :updated);
	`
	p.ID = generateId()
	p.Created = time.Now().UTC().Unix()
	p.Updated = time.Now().UTC().Unix()
	_, err := db.conn.NamedExec(q, p)
	return err
}

func (db PokedexesDB) Update(p types.Pokedex) error {
	q := `
	UPDATE pokedexes
	SET config = :config,
		updated = :updated
	WHERE id = :id;
	`
	p.Updated = time.Now().UTC().Unix()
	_, err := db.conn.NamedExec(q, p)
	return err
}

func (db PokedexesDB) Delete(id string) error {
	q := "DELETE FROM pokedexes WHERE id = $1;"
	_, err := db.conn.Exec(q, id)
	return err
}

type PokedexEntriesDB struct {
	conn *sqlx.DB
}

func (db PokedexesDB) Entries() PokedexEntriesDB {
	return PokedexEntriesDB(db)
}

func (db PokedexEntriesDB) Catch(pokedexID, pokemonID string) error {
	q := `
	INSERT INTO pokedex_entries (pokedex_id, pokemon_id, form_id, created, updated)
	VALUES (:pokedex_id, :pokemon_id, NULL, :created, :updated);
	`
	entry := types.PokedexEntry{
		PokedexID: pokedexID,
		PokemonID: pokemonID,
		Created:   time.Now().UTC().Unix(),
		Updated:   time.Now().UTC().Unix(),
	}

	_, err := db.conn.NamedExec(q, entry)
	return err
}

func (db PokedexEntriesDB) Release(pokedexID, pokemonID string) error {
	q := `DELETE FROM pokedex_entries WHERE pokedex_id = $1, pokemon_id = $2;`
	_, err := db.conn.Exec(q, pokedexID, pokemonID)
	return err
}

func (db PokedexEntriesDB) CatchForm(pokedexID, pokemonID, formID string) error {
	q := `
	INSERT INTO pokedex_entries (pokedex_id, pokemon_id, form_id, created, updated)
	VALUES (:pokedex_id, :pokemon_id, :form_id, :created, :updated);
	`
	entry := types.PokedexEntry{
		PokedexID: pokedexID,
		PokemonID: pokemonID,
		FormID:    formID,
		Created:   time.Now().UTC().Unix(),
		Updated:   time.Now().UTC().Unix(),
	}

	_, err := db.conn.NamedExec(q, entry)
	return err
}

func (db PokedexesDB) ReleaseForm(pokedexID, pokemonID, formID string) error {
	q := `DELETE FROM pokedex_entries WHERE pokedex_id = $1, pokemon_id = $2, form_id = $3;`
	_, err := db.conn.Exec(q, pokedexID, pokemonID, formID)
	return err
}
