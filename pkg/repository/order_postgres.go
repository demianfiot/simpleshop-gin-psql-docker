package repository

import (
	"context"
	"database/sql"
	"fmt"
	"prac/todo"

	"github.com/jmoiron/sqlx"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) CreateOrder(ctx context.Context, order todo.Order, items []todo.OrderItem) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// zamovlenn9
	var orderID int
	query := `
		INSERT INTO orders (user_id, status, total)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err = tx.QueryRow(
		query,
		order.UserID,
		order.Status,
		order.Total,
	).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}

	// + tovary
	itemsQuery := `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4)
	`
	for _, item := range items {
		_, err := tx.Exec(itemsQuery, orderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return 0, fmt.Errorf("failed to add order item: %w", err)
		}

		// onovlenn9 skladu tovaru
		updateStockQuery := `
			UPDATE products 
			SET stock = stock - $1 
			WHERE id = $2 AND stock >= $1
		`
		result, err := tx.Exec(updateStockQuery, item.Quantity, item.ProductID)
		if err != nil {
			return 0, fmt.Errorf("failed to update product stock: %w", err)
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return 0, fmt.Errorf("not enough stock for product %d", item.ProductID)
		}
	}

	// tranzakci9
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return orderID, nil
}

func (r *OrderPostgres) GetUserOrders(ctx context.Context, userID uint) ([]todo.Order, error) {
	var orders []todo.Order
	query := `
		SELECT id, user_id, status, total, created_at 
		FROM orders 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`
	err := r.db.Select(&orders, query, userID)
	return orders, err
}

func (r *OrderPostgres) GetOrderByID(ctx context.Context, orderID uint) (todo.Order, error) {
	var order todo.Order

	query := `
		SELECT id, user_id, status, total, created_at
		FROM orders
		WHERE id = $1
	`
	err := r.db.Get(&order, query, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return todo.Order{}, ErrOrderNotFound
		}
		return todo.Order{}, err
	}

	var items []todo.OrderItem
	itemsQuery := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price, p.name as product_name
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		WHERE oi.order_id = $1
	`
	err = r.db.Select(&items, itemsQuery, orderID)
	if err != nil {
		return todo.Order{}, err
	}

	order.Items = items
	return order, nil
}

func (r *OrderPostgres) UpdateOrderStatus(ctx context.Context, orderID uint, status string) error {
	query := `
		UPDATE orders
		SET status = $1
		WHERE id = $2
	`
	result, err := r.db.Exec(query, status, orderID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrOrderNotFound
	}

	return nil
}

func (r *OrderPostgres) GetAllOrders(ctx context.Context) ([]todo.Order, error) {
	var orders []todo.Order
	query := `
        SELECT o.id, o.user_id, o.status, o.total, o.created_at, u.name as user_name
        FROM orders o
        JOIN users u ON o.user_id = u.id
        ORDER BY o.created_at DESC
    `
	err := r.db.Select(&orders, query)
	return orders, err
}
