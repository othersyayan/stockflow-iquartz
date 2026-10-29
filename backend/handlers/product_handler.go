package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"backend/models"
	"backend/repositories"
	"backend/services"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetProducts(w, r)
	case http.MethodPost:
		h.CreateProduct(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "METHOD_NOT_ALLOWED", Message: "Method not allowed"},
		})
	}
}

func (h *ProductHandler) HandleProductById(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		h.DeleteProduct(w, r)
	case http.MethodPut:
		h.UpdateProduct(w, r)
	case http.MethodPost:
		h.AdjustStock(w, r)
	case http.MethodGet:
		h.GetMovements(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "METHOD_NOT_ALLOWED", Message: "Method not allowed"},
		})
	}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}

	products, total, err := h.service.GetProducts(search, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  products,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProductReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Invalid payload"},
		})
		return
	}

	if strings.TrimSpace(req.SKU) == "" || strings.TrimSpace(req.Name) == "" || req.Price < 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "SKU and Name are required, Price must be >= 0"},
		})
		return
	}

	p, err := h.service.CreateProduct(req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "DUPLICATE_SKU", Message: "SKU already exists"},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productId := r.PathValue("id")

	if productId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Product ID is required"},
		})
		return
	}

	id, errId := strconv.Atoi(productId)

	if errId != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Invalid Product ID"},
		})
		return
	}

	err := h.service.DeleteProduct(id)

	if err == repositories.ErrNotFound {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "PRODUCT_NOT_FOUND", Message: "Product not found"},
		})
		return
	}

	if err == repositories.ErrProductHasHistory {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "PRODUCT_HAS_HISTORY", Message: "Product already has history movement"},
		})
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productId := r.PathValue("id")

	if productId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Product ID is required"},
		})
		return
	}

	id, errId := strconv.Atoi(productId)

	if errId != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Invalid Product ID"},
		})
		return
	}

	var req models.UpdateProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Invalid payload"},
		})
		return
	}

	p, err := h.service.UpdateProduct(id, req)
	if err == repositories.ErrNotFound {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "PRODUCT_NOT_FOUND", Message: "Product not found"},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *ProductHandler) AdjustStock(w http.ResponseWriter, r *http.Request) {
	productId := r.PathValue("id")

	if productId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Product ID is required"},
		})
		return
	}

	id, errId := strconv.Atoi(productId)

	if errId != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Invalid Product ID"},
		})
		return
	}

	var req models.AdjustStockReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || (req.Type != "in" && req.Type != "out") || req.Qty <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Invalid stock payload"},
		})
		return
	}

	p, err := h.service.AdjustStock(id, req)
	if err == repositories.ErrNotFound {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "PRODUCT_NOT_FOUND", Message: "Product not found"},
		})
		return
	}
	if err != nil && strings.Contains(err.Error(), "INSUFFICIENT_STOCK") {
		var available int
		fmt.Sscanf(err.Error(), "INSUFFICIENT_STOCK: available %d", &available)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "INSUFFICIENT_STOCK", Available: &available},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *ProductHandler) GetMovements(w http.ResponseWriter, r *http.Request) {
	productId := r.PathValue("id")

	if productId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Product ID is required"},
		})
		return
	}

	id, errId := strconv.Atoi(productId)

	if errId != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "VALIDATION_ERROR", Message: "Invalid Product ID"},
		})
		return
	}

	m, err := h.service.GetMovements(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: models.ErrorDetail{Code: "INTERNAL_SERVER_ERROR"},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": m,
	})
}
