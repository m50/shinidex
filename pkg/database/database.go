package database

import (
	crand "crypto/rand"
	"database/sql"
	"database/sql/driver"
	"math/big"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/m50/shinidex/pkg/config"
	"github.com/tursodatabase/go-libsql"
)

type DBContext interface {
	DB() *Database
}

type Database struct {
	conn *sqlx.DB
}

func NewFromLoadedConfig() (*Database, error) {
	file := config.Loaded.DBPath
	dbUrl := config.Loaded.TursoURL
	authToken := config.Loaded.TursoAuthToken
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
		k, _ := crand.Int(crand.Reader, big.NewInt(now))
		r := k.Int64()
		randComponent += strconv.FormatInt(r, 36)
	}

	max := len(randComponent) - (32 - len(timeComponent))
	start, _ := crand.Int(crand.Reader, big.NewInt(int64(max)))
	randComponent = randComponent[start.Int64():]

	id := timeComponent + randComponent
	return id[:32]
}
