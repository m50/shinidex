package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/database/passwords"
	"github.com/m50/shinidex/pkg/types"
	"github.com/m50/shinidex/pkg/web"
)

func logger() *log.Logger {
	logger := log.New("shinidex")
	logger.SetHeader("[${time_rfc3339}] ${level} ${short_file}:${line}")
	logger.SetLevel(log.DEBUG)
	return logger
}

func main() {
	if _, err := os.Stat("./.env"); err == nil {
		if err = godotenv.Load("./.env"); err != nil {
			log.Fatalf("error loading .env file: %s", err)
		}
	}

	logger := logger()

	db, err := database.NewFromEnv()
	if err != nil {
		logger.Fatalf("DB failed to open: %s", err)
		return
	}
	defer db.Close()
	db.AttachLogger(logger)
	if err = db.Migrate("./migrations"); err != nil {
		logger.Fatalf("Failed to migrate: %s", err)
		return
	}

	addTestUser(db, logger)

	e := web.New(db, logger)
	if err := e.Start(":1323"); err != nil {
		logger.Fatal(err)
	}
}

func addTestUser(db *database.Database, logger *log.Logger) {
	p, err := passwords.HashPassword("test")
	if err != nil {
		logger.Fatalf("failed to hash password for test user")
		return
	}
	if _, err := db.Users().Insert(types.User{
		Email: "test@test.com",
		Password: p,
	}); err != nil {
		logger.Fatalf("Failed to insert test user")
		return
	}
}