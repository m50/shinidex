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

func (db *Database) Migrate(path string) error {
	path = strings.TrimRight(path, "/")
	files, err := os.ReadDir(path)
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
		var res int
		if err := tx.Get(&res, "SELECT count(*) FROM migrations WHERE id = $1;", f.Name()); err != nil {
			tx.Rollback()
			return err
		}
		if res == 1 {
			continue
		}
		db.info("Migrating " + f.Name() + "...")

		sql, err := os.ReadFile(path + "/" + f.Name())
		if err != nil {
			tx.Rollback()
			return err
		}
		queries := strings.Split(string(sql), ";")
		for _, q := range queries {
			if strings.TrimSpace(q) == "" {
				continue
			}
			r, err := tx.Exec(q)
			if err != nil {
				tx.Rollback()
				return err
			}
			cmd := strings.ToLower(q[:6])
			if cmd == "insert" || cmd == "delete" || cmd == "update" {
				rows, err := r.RowsAffected()
				if err != nil {
					tx.Rollback()
					return err
				}
				if rows < 1 {
					tx.Rollback()
					return fmt.Errorf("Insert query failed to %s rows for %s", cmd, f.Name())
				}
			}
		}
		_, err = tx.Exec("INSERT INTO migrations (id) VALUES ($1);", f.Name())
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}
