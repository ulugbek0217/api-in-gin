package main

import (
	"api-in-gin/internal/database"
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("error connecting to database: %v\n", err)
	}
	defer db.Close(context.Background())

	models := database.NewModels(db)
	app := &application{
		port:      os.Getenv("PORT"),
		jwtSecret: os.Getenv("JWT_SECRET"),
	}
}
