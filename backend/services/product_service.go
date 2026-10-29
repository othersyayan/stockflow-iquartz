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

func (s *ProductService) CreateProduct(req models.CreateProductReq) (*models.Product, error) {
	products, _, err := s.productRepo.GetProducts(req.SKU, 1, 1)
	if err == nil && len(products) > 0 {
		return nil, repositories.ErrDuplicateSKU
	}

	return s.productRepo.CreateProduct(req)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.productRepo.DeleteProduct(id)
}

func (s *ProductService) UpdateProduct(id int, req models.UpdateProductReq) (*models.Product, error) {
	return s.productRepo.UpdateProduct(id, req)
}

func (s *ProductService) AdjustStock(id int, req models.AdjustStockReq) (*models.Product, error) {
	product, err := s.productRepo.AdjustStock(id, req)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetMovements(id int) ([]models.StockMovement, error) {
	return s.productRepo.GetMovements(id)
}
