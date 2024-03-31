package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/database/passwords"
	"github.com/m50/shinidex/pkg/types"
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
	p, err := passwords.HashPassword("test")
	if err != nil {
		log.Fatalf("failed to hash password for test user")
		return
	}
	if _, err := db.Users().Insert(types.User{
		Email: "test@test.com",
		Password: p,
	}); err != nil {
		log.Fatalf("Failed to insert test user")
		return
	}

	e := web.New(db)
	e.Logger.Fatal(e.Start(":1323"))
}