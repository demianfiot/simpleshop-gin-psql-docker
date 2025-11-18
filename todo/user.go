package todo

import "time"

type User struct {
	ID           uint      `json:"id" db:"id"`
	Name         string    `json:"name" db:"name" binding:"required"`
	Email        string    `json:"email" db:"email" binding:"required"`
	PasswordHash string    `json:"password" db:"password_hash" binding:"required"`
	Role         string    `json:"role" db:"role"` // customer, seller, admin
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type UpdateUserInput struct {
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Role  string `json:"role" db:"role"`
}
