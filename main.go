package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/web"
)

func main() {
	if _, err := os.Stat("./.env"); err == nil {
		if err = godotenv.Load("./.env"); err != nil {
			log.Fatalf("error loading .env file: %s", err)
		}
	}

	db, err := database.NewFromEnv()
	if err != nil {
		log.Fatalf("DB failed to open: %s", err)
		return
	}
	defer db.Close()
	if err = db.Migrate("./migrations"); err != nil {
		log.Fatalf("Failed to migrate: %s", err)
		return
	}

	e := web.New()
	e.Logger.Info("test")
	// e.Logger.Fatal(e.Start(":1323"))
}