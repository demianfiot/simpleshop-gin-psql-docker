package todo

import "time"

type OrderItem struct {
	ID          uint      `json:"id" db:"id"`
	OrderID     uint      `json:"order_id" db:"order_id"`
	ProductID   uint      `json:"product_id" db:"product_id"`
	Quantity    int       `json:"quantity" db:"quantity"`
	Price       float64   `json:"price" db:"price"`
	ProductName string    `json:"product_name" db:"product_name"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"` // ← ДОДАТИ!
}
