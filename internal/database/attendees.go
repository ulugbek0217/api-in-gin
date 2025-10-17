package database

import "github.com/jackc/pgx/v5"

type AttendeeModel struct {
	DB *pgx.Conn
}

type Attendee struct {
	ID     int `json:"id"`
	UserID int `json:"userID" binding:"required"`
}
