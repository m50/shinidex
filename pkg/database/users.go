package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/m50/shinidex/pkg/types"
)

type UserDB struct {
	conn *sqlx.DB
}

func (db Database) Users() UserDB {
	return UserDB(db)
}

func (db UserDB) FindByID(ctx context.Context, id string) (types.User, error) {
	user := types.User{}
	err := db.conn.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
	return user, err
}

func (db UserDB) FindByEmail(ctx context.Context, email string) (types.User, error) {
	user := types.User{}
	err := db.conn.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email)
	return user, err
}

func (db UserDB) FindOrMake(ctx context.Context, user types.User) (types.User, error) {
	u, err := db.FindByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return user, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		id, err := db.Insert(ctx, user)
		if err != nil {
			return user, err
		}

		return db.FindByID(ctx, id)
	}
	return u, nil
}

func (db UserDB) Insert(ctx context.Context, user types.User) (string, error) {
	query := `
	INSERT INTO users (id, email, password, created, updated)
	VALUES (:id, :email, :password, :created, :updated);
	`
	user.Created = time.Now()
	user.Updated = time.Now()
	user.ID = generateId()
	_, err := db.conn.NamedExecContext(ctx, query, user)
	return user.ID, err
}

func (db UserDB) Update(ctx context.Context, user types.User) error {
	query := `
	UPDATE users
	SET email = :email,
		password = :password,
		created = :created,
		updated = :updated
	WHERE id = :id;
	`

	user.Updated = time.Now()
	_, err := db.conn.NamedExecContext(ctx, query, user)
	return err
}

func (db UserDB) Delete(ctx context.Context, id string) error {
	_, err := db.conn.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}
