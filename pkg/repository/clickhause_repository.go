package repository

import (
	"context"
	"prac/todo"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type OrderAnalyticsRepository struct {
	conn clickhouse.Conn
}

func NewOrderAnalyticsRepository(conn clickhouse.Conn) *OrderAnalyticsRepository {
	return &OrderAnalyticsRepository{conn: conn}
}

func (r *OrderAnalyticsRepository) InsertOrder(
	ctx context.Context,
	event todo.OrderCreatedEvent,
) error {

	return r.conn.Exec(ctx,
		`
		INSERT INTO orders_analytics
		(order_id, user_id, total, created_at)
		VALUES (?, ?, ?, now())
		`,
		event.OrderID,
		event.UserID,
		event.Total,
	)
}

// docker exec -it clickhouse clickhouse-client \
//   --query "$(cat migrations/clickhouse/001_create_orders_analytics.sql)"
