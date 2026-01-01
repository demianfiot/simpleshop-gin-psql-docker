package service

import (
	"context"
	"prac/pkg/repository"
	"prac/todo"
)

type Authorization interface {
	CreateUser(ctx context.Context, user todo.User) (int, error)
	GenerateToken(ctx context.Context, email, password string) (string, error)
	ParseToken(ctx context.Context, token string) (uint, string, error)
	GetUserByID(ctx context.Context, userID uint) (todo.User, error)
}
type User interface {
	CreateUser(ctx context.Context, user todo.User) (int, error)
	GetAllUsers(ctx context.Context) ([]todo.User, error)
	GetUserByID(ctx context.Context, userID uint) (todo.User, error)
	GetUserByEmail(ctx context.Context, email string) (todo.User, error)
	UpdateUser(ctx context.Context, userID uint, input todo.UpdateUserInput) (todo.UpdateUserInput, error)
	DeleteUser(ctx context.Context, id int) error
}

type Product interface {
	CreateProduct(ctx context.Context, input todo.CreateProductInput, sellerID uint) (int, error)
	GetAllProducts(ctx context.Context) ([]todo.Product, error)
	GetProductByID(ctx context.Context, productID uint) (todo.Product, error)
	UpdateProduct(ctx context.Context, productID uint, input todo.UpdateProductInput, sellerID uint) (todo.Product, error)
	DeleteProduct(ctx context.Context, id int) error
}
type Order interface {
	CreateOrder(ctx context.Context, userID uint, items []todo.OrderItem) (int, error)
	GetUserOrders(ctx context.Context, userID uint) ([]todo.Order, error)
	GetOrderByID(ctx context.Context, orderID uint) (todo.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uint, status string) error
	GetAllOrders(ctx context.Context) ([]todo.Order, error)
}

type Service struct {
	Authorization
	User
	Product
	Order
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		User:          NewUserService(repos.User, repos.CacheRepository),
		Product:       NewProductService(repos.Product, repos.CacheRepository),
		Order:         NewOrderService(repos.Order),
	}
}
