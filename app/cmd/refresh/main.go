package main

import (
	"app/config"
	"app/db/seed"
	"fmt"
	"log"
	"os"
)

func main() {
	dbName := os.Getenv("MYSQL_DATABASE")
	if dbName == "" {
		log.Fatal("MYSQL_DATABASE is not set")
	}

	rootDB, err := config.NewMysqlConnectionWithoutDB()
	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}

	fmt.Printf("Dropping database: %s\n", dbName)
	if err := rootDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", dbName)).Error; err != nil {
		log.Fatalf("failed to drop database: %s", err)
	}

	fmt.Printf("Creating database: %s\n", dbName)
	if err := rootDB.Exec(fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)).Error; err != nil {
		log.Fatalf("failed to create database: %s", err)
	}

	db, err := config.NewMysqlConnection()
	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}

	fmt.Println("Running migrations...")
	if err := config.ExecuteMigrate(db); err != nil {
		log.Fatalf("failed to migrate: %s", err)
	}

	fmt.Println("Running seeds...")
	if err := seed.Run(db); err != nil {
		log.Fatalf("failed to seed: %s", err)
	}

	fmt.Println("Done.")
}
