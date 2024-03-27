package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/tursodatabase/go-libsql"
)

type Database struct {
	conn *sqlx.DB
}

func New() (*Database, error) {
	file := os.Getenv("DB_PATH")
	if file == "" {
		return nil, fmt.Errorf("DB_PATH envvar needs to be set")
	}
	dbUrl := os.Getenv("TURSO_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	var db *sqlx.DB
	var err error
	if dbUrl != "" && authToken != "" {
		var connector driver.Connector
		connector, err = libsql.NewEmbeddedReplicaConnector(file, dbUrl, libsql.WithAuthToken(authToken))
		if err != nil {
			return nil, err
		}
		dbi := sql.OpenDB(connector)
		db = sqlx.NewDb(dbi, "libsql")
	} else {
		db, err = sqlx.Open("libsql", file)
		if err != nil {
			return nil, err
		}
	}

	return &Database{
		conn: db,
	}, nil
}

func (db *Database) Close() {
	db.conn.Close()
}
