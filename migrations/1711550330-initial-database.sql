PRAGMA foreign_keys=on;

CREATE TABLE users (
	id TEXT PRIMARY KEY,
	email TEXT NOT NULL,
	password TEXT NOT NULL,
	created INTEGER NOT NULL,
	updated INTEGER NOT NULL
);

CREATE TABLE pokemon (
	id TEXT PRIMARY KEY,
	name TEXT
);

CREATE TABLE pokemon_forms (
	id TEXT,
	pokemon_id TEXT,
	FOREIGN KEY(pokemon_id) REFERENCES pokemon(id),
	CONSTRAINT pokemon_form_id_cstrt PRIMARY KEY (id, pokemon_id)
);

CREATE TABLE dex (
	id TEXT PRIMARY KEY,
	owner_id TEXT NOT NULL,
	config TEXT NOT NULL,
	created INTEGER NOT NULL,
	updated INTEGER NOT NULL,
	FOREIGN KEY(owner_id) REFERENCES users(id)
);

CREATE TABLE pokedex_entry (
	dex_id TEXT NOT NULL,
	pokemon_id TEXT NOT NULL,
	form_id TEXT,
	created INTEGER NOT NULL,
	updated INTEGER NOT NULL,
	FOREIGN KEY(dex_id) REFERENCES dex(id) ON DELETE CASCADE,
	FOREIGN KEY(pokemon_id) REFERENCES pokemon(id)
	CONSTRAINT pokemon_entry_id_cstrt PRIMARY KEY (dex_id, pokemon_id)
);