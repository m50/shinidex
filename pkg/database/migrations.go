package database

import (
	"embed"
	"fmt"
	"strings"

	"github.com/gookit/slog"
)

const migrationSchema = `
CREATE TABLE migrations (
	id TEXT PRIMARY KEY
);
`

//go:embed migrations/*.sql
var migrationFS embed.FS

func (db *Database) Migrate() error {
	files, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return err
	}
	if err = db.conn.Ping(); err != nil {
		return err
	}
	_, err = db.conn.Exec("SELECT * FROM migrations;")
	if err != nil {
		if _, err = db.conn.Exec(migrationSchema); err != nil {
			slog.Warnf("possible failure in migrating: %v", err)
		}
	}

	tx := db.conn.MustBegin()
	defer tx.Rollback()
	for _, f := range files {
		var res int
		if err := tx.Get(&res, "SELECT count(*) FROM migrations WHERE id = $1;", f.Name()); err != nil {
			tx.Rollback()
			return err
		}
		if res == 1 {
			continue
		}
		slog.Debugf("migrating %s...", f.Name())

		sql, err := migrationFS.ReadFile("migrations/" + f.Name())
		if err != nil {
			return err
		}
		queries := strings.Split(string(sql), ";")
		for _, q := range queries {
			if strings.TrimSpace(q) == "" {
				continue
			}
			r, err := tx.Exec(q)
			if err != nil {
				return err
			}
			cmd := strings.ToLower(q[:6])
			if cmd == "insert" || cmd == "delete" || cmd == "update" {
				rows, err := r.RowsAffected()
				if err != nil {
					return err
				}
				if rows < 1 {
					return fmt.Errorf("Insert query failed to %s rows for %s", cmd, f.Name())
				}
			}
		}
		_, err = tx.Exec("INSERT INTO migrations (id) VALUES ($1);", f.Name())
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}
