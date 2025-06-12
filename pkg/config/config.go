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
	WebAddress           string
}

var Loaded Conf

func LoadConfigFromEnv() error {
	c := Conf{
		DBPath:               os.Getenv("DB_PATH"),
		AuthKey:              []byte(os.Getenv("AUTH_KEY")),
		TursoAuthToken:       os.Getenv("TURSO_AUTH_TOKEN"),
		TursoURL:             os.Getenv("TURSO_URL"),
		DisallowRegistration: os.Getenv("DISALLOW_REGISTRATION") == "true",
		WebAddress:           os.Getenv("WEB_ADDRESS"),
	}
	if c.DBPath == "" {
		return errors.New("DB_PATH envvar needs to be set")
	}
	if len(c.AuthKey) == 0 {
		c.AuthKey = make([]byte, 32)
	}
	if c.WebAddress == "" {
		c.WebAddress = ":1323"
	}

	Loaded = c
	return nil
}
