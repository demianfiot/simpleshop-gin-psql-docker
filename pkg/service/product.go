package service

import (
	"prac/pkg/repository"
	"prac/todo"
)

type ProductService struct {
	repo repository.Product
}

func NewProductService(repo repository.Product) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(input todo.CreateProductInput, sellerID uint) (int, error) {
	product := todo.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		Category:    input.Category,
		SellerID:    sellerID,
	}

	return s.repo.CreateProduct(product)
}

func (s *ProductService) GetAllProducts() ([]todo.Product, error) {
	return s.repo.GetAllProducts()
}

func (s *ProductService) GetProductByID(userID uint) (todo.Product, error) {
	return s.repo.GetProductByID(userID)
}
func (s *ProductService) UpdateProduct(productID uint, input todo.UpdateProductInput, currentUserID uint) (todo.Product, error) {
	product := todo.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		Category:    input.Category,
	}

	return s.repo.UpdateProduct(productID, product, currentUserID)
}
func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.DeleteProduct(id)
}
