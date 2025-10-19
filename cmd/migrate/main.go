package main

import (
	"api-in-gin/internal/env"
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib" // Important: register pgx driver
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide migration direction 'up' or 'down'")
	}

	log.Printf("ENV loaded: %s\n", env.GetEnvString("DATABASE_URL", ""))

	direction := os.Args[1]

	// Use "pgx" as driver name (matches the stdlib registration)
	db, err := sql.Open("pgx", env.GetEnvString("DATABASE_URL", "postgres://root:secret@localhost:5432/events?sslmode=disable"))
	if err != nil {
		log.Fatalf("error connecting to db: %v\n", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("error pinging database: %v\n", err)
	}
	log.Println("Database connection established")

	// Create pgx driver instance
	instance, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		log.Fatalf("error creating driver instance: %v\n", err)
	}

	// Open file source for migrations
	fSrc, err := (&file.File{}).Open("cmd/migrate/migrations")
	if err != nil {
		log.Fatalf("error opening migrations directory: %v\n", err)
	}

	// Create migrate instance - use "postgres" as database name
	m, err := migrate.NewWithInstance("file", fSrc, "postgres", instance)
	if err != nil {
		log.Fatalf("error creating migration instance: %v\n", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("error closing source: %v\n", srcErr)
		}
		if dbErr != nil {
			log.Printf("error closing database: %v\n", dbErr)
		}
	}()

	// Execute migration
	switch direction {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("error migrating db: %v\n", err)
		}
		if err == migrate.ErrNoChange {
			log.Println("No migrations to apply")
		} else {
			log.Println("Migrations applied successfully")
		}
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("error rolling back migration: %v\n", err)
		}
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No migrations to roll back")
		} else {
			log.Println("Migrations rolled back successfully")
		}
	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'")
	}
}
