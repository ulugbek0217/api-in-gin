package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type AttendeeModel struct {
	DB *pgx.Conn
}

type Attendee struct {
	ID      int `json:"id"`
	EventID int `json:"eventID" binding:"required"`
	UserID  int `json:"userID" binding:"required"`
}

func (m *AttendeeModel) Insert(attendee *Attendee) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO attendees (event_id, user_id) VALUES ($1, $2) RETURNING id"
	err := m.DB.QueryRow(ctx, query, attendee.EventID, attendee.UserID).Scan(&attendee.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m *AttendeeModel) GetByEventAndAttendee(eventID, attendeeID int) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var attendee Attendee
	query := "SELECT * FROM attendees WHERE event_id = $1 AND user_id = $2"
	err := m.DB.QueryRow(ctx, query, eventID, attendeeID).Scan(&attendee.ID, &attendee.EventID, &attendee.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &attendee, nil
}

func (m *AttendeeModel) GetAttendeesByEvent(eventID int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT u.id, u.email, u.name FROM users u
	JOIN attendees a ON u.id = a.user_id
	WHERE a.event_id = $1
	`
	rows, err := m.DB.Query(ctx, query, eventID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Email, &user.Name)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (m *AttendeeModel) Delete(eventID, attendeeID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM attendees WHERE event_id = $1 AND user_id = $2"
	_, err := m.DB.Exec(ctx, query, eventID, attendeeID)
	if err != nil {
		return err
	}

	return nil
}

func (m *AttendeeModel) GetEventsByAttendee(id int) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT e.id, e.owner_id, e.name, e.description, e.date, e.location FROM events e
	JOIN attendees a ON e.id = a.event_id
	AND a.user_id = $1
	`
	rows, err := m.DB.Query(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		event := Event{}
		err := rows.Scan(&event.ID, &event.OwnerID, &event.Name, &event.Description, &event.Date, &event.Location)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}
