package config

import (
	"errors"
	"os"
)

type Conf struct {
	DisallowRegistration bool
	DBPath               string
	AuthKey              []byte
	TursoAuthToken       string
	TursoURL             string
}

var Loaded Conf

func LoadConfigFromEnv() error {
	c := Conf{
		DBPath:               os.Getenv("DB_PATH"),
		AuthKey:              []byte(os.Getenv("AUTH_KEY")),
		TursoAuthToken:       os.Getenv("TURSO_AUTH_TOKEN"),
		TursoURL:             os.Getenv("TURSO_URL"),
		DisallowRegistration: os.Getenv("DISALLOW_REGISTRATION") == "true",
	}
	if c.DBPath == "" {
		return errors.New("DB_PATH envvar needs to be set")
	}
	if len(c.AuthKey) == 0 {
		c.AuthKey = make([]byte, 32)
	}

	Loaded = c
	return nil
}
