package database

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	Id       string
	Email    string
	Password string
	Created  time.Time
	Updated  time.Time
}

type UserDB struct {
	conn *sqlx.DB
}

func (db *Database) Users() *UserDB {
	return &UserDB{conn: db.conn}
}

func (db *UserDB) FindById(id string) (User, error) {
	user := User{}
	err := db.conn.Get(user, "SELECT * FROM users WHERE id = $1", id)
	return user, err
}

func (db *UserDB) FindByEmail(email string) (User, error) {
	user := User{}
	err := db.conn.Get(user, "SELECT * FROM users WHERE email = $1", email)
	return user, err
}

func (db *UserDB) FindByUsername(username string) (User, error) {
	user := User{}
	err := db.conn.Get(user, "SELECT * FROM users WHERE username = $1", username)
	return user, err
}

func (db *UserDB) Insert(user User) error {
	query := `
	INSERT INTO users (id, email, password, created, updated)
	VALUES (:id, :email, :password, :created, :updated);
	`
	user.Created = time.Now()
	user.Updated = time.Now()
	user.Id = generateId()
	_, err := db.conn.NamedExec(query, user)
	return err
}

func (db *UserDB) Update(user User) error {
	query := `
	UPDATE users
	SET email = :email,
		password = :password,
		created = :created,
		updated = :updated,
	WHERE id = :id;
	`

	user.Updated = time.Now()
	_, err := db.conn.NamedExec(query, user)
	return err
}

func (db *UserDB) Delete(id string) error {
	_, err := db.conn.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}