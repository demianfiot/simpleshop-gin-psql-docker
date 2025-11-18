package repository

import (
	"context"
	"prac/todo"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	DeleteByPattern(ctx context.Context, pattern string) error
}

type Autorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(email, password string) (todo.User, error)
	GetUserByID(userID uint) (todo.User, error)
}

type User interface {
	CreateUser(user todo.User) (int, error)
	GetAllUsers() ([]todo.User, error)
	GetUserByID(id uint) (todo.User, error)
	GetUserByEmail(email string) (todo.User, error)
	UpdateUser(userid uint, user todo.UpdateUserInput) (todo.UpdateUserInput, error)
	DeleteUser(id int) error
}

type Product interface {
	CreateProduct(product todo.Product) (int, error)
	GetAllProducts() ([]todo.Product, error)
	GetProductByID(productID uint) (todo.Product, error)
	UpdateProduct(productID uint, product todo.Product, currentUserID uint) (todo.Product, error)
	DeleteProduct(id int) error
}
type Order interface {
	CreateOrder(order todo.Order, items []todo.OrderItem) (int, error)
	GetUserOrders(userID uint) ([]todo.Order, error)
	GetOrderByID(orderID uint) (todo.Order, error)
	UpdateOrderStatus(orderID uint, status string) error
	GetAllOrders() ([]todo.Order, error)
}

type Repository struct {
	CacheRepository
	Autorization
	User
	Product
	Order
}

func NewRepository(db *sqlx.DB, redisClient *redis.Client) *Repository {
	return &Repository{
		CacheRepository: NewCacheRepository(redisClient),
		Autorization:    NewAuthPostgres(db),
		User:            NewUserPostgres(db),
		Product:         NewProductPostgres(db),
		Order:           NewOrderPostgres(db),
	}
}
