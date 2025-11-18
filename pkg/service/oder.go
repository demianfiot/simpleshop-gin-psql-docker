package service

import (
	"prac/pkg/repository"
	"prac/todo"
)

type OrderService struct {
	repo repository.Order
}

func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{repo: repo}
}
func (s *OrderService) CreateOrder(userID uint, items []todo.OrderItem) (int, error) {
	// Розрахунок загальної суми
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}

	order := todo.Order{
		UserID: userID,
		Status: "pending",
		Total:  total,
	}

	return s.repo.CreateOrder(order, items)
}

func (s *OrderService) GetUserOrders(userID uint) ([]todo.Order, error) {
	return s.repo.GetUserOrders(userID)
}

func (s *OrderService) GetAllOrders() ([]todo.Order, error) {
	return s.repo.GetAllOrders()
}

func (s *OrderService) GetOrderByID(orderID uint) (todo.Order, error) {
	return s.repo.GetOrderByID(orderID)
}

func (s *OrderService) UpdateOrderStatus(orderID uint, status string) error {
	return s.repo.UpdateOrderStatus(orderID, status)
}
