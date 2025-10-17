package database

import "github.com/jackc/pgx/v5"

type Models struct {
	Users     UserModel
	Events    EventModel
	Attendees AttendeeModel
}

func NewModels(db *pgx.Conn) Models {
	return Models{
		Users:     UserModel{DB: db},
		Events:    EventModel{DB: db},
		Attendees: AttendeeModel{DB: db},
	}
}
