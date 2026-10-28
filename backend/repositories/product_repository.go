package repositories

import (
	"database/sql"
	"errors"

	"backend/models"
)

var (
	ErrNotFound          = errors.New("PRODUCT_NOT_FOUND")
	ErrProductHasHistory = errors.New("PRODUCT_HAS_HISTORY")
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

func (r *ProductRepository) CreateProduct(req models.CreateProductReq) (*models.Product, error) {
	query := `
		INSERT INTO products (sku, name, price)
		VALUES ($1, $2, $3)
		RETURNING id, sku, name, price, stock, is_active, created_at
	`
	var p models.Product
	err := r.db.QueryRow(query, req.SKU, req.Name, req.Price).Scan(
		&p.ID, &p.SKU, &p.Name, &p.Price, &p.Stock, &p.IsActive, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) DeleteProduct(id int) error {
	var count int
	checkQuery := `SELECT COUNT(*) FROM stock_movements WHERE product_id = $1`

	err := r.db.QueryRow(checkQuery, id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return ErrProductHasHistory
	}

	res, err := r.db.Exec(`DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
