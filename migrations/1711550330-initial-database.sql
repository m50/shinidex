CREATE TABLE users (
	id TEXT,
	email TEXT
);

CREATE TABLE pokemon (
	id TEXT,
	name TEXT
);

CREATE TABLE dex (
	id TEXT,
	owner_id TEXT
);

CREATE TABLE pokedex_entry (
	dex_id TEXT,
	pokemon_id TEXT,
	caught BOOLEAN NOT NULL CHECK (caught IN (0, 1))
);