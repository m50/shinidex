package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/tursodatabase/go-libsql"
)

type DBContext interface {
	DB() *Database
}

type Database struct {
	conn *sqlx.DB
}

func NewFromEnv() (*Database, error) {
	file := os.Getenv("DB_PATH")
	if file == "" {
		return nil, fmt.Errorf("DB_PATH envvar needs to be set")
	}
	dbUrl := os.Getenv("TURSO_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")
	return New(file, dbUrl, authToken)
}

func NewLocal(file string) (*Database, error) {
	return New(file, "", "")
}

func New(file, dbUrl, authToken string) (*Database, error) {
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

func generateId() string {
	now := time.Now().UnixMicro()
	timeComponent := strconv.FormatInt(now, 36)
	timeComponent = timeComponent[2:]

	randComponent := ""
	for i := 0; i < 5; i++ {
		r := rand.Uint64()
		randComponent += strconv.FormatUint(r, 36)
	}

	max := len(randComponent) - (32 - len(timeComponent))
	start := rand.Intn(max)
	randComponent = randComponent[start:]

	id := timeComponent + randComponent
	return id[:32]
}