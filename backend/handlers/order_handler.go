package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"backend/models"
	"backend/repositories"
	"backend/services"
)

type OrderHandler struct {
	service    *services.OrderService
	repository *repositories.OrderRepository
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) HandleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateOrder(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "METHOD_NOT_ALLOWED", Message: "Method not allowed"},
		})
	}
}

func (h *OrderHandler) HandleOrderById(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetOrder(w, r)
	case http.MethodPost:
		h.CancelOrder(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "METHOD_NOT_ALLOWED", Message: "Method not allowed"},
		})
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req models.CreateOrderReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Invalid request body"},
		})
		return
	}

	if strings.TrimSpace(req.Customer) == "" || len(req.Items) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Customer and non empty items are required"},
		})
		return
	}

	skuMap := make(map[string]bool)

	for _, item := range req.Items {
		if item.SKU == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "SKU cannot be empty"},
			})
			return
		}

		if item.Qty <= 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Item Qty must be > 0"},
			})
			return
		}

		if skuMap[item.SKU] {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error: models.ErrorDetail{Code: "DUPLICATE_SKU_IN_ORDER", Message: "Duplicate SKUs in single order"},
			})
			return
		}

		skuMap[item.SKU] = true
	}

	order, errSKU, avail, err := h.repository.CreateOrder(req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		if err.Error() == "INVALID_SKU" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error: models.ErrorDetail{Code: "INVALID_SKU", Message: "SKU " + errSKU + " does not exist"},
			})
			return
		}

		if err == repositories.ErrInsufficientStock {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:      "INSUFFICIENT_STOCK",
					SKU:       errSKU,
					Available: &avail,
				},
			})
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderId := r.PathValue("id")

	if orderId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Order ID is required"},
		})
		return
	}

	id, errId := strconv.Atoi(orderId)

	if errId != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Invalid Order ID"},
		})
		return
	}

	err := h.repository.CancelOrder(id)

	if err == repositories.ErrNotFound {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "ORDER_NOT_FOUND", Message: "Order not found"},
		})
		return
	}

	if err != nil && err.Error() == "ALREADY_CANCELLED" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "ALREADY_CANCELLED", Message: "Order " + orderId + " is already cancelled"},
		})
		return
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "INTERNAL_SERVER_ERROR"},
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderId := r.PathValue("id")

	if orderId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Order ID is required"},
		})
		return
	}

	id, errId := strconv.Atoi(orderId)

	if errId != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Invalid Order ID"},
		})
		return
	}

	order, err := h.repository.GetOrder(id)
	if err == repositories.ErrNotFound {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "ORDER_NOT_FOUND", Message: "Order " + orderId + " is not found"},
		})
		return
	}
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "INTERNAL_SERVER_ERROR"},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
