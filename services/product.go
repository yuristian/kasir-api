package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	repository *repositories.ProductRepository
}

func NewProductService(repository *repositories.ProductRepository) *ProductService {
	return &ProductService{repository: repository}
}

func (s *ProductService) GetAllProducts() ([]*models.ProductWithCategory, error) {
	return s.repository.GetAllProducts()
}

func (s *ProductService) GetAllProductsByCategoryID(categoryID int) ([]*models.ProductWithCategory, error) {
	return s.repository.GetAllProductsByCategoryID(categoryID)
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	return s.repository.CreateProduct(product)
}

func (s *ProductService) GetProductByID(id int) (*models.ProductWithCategory, error) {
	return s.repository.GetProductByID(id)
}

func (s *ProductService) UpdateProduct(product *models.Product) error {
	return s.repository.UpdateProduct(product)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.repository.DeleteProduct(id)
}
