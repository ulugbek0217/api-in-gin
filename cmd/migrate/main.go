package main

import (
	"api-in-gin/internal/env"
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib" // Important: register pgx driver
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide migration direction 'up' or 'down'")
	}

	log.Printf("ENV loaded: %s\n", env.GetEnvString("DATABASE_URL", "postgresql://root:secret@localhost:5432/events?sslmode=disable"))

	direction := os.Args[1]

	// Use "pgx" as driver name (matches the stdlib registration)
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
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

	//// Get current version
	//version, dirty, err := m.Version()
	//if err != nil && err != migrate.ErrNilVersion {
	//	log.Printf("error getting version: %v\n", err)
	//} else if err == migrate.ErrNilVersion {
	//	log.Println("No migrations applied yet")
	//} else {
	//	log.Printf("Current version: %d (dirty: %v)\n", version, dirty)
	//}

	// Execute migration
	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("error migrating db: %v\n", err)
		}
		if err == migrate.ErrNoChange {
			log.Println("No migrations to apply")
		} else {
			log.Println("Migrations applied successfully")
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("error rolling back migration: %v\n", err)
		}
		if err == migrate.ErrNoChange {
			log.Println("No migrations to roll back")
		} else {
			log.Println("Migrations rolled back successfully")
		}
	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'")
	}

	//// Get new version
	//newVersion, _, err := m.Version()
	//if err != nil && err != migrate.ErrNilVersion {
	//	log.Printf("error getting new version: %v\n", err)
	//} else if err != migrate.ErrNilVersion {
	//	log.Printf("New version: %d\n", newVersion)
	//}
}
