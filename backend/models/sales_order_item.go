package models

type SalesOrderItem struct {
	ID        int    `json:"id,omitempty"`
	OrderID   int    `json:"order_id,omitempty"`
	SKU       string `json:"sku"`
	Qty       int    `json:"qty"`
	UnitPrice int    `json:"unit_price"`
}

type CreateOrderItemReq struct {
	SKU string `json:"sku"`
	Qty int    `json:"qty"`
}
