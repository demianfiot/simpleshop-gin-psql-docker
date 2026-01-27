package service

import (
	"context"
	"prac/pkg/repository"
	"prac/pkg/service/events"
	"prac/todo"
	"time"
)

type OrderService struct {
	repo     repository.Order
	producer events.Producer
}

func NewOrderService(repo repository.Order, producer events.Producer) *OrderService {
	return &OrderService{repo: repo, producer: producer}
}
func (s *OrderService) CreateOrder(ctx context.Context, userID uint, items []todo.OrderItem) (int, error) {
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}

	order := todo.Order{
		UserID: userID,
		Status: "pending",
		Total:  total,
	}

	orderID, err := s.repo.CreateOrder(ctx, order, items)
	if err != nil {
		return 0, err
	}
	event := todo.OrderCreatedEvent{
		OrderID:   orderID,
		UserID:    userID,
		Total:     total,
		CreatedAt: time.Now(),
	}

	if err := s.producer.PublishOrderCreated(ctx, event); err != nil {
		//
	}
	return orderID, err
}

func (s *OrderService) GetUserOrders(ctx context.Context, userID uint) ([]todo.Order, error) {
	return s.repo.GetUserOrders(ctx, userID)
}

func (s *OrderService) GetAllOrders(ctx context.Context) ([]todo.Order, error) {
	return s.repo.GetAllOrders(ctx)
}

func (s *OrderService) GetOrderByID(ctx context.Context, orderID uint) (todo.Order, error) {
	return s.repo.GetOrderByID(ctx, orderID)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uint, status string) error {
	return s.repo.UpdateOrderStatus(ctx, orderID, status)
}
