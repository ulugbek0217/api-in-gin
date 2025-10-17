package database

import (
	"github.com/jackc/pgx/v5"
)

type UserModel struct {
	DB *pgx.Conn
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=3,max=64"`
	Password string `json:"-"`
}
