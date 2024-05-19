PRAGMA foreign_keys=on;

CREATE TABLE users (
	id TEXT PRIMARY KEY,
	email TEXT NOT NULL,
	password TEXT NOT NULL,
	created DATETIME NOT NULL,
	updated DATETIME NOT NULL
);

CREATE TABLE pokemon (
	id TEXT PRIMARY KEY,
	national_dex_number INTEGER NOT NULL,
	name TEXT NOT NULL,
	shiny_locked INTEGER(1) NOT NULL
);

CREATE TABLE pokemon_forms (
	id TEXT NOT NULL,
	pokemon_id TEXT NOT NULL,
	name TEXT NOT NULL,
	shiny_locked INTEGER(1) NOT NULL,
	FOREIGN KEY(pokemon_id) REFERENCES pokemon(id),
	CONSTRAINT pokemon_form_id_cstrt PRIMARY KEY (id, pokemon_id)
);

CREATE TABLE pokedexes (
	id TEXT PRIMARY KEY,
	owner_id TEXT NOT NULL,
	name TEXT NOT NULL,
	config TEXT NOT NULL,
	created DATETIME NOT NULL,
	updated DATETIME NOT NULL,
	FOREIGN KEY(owner_id) REFERENCES users(id)
);

CREATE TABLE pokedex_entries (
	pokedex_id TEXT NOT NULL,
	pokemon_id TEXT NOT NULL,
	form_id TEXT NOT NULL,
	created DATETIME NOT NULL,
	updated DATETIME NOT NULL,
	FOREIGN KEY(pokedex_id) REFERENCES pokedexes(id) ON DELETE CASCADE,
	FOREIGN KEY(pokemon_id) REFERENCES pokemon(id),
	CONSTRAINT pokemon_entry_id_cstrt PRIMARY KEY (pokedex_id, pokemon_id, form_id)
);