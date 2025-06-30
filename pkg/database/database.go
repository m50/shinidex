package database

import (
	crand "crypto/rand"
	"math/big"
	"strconv"
	"time"

	"github.com/gookit/slog"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/xo/dburl"
)

type DBContext interface {
	DB() *Database
}

type Database struct {
	conn *sqlx.DB
}

func NewFromLoadedConfig() (*Database, error) {
	return New(viper.GetString("db-url"))
}

func New(dbURL string) (*Database, error) {
	slog.Debugf("connecting to database at %v", dbURL)
	url, err := dburl.Parse(dbURL)
	if err != nil {
		slog.Errorf("unable to parse database url: %s", err)
		return nil, err
	}
	driver := url.Driver
	if url.GoDriver != "" {
		driver = url.GoDriver
	}
	db, err := sqlx.Connect(driver, url.DSN)
	if err != nil {
		slog.Errorf("unable to connect to database: %s", err)
		return nil, err
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
