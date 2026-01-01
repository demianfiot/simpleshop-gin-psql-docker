package repository

import (
	"context"
	"database/sql"
	"fmt"
	"prac/todo"

	"github.com/jmoiron/sqlx"
)

type ProductPostgres struct {
	db *sqlx.DB
}

func NewProductPostgres(db *sqlx.DB) *ProductPostgres {
	return &ProductPostgres{db: db}
}

func (r *ProductPostgres) CreateProduct(ctx context.Context, product todo.Product) (int, error) {
	query := `
	INSERT INTO products (name, description, price, stock, category, seller_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at
	`
	err := r.db.QueryRow(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.Category,
		product.SellerID,
	).Scan(&product.ID, &product.CreatedAt)
	if err != nil {
		if err.Error() == "duplicate key value violates unique constraint " {
			return 0, ErrProductAlreadyExists
		}
		return 0, fmt.Errorf("failed to create product: %w", err)
	}
	return int(product.ID), nil
}
func (r *ProductPostgres) GetAllProducts(ctx context.Context) ([]todo.Product, error) {
	var products []todo.Product
	query := `SELECT id, name, description, price, stock, category, seller_id FROM products ORDER BY id`
	err := r.db.Select(&products, query)
	return products, err
}
func (r *ProductPostgres) GetProductByID(ctx context.Context, productID uint) (todo.Product, error) {
	var product todo.Product
	query :=
		`
		SELECT id, name, description, price, stock, category, seller_id
		FROM products
		WHERE id = $1
	`
	err := r.db.Get(&product, query, productID)
	return product, err
}
func (r *ProductPostgres) UpdateProduct(ctx context.Context, productID uint, product todo.Product) (todo.Product, error) {
	query := `
        UPDATE products
        SET name = $1, description = $2, price = $3, stock = $4, category = $5
        WHERE id = $6 AND (seller_id = $7 OR EXISTS (SELECT 1 FROM users WHERE id = $7 AND role = 'admin'))
        RETURNING id, name, description, price, stock, category, seller_id, created_at
    `

	var updatedProduct todo.Product
	err := r.db.QueryRow(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.Category,
		productID,
		product.SellerID, // check seller / admin
	).Scan(
		&updatedProduct.ID,
		&updatedProduct.Name,
		&updatedProduct.Description,
		&updatedProduct.Price,
		&updatedProduct.Stock,
		&updatedProduct.Category,
		&updatedProduct.SellerID,
		&updatedProduct.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return todo.Product{}, ErrProductNotFound
		}
		return todo.Product{}, err
	}

	return updatedProduct, nil
}
func (r *ProductPostgres) DeleteProduct(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`
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
