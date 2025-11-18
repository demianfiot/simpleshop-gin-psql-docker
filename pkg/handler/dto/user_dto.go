package dto

import "time"

type CreateUserInput struct {
	ID           uint      `json:"id" db:"id"`
	Name         string    `json:"name" db:"name" binding:"required"`
	Email        string    `json:"email" db:"email" binding:"required"`
	PasswordHash string    `json:"password" db:"password_hash" binding:"required"`
	Role         string    `json:"role" db:"role" binding:"required"` // customer, seller, admin
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// type SignUpInput struct { // not used
// 	Name     string `json:"name" binding:"required"`
// 	Email    string `json:"email" binding:"required,email"`
// 	Password string `json:"password" binding:"required,min=6"`
// 	// Role не включаємо - встановлюється автоматично
// }
