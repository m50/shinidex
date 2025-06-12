package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/m50/shinidex/pkg/config"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/database/passwords"
	imgdownloader "github.com/m50/shinidex/pkg/img-downloader"
	l "github.com/m50/shinidex/pkg/logger"
	"github.com/m50/shinidex/pkg/types"
	"github.com/m50/shinidex/pkg/web"
)

func logger() *log.Logger {
	logger := log.New("shinidex")
	logger.SetHeader("[${time_rfc3339}] ${level}")
	logger.SetLevel(log.DEBUG)
	return logger
}

func main() {
	if _, err := os.Stat("./.env"); err == nil {
		if err = godotenv.Load("./.env"); err != nil {
			log.Fatalf("error loading .env file: %s", err)
		}
	}
	if err := config.LoadConfigFromEnv(); err != nil {
		log.Fatal(err)
	}

	logger := logger()
	l.SetDefaultLogger(logger)

	db, err := database.NewFromLoadedConfig()
	if err != nil {
		logger.Fatalf("DB failed to open: %s", err)
		return
	}
	defer db.Close()
	if err = db.Migrate("./migrations"); err != nil {
		logger.Fatalf("Failed to migrate: %s", err)
		return
	}

	addTestUser(db, logger)

	go imgdownloader.DownloadImages(db, logger)

	e := web.New(db, logger)
	if err := e.Start(config.Loaded.WebAddress); err != nil {
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
		Email:    "test@test.com",
		Password: p,
	}); err != nil {
		logger.Fatalf("Failed to insert test user")
		return
	}
}
