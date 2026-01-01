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

type Authorization interface {
	CreateUser(ctx context.Context, user todo.User) (int, error)
	GetUser(ctx context.Context, email, password string) (todo.User, error)
	GetUserByID(ctx context.Context, userID uint) (todo.User, error)
}
type User interface {
	CreateUser(ctx context.Context, user todo.User) (int, error)
	GetAllUsers(ctx context.Context) ([]todo.User, error)
	GetUserByID(ctx context.Context, id uint) (todo.User, error)
	GetUserByEmail(ctx context.Context, email string) (todo.User, error)
	UpdateUser(ctx context.Context, userid uint, user todo.UpdateUserInput) (todo.UpdateUserInput, error)
	DeleteUser(ctx context.Context, id int) error
}

type Product interface {
	CreateProduct(ctx context.Context, product todo.Product) (int, error)
	GetAllProducts(ctx context.Context) ([]todo.Product, error)
	GetProductByID(ctx context.Context, productID uint) (todo.Product, error)
	UpdateProduct(ctx context.Context, productID uint, product todo.Product) (todo.Product, error)
	DeleteProduct(ctx context.Context, id int) error
}
type Order interface {
	CreateOrder(ctx context.Context, order todo.Order, items []todo.OrderItem) (int, error)
	GetUserOrders(ctx context.Context, userID uint) ([]todo.Order, error)
	GetOrderByID(ctx context.Context, orderID uint) (todo.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uint, status string) error
	GetAllOrders(ctx context.Context) ([]todo.Order, error)
}

type Repository struct {
	CacheRepository
	Authorization
	User
	Product
	Order
}

func NewRepository(db *sqlx.DB, redisClient *redis.Client) *Repository {
	return &Repository{
		CacheRepository: NewCacheRepository(redisClient),
		Authorization:   NewAuthPostgres(db),
		User:            NewUserPostgres(db),
		Product:         NewProductPostgres(db),
		Order:           NewOrderPostgres(db),
	}
}
