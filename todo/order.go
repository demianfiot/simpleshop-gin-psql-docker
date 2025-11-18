package todo

import "time"

type Order struct {
	ID        uint      `json:"id" gorm:"primaryKey" db:"id"`
	UserID    uint      `json:"user_id" db:"user_id"`
	Status    string    `json:"status" db:"status"` // pending, paid, shipped, delivered
	Total     float64   `json:"total" db:"total"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	// UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
	Items    []OrderItem `json:"items" gorm:"foreignKey:OrderID" db:"-"`
	UserName string      `json:"user_name,omitempty" db:"user_name"` // + !!
}
