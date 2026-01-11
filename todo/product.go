package todo

import "time"

type Product struct {
	ID          uint      `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" binding:"required"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price" binding:"required"`
	Stock       int       `json:"stock" db:"stock" binding:"required"`
	Category    string    `json:"category" db:"category"`
	SellerID    uint      `json:"seller_id" db:"seller_id"` // +!!
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
type CreateProductInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"required,min=0"`
	Category    string  `json:"category"`
	// SellerID - отримуємо з контексту
}
type UpdateProductInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"required,min=0"`
	Category    string  `json:"category"`
	// SellerID - не можна змінювати власника
}
