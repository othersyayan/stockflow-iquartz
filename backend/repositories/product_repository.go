package repositories

import (
	"database/sql"
	"errors"

	"backend/models"
)

var (
	ErrNotFound          = errors.New("PRODUCT_NOT_FOUND")
	ErrProductHasHistory = errors.New("PRODUCT_HAS_HISTORY")
	ErrInsufficientStock = errors.New("INSUFFICIENT_STOCK")
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProducts(search string, page, limit int) ([]models.Product, int, error) {
	offset := (page - 1) * limit
	queryCount := `SELECT COUNT(*) FROM products WHERE LOWER(sku) LIKE LOWER($1) OR LOWER(name) LIKE LOWER($1)`

	searchParam := "%" + search + "%"
	var total int
	if err := r.db.QueryRow(queryCount, searchParam).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			p.id, p.sku, p.name, p.price, p.stock, p.is_active, p.created_at,
			MAX(sm.created_at) as last_movement_at
		FROM products p
		LEFT JOIN stock_movements sm ON p.id = sm.product_id
		WHERE LOWER(p.sku) LIKE LOWER($1) OR LOWER(p.name) LIKE LOWER($1)
		GROUP BY p.id
		ORDER BY p.id DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, searchParam, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var lastMove sql.NullTime
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Price, &p.Stock, &p.IsActive, &p.CreatedAt, &lastMove); err != nil {
			return nil, 0, err
		}
		if lastMove.Valid {
			p.LastMovementAt = &lastMove.Time
		}
		products = append(products, p)
	}

	return products, total, nil
}
