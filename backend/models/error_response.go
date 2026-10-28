package models

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code      string `json:"code"`
	Message   string `json:"message,omitempty"`
	SKU       string `json:"sku,omitempty"`
	Available *int   `json:"available,omitempty"`
}
