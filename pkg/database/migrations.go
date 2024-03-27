package database

import (
	"fmt"
	"os"
	"strings"
)

const migrationSchema = `
CREATE TABLE migrations (
	id TEXT PRIMARY KEY
);
`

func (db *Database) Migrate() error {
	files, err := os.ReadDir("./migrations")
	if err != nil {
		return err
	}
	if err = db.conn.Ping(); err != nil {
		return err
	}
	_, err = db.conn.Exec("SELECT * FROM migrations;")
	if err != nil {
		if _, err = db.conn.Exec(migrationSchema); err != nil {
			return err
		}
	}

	tx := db.conn.MustBegin()
	for _, f := range files {
		fmt.Println("Migrating " + f.Name())
		res := &struct {
			C int `db:"c"`
		}{}
		if err := db.conn.Get(res, "SELECT count(*) as c FROM migrations WHERE id = $1;", f.Name()); err != nil {
			tx.Rollback()
			return err
		}
		if res.C == 1 {
			continue
		}

		sql, err := os.ReadFile("./migrations/" + f.Name())
		if err != nil {
			tx.Rollback()
			return err
		}
		queries := strings.Split(string(sql), ";")
		for _, q := range queries {
			_, err = db.conn.Exec(q + ";")
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		_, err = db.conn.Exec("INSERT INTO migrations (id) VALUES ($1);", f.Name())
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}
