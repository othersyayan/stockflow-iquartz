package services

import (
	"backend/models"
	"backend/repositories"
)

type ProductService struct {
	productRepo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{productRepo: repo}
}

func (s *ProductService) GetProducts(search string, page, limit int) ([]models.Product, int, error) {
	return s.productRepo.GetProducts(search, page, limit)
}
