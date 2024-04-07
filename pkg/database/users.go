package database

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/m50/shinidex/pkg/types"
)

type UserDB struct {
	conn *sqlx.DB
	logger *log.Logger
}

func (db Database) Users() UserDB {
	return UserDB(db)
}

func (db UserDB) FindByID(id string) (types.User, error) {
	user := types.User{}
	err := db.conn.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	return user, err
}

func (db UserDB) FindByEmail(email string) (types.User, error) {
	user := types.User{}
	err := db.conn.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	return user, err
}

func (db UserDB) Insert(user types.User) (string, error) {
	query := `
	INSERT INTO users (id, email, password, created, updated)
	VALUES (:id, :email, :password, :created, :updated);
	`
	user.Created = time.Now().UTC().Unix()
	user.Updated = time.Now().UTC().Unix()
	user.ID = generateId()
	_, err := db.conn.NamedExec(query, user)
	return user.ID, err
}

func (db UserDB) Update(user types.User) error {
	query := `
	UPDATE users
	SET email = :email,
		password = :password,
		created = :created,
		updated = :updated
	WHERE id = :id;
	`

	user.Updated = time.Now().UTC().Unix()
	_, err := db.conn.NamedExec(query, user)
	return err
}

func (db UserDB) Delete(id string) error {
	_, err := db.conn.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}