package models

import "time"

type SalesOrder struct {
	ID          int              `json:"id"`
	Customer    string           `json:"customer"`
	Status      string           `json:"status"`
	Total       int              `json:"total"`
	CreatedAt   time.Time        `json:"created_at"`
	CancelledAt *time.Time       `json:"cancelled_at,omitempty"`
	Items       []SalesOrderItem `json:"items,omitempty"`
}

type CreateOrderReq struct {
	Customer string               `json:"customer"`
	Items    []CreateOrderItemReq `json:"items"`
}
