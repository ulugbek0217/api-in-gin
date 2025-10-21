package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserModel struct {
	DB *pgx.Conn
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=2,max=64"`
	Password string `json:"-"`
}

func (m *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO users (email, name, password) VALUES ($1, $2, $3) RETURNING id"

	return m.DB.QueryRow(ctx, query, user.Email, user.Name, user.Password).Scan(&user.ID)
}

func (m *UserModel) getUser(query string, args ...interface{}) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := m.DB.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) Get(id int) (*User, error) {
	query := "SELECT * FROM users WHERE id = $1"

	return m.getUser(query, id)
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := "SELECT * FROM users WHERE email = $1"

	return m.getUser(query, email)
}
