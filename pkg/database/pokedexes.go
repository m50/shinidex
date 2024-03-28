package database

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

type Pokedex struct {
	ID      string
	OwnerID string `db:"owner_id"`
	config  string
	Created int64
	Updated int64
}

func NewPokedex(ownerID string, config PokedexConfig) (Pokedex, error) {
	c, err := json.Marshal(config)
	if err != nil {
		return Pokedex{}, err
	}
	return Pokedex{
		OwnerID: ownerID,
		config:  string(c),
	}, nil
}

type FormLocation int
const (
	Off FormLocation = iota
	Inline
	After
)

type PokedexConfig struct {
	Shiny         bool
	GenderForms   FormLocation
	RegionalForms FormLocation
}

func (p Pokedex) Config() (PokedexConfig, error) {
	var conf PokedexConfig
	err := json.Unmarshal([]byte(p.config), &conf)
	return conf, err
}
func (p *Pokedex) UpdateConfig(config PokedexConfig) error {
	c, err := json.Marshal(config)
	if err != nil {
		return err
	}
	p.config = string(c)
	return nil
}

type PokedexEntry struct {
	PokedexID string `db:"pokedex_id"`
	PokemonID string `db:"pokemon_id"`
	FormID    string `db:"form_id"`
	Created   int64
	Updated   int64
}

type PokedexesDB struct {
	conn *sqlx.DB
}

func (db Database) Pokedexes() PokedexesDB {
	return PokedexesDB(db)
}

func (db PokedexesDB) FindByOwnerID(id string) ([]Pokedex, error) {
	dexes := []Pokedex{}
	err := db.conn.Select(dexes, "SELECT * FROM pokedexes WHERE owner_id = $1;", id)
	return dexes, err
}

func (db PokedexesDB) FindByID(id string) (Pokedex, error) {
	dex := Pokedex{}
	err := db.conn.Get(dex, "SELECT * FROM pokedexes WHERE id = $1;", id)
	return dex, err
}

func (db PokedexesDB) Insert(p Pokedex) error {
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

func (db PokedexesDB) Update(p Pokedex) error {
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
	entry := PokedexEntry{
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
	entry := PokedexEntry{
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
