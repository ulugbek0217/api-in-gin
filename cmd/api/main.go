package main

import (
	"api-in-gin/internal/database"
	"api-in-gin/internal/env"
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
)

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	db, err := pgx.Connect(context.Background(), env.GetEnvString("DATABASE_URL", "postgresql://root:secret@localhost:5432/events?sslmode=disable"))
	if err != nil {
		log.Fatalf("error connecting to database: %v\n", err)
	}
	defer db.Close(context.Background())

	models := database.NewModels(db)
	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-123456"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}
