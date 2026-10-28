package services

import (
	"errors"

	"backend/models"
	"backend/repositories"
)

var (
	ErrInvalidRequest = errors.New("INVALID_REQUEST")
	ErrDuplicateSKU   = errors.New("DUPLICATE_SKU")
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

func (s *ProductService) CreateProduct(req models.CreateProductReq) (*models.Product, error) {
	if req.Name == "" || req.Price <= 0 {
		return nil, ErrInvalidRequest
	}

	products, _, err := s.productRepo.GetProducts(req.SKU, 1, 1)
	if err == nil && len(products) > 0 {
		return nil, ErrDuplicateSKU
	}
	return s.productRepo.CreateProduct(req)
}

func (s *ProductService) DeleteProduct(id int) error {
	err := s.productRepo.DeleteProduct(id)

	return err
}
