package main

import (
	"app/config"
	"app/db/seed"
	"fmt"
	"log"
)

func main() {
	db, err := config.NewMysqlConnection()
	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}

	if err := config.ExecuteMigrate(db); err != nil {
		log.Fatalf("failed to migrate: %s", err)
	}

	if err := seed.Run(db); err != nil {
		log.Fatalf("failed to seed: %s", err)
	}

	fmt.Println("Seed completed successfully")
}
