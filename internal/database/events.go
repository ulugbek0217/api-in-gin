package database

import (
	"time"

	"github.com/jackc/pgx/v5"
)

type EventModel struct {
	DB *pgx.Conn
}

type Event struct {
	ID          int       `json:"id"`
	OwnerID     string    `json:"ownerID" binding:"required"`
	Name        string    `json:"name" binding:"required,min=3"`
	Description string    `json:"description" binding:"required,min=10"`
	Date        time.Time `json:"date" binding:"required,datetime=2006-01-02"`
	Location    string    `json:"location" binding:"required,min=3"`
}
