package repository

import (
	"context"
	"database/sql"
	"fmt"
	"prac/todo"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(ctx context.Context, user todo.User) (int, error) {
	query := `
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	err := r.db.QueryRow(
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return 0, ErrUserAlreadyExists
		}
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return int(user.ID), nil
}
func (r *AuthPostgres) GetUser(ctx context.Context, email, password string) (todo.User, error) {
	var user todo.User
	query := `
		SELECT id, name, role, created_at
		FROM users
		WHERE email = $1 AND password_hash = $2
	`

	err := r.db.Get(&user, query, email, password)

	return user, err
}
func (r *AuthPostgres) GetUserByID(ctx context.Context, userID uint) (todo.User, error) {
	var user todo.User
	query := `
		SELECT id, name, email, password_hash, role, created_at 
		FROM users 
		WHERE id = $1
	`
	err := r.db.Get(&user, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return todo.User{}, ErrUserNotFound
		}
		return todo.User{}, err
	}
	return user, nil
}
