package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"backend/models"
)

var (
	ErrAlreadyCancelled   = errors.New("ALREADY_CANCELLED")
	ErrInvalidSKU         = errors.New("INVALID_SKU")
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(req models.CreateOrderReq) (*models.SalesOrder, string, int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, "", 0, err
	}
	defer tx.Rollback()

	totalOrderAmount := 0
	type itemDetail struct {
		productID int
		sku       string
		qty       int
		unitPrice int
	}
	var itemsToProcess []itemDetail

	for _, item := range req.Items {
		var prodID, currentStock, price int
		queryLock := `SELECT id, stock, price FROM products WHERE sku = $1 FOR UPDATE`
		err := tx.QueryRow(queryLock, item.SKU).Scan(&prodID, &currentStock, &price)
		if err == sql.ErrNoRows {
			return nil, item.SKU, 0, ErrInvalidSKU
		} else if err != nil {
			return nil, "", 0, err
		}

		if currentStock < item.Qty {
			return nil, item.SKU, currentStock, ErrInsufficientStock
		}

		res, err := tx.Exec(`UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1`, item.Qty, prodID)
		if err != nil {
			return nil, "", 0, err
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return nil, item.SKU, currentStock, ErrInsufficientStock
		}

		_, err = tx.Exec(`INSERT INTO stock_movements (product_id, type, qty, note) VALUES ($1, 'out', $2, $3)`, prodID, item.Qty, "Sales Order Creation")
		if err != nil {
			return nil, "", 0, err
		}

		totalOrderAmount += price * item.Qty
		itemsToProcess = append(itemsToProcess, itemDetail{
			productID: prodID,
			sku:       item.SKU,
			qty:       item.Qty,
			unitPrice: price,
		})
	}

	var order models.SalesOrder
	queryOrder := `INSERT INTO sales_orders (customer, total) VALUES ($1, $2) RETURNING id, customer, status, total, created_at`

	err = tx.QueryRow(queryOrder, req.Customer, totalOrderAmount).Scan(&order.ID, &order.Customer, &order.Status, &order.Total, &order.CreatedAt)
	if err != nil {
		return nil, "", 0, err
	}

	for _, item := range itemsToProcess {
		var orderItem models.SalesOrderItem
		queryItem := `INSERT INTO sales_order_items (order_id, sku, qty, unit_price) VALUES ($1, $2, $3, $4) RETURNING id, order_id, sku, qty, unit_price`

		err := tx.QueryRow(queryItem, order.ID, item.sku, item.qty, item.unitPrice).Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.SKU, &orderItem.Qty, &orderItem.UnitPrice)
		if err != nil {
			return nil, "", 0, err
		}

		order.Items = append(order.Items, orderItem)
	}

	if err := tx.Commit(); err != nil {
		return nil, "", 0, err
	}

	return &order, "", 0, nil
}

func (r *OrderRepository) CancelOrder(orderId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var status string
	err = tx.QueryRow(`SELECT status FROM sales_orders WHERE id = $1 FOR UPDATE`, orderId).Scan(&status)

	if err == sql.ErrNoRows {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	if status == "CANCELLED" {
		return ErrAlreadyCancelled
	}

	rows, err := tx.Query(`SELECT sku, qty FROM sales_order_items WHERE order_id = $1`, orderId)
	if err != nil {
		return err
	}
	defer rows.Close()

	type restoreItem struct {
		sku string
		qty int
	}

	var items []restoreItem

	for rows.Next() {
		var item restoreItem
		if err := rows.Scan(&item.sku, &item.qty); err != nil {
			return err
		}
		items = append(items, item)
	}

	for _, item := range items {
		var prodID int
		err := tx.QueryRow(`UPDATE products SET stock = stock + $1 WHERE sku = $2 RETURNING id`, item.qty, item.sku).Scan(&prodID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`INSERT INTO stock_movements (product_id, type, qty, note) VALUES ($1, 'in', $2, $3)`, prodID, item.qty, fmt.Sprintf("Order #%d Cancellation", orderId))
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(`UPDATE sales_orders SET status = 'CANCELLED', cancelled_at = NOW() WHERE id = $1`, orderId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *OrderRepository) GetOrder(orderId int) (*models.SalesOrder, error) {
	var order models.SalesOrder
	var cancelledAt sql.NullTime

	query := `SELECT id, customer, status, total, created_at, cancelled_at FROM sales_orders WHERE id = $1`
	err := r.db.QueryRow(query, orderId).Scan(&order.ID, &order.Customer, &order.Status, &order.Total, &order.CreatedAt, &cancelledAt)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	if cancelledAt.Valid {
		order.CancelledAt = &cancelledAt.Time
	}

	rows, err := r.db.Query(`SELECT id, order_id, sku, qty, unit_price FROM sales_order_items WHERE order_id = $1`, orderId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var item models.SalesOrderItem

		if err := rows.Scan(&item.ID, &item.OrderID, &item.SKU, &item.Qty, &item.UnitPrice); err != nil {
			return nil, err
		}

		order.Items = append(order.Items, item)
	}

	return &order, nil
}
