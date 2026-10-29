package services

import (
	"backend/models"
	"backend/repositories"
)

type OrderService struct {
	orderRepo *repositories.OrderRepository
}

func NewOrderService(orderRepo *repositories.OrderRepository) *OrderService {
	return &OrderService{orderRepo: orderRepo}
}

func (s *OrderService) CreateOrder(req models.CreateOrderReq) (*models.SalesOrder, string, int, error) {
	return s.orderRepo.CreateOrder(req)
}

func (s *OrderService) CancelOrder(orderId int) error {
	return s.orderRepo.CancelOrder(orderId)
}

func (s *OrderService) GetOrder(orderId int) (*models.SalesOrder, error) {
	return s.orderRepo.GetOrder(orderId)
}
