package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"prac/pkg/repository"
	"prac/todo"
	"syscall"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/segmentio/kafka-go"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "order_events",
		GroupID: "analytics-service", // dl9 grupi - mashtabuvann9
	})
	defer reader.Close()
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "analytics",
			Password: "secret",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewOrderAnalyticsRepository(conn)

	// graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("shutting down consumer...")
		cancel()

	}()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Println("reader stopped:", err)
			return
		}

		var event todo.OrderCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("invalid message:", err)
			continue
		}
		if err := repo.InsertOrder(ctx, event); err != nil {
			log.Println("failed to insert into clickhouse:", err)
		}
		log.Printf(
			"Order %d created by user %d, total %.2f",
			event.OrderID,
			event.UserID,
			event.Total,
		)
	}
}
