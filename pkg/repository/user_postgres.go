package repository

import (
	"fmt"
	"prac/todo"

	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}
func (r *UserPostgres) CreateUser(user todo.User) (int, error) {
	query := `
	INSERT INTO users (name, email, password_hash, role)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`
	err := r.db.QueryRow(
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Role,
	).Scan(&user.ID)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return 0, ErrUserAlreadyExists
		}
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return int(user.ID), nil
}
func (r *UserPostgres) GetAllUsers() ([]todo.User, error) {
	var users []todo.User
	query := `SELECT id, name, email, created_at, role FROM users ORDER BY id`
	err := r.db.Select(&users, query)
	return users, err
}

func (r *UserPostgres) GetUserByID(id uint) (todo.User, error) {
	var user todo.User
	query :=
		`
		SELECT id, name, email, password_hash, role, created_at
		FROM users
		WHERE id = $1
	`
	err := r.db.Get(&user, query, id)
	return user, err
}
func (r *UserPostgres) GetUserByEmail(email string) (todo.User, error) {
	var user todo.User
	query :=
		`
		SELECT id, name, email, password_hash, role, is_active, created_at
		FROM users
		WHERE email = $1
	`
	err := r.db.Get(&user, query, email)
	return user, err
}
func (r *UserPostgres) UpdateUser(userid uint, user todo.UpdateUserInput) (todo.UpdateUserInput, error) {
	query := `
		UPDATE users
		SET name = $1, email = $2, role = $3
		WHERE id = $4
	`
	result, err := r.db.Exec(query, user.Name, user.Email, user.Role, userid)
	if err != nil {
		return todo.UpdateUserInput{}, err
	}
	// check zmin
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return todo.UpdateUserInput{}, err
	}
	if rowsAffected == 0 {
		return todo.UpdateUserInput{}, ErrProductNotFound
	}
	return user, err
}
func (r *UserPostgres) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrProductNotFound
	}
	return err
}
