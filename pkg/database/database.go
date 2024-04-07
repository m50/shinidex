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
	"github.com/labstack/gommon/log"
	"github.com/tursodatabase/go-libsql"
)

type DBContext interface {
	DB() *Database
}

type Database struct {
	conn *sqlx.DB
	logger *log.Logger
}

func (db Database) debug(i ...interface{}) {
	if db.logger == nil {
		return
	}
	db.logger.Debug(i...)
}

func (db Database) debugf(msg string, i ...interface{}) {
	if db.logger == nil {
		return
	}
	db.logger.Debugf(msg, i...)
}

func (db Database) info(i ...interface{}) {
	if db.logger == nil {
		return
	}
	db.logger.Info(i...)
}

func (db Database) infof(msg string, i ...interface{}) {
	if db.logger == nil {
		return
	}
	db.logger.Infof(msg, i...)
}

func (db Database) warn(i ...interface{}) {
	if db.logger == nil {
		return
	}
	db.logger.Warn(i...)
}

func (db Database) warnf(msg string, i ...interface{}) {
	if db.logger == nil {
		return
	}
	db.logger.Warnf(msg, i...)
}

func (db Database) error(i ...interface{}) {
	if db.logger == nil {
		return
	}
	db.logger.Error(i...)
}

func (db Database) errorf(msg string, i ...interface{}) {
	if db.logger == nil {
		return
	}
	db.logger.Errorf(msg, i...)
}

func (db *Database) AttachLogger(logger *log.Logger) {
	db.logger = logger
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
	id = id[:32]
	return id
}