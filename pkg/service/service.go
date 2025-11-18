package service

import (
	"context"
	"prac/pkg/handler/dto"
	"prac/pkg/repository"
	"prac/todo"
)

type Autorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(token string) (uint, string, error)
	GetUserByID(userID uint) (todo.User, error)
}
type User interface {
	CreateUser(ctx context.Context, user dto.CreateUserInput) (int, error)
	GetAllUsers(ctx context.Context) ([]todo.User, error)
	GetUserByID(ctx context.Context, userID uint) (todo.User, error)
	GetUserByEmail(ctx context.Context, email string) (todo.User, error)
	UpdateUser(ctx context.Context, userID uint, input todo.UpdateUserInput) (todo.UpdateUserInput, error)
	DeleteUser(ctx context.Context, id int) error
}

type Product interface {
	CreateProduct(input todo.CreateProductInput, sellerID uint) (int, error)
	GetAllProducts() ([]todo.Product, error)
	GetProductByID(productID uint) (todo.Product, error)
	UpdateProduct(productID uint, input todo.UpdateProductInput, currentUserID uint) (todo.Product, error)
	DeleteProduct(id int) error
}
type Order interface {
	CreateOrder(userID uint, items []todo.OrderItem) (int, error)
	GetUserOrders(userID uint) ([]todo.Order, error)
	GetOrderByID(orderID uint) (todo.Order, error)
	UpdateOrderStatus(orderID uint, status string) error
	GetAllOrders() ([]todo.Order, error)
}

type Service struct {
	Autorization
	User
	Product
	Order
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Autorization: NewAuthService(repos.Autorization),
		User:         NewUserService(repos.User, repos.CacheRepository),
		Product:      NewProductService(repos.Product),
		Order:        NewOrderService(repos.Order),
	}
}
