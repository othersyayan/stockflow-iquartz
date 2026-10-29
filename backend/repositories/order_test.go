package repositories_test

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"

	"backend/models"
	"backend/repositories"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	dbURL := "postgresql://root:postgres@localhost:5432/stockflow_iquartz_test?sslmode=disable"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test db: %v", err)
	}

	_, err = db.Exec(`
		TRUNCATE TABLE sales_order_items, sales_orders, stock_movements, products RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		t.Fatalf("Failed to truncate tables: %v", err)
	}

	return db
}

func TestCreateOrder_InsufficientStock(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	prodRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	_, err := prodRepo.CreateProduct(models.CreateProductReq{
		SKU:   "TEST-SKU-1",
		Name:  "Test Item",
		Price: 10000,
	})
	if err != nil {
		t.Fatalf("Failed to seed product: %v", err)
	}

	_, err = prodRepo.AdjustStock(1, models.AdjustStockReq{
		Type: "in",
		Qty:  5,
		Note: "Initial stock",
	})
	if err != nil {
		t.Fatalf("Failed to adjust initial stock: %v", err)
	}

	orderReq := models.CreateOrderReq{
		Customer: "Budi",
		Items: []models.CreateOrderItemReq{
			{SKU: "TEST-SKU-1", Qty: 10},
		},
	}

	_, _, _, err = orderRepo.CreateOrder(orderReq)

	if err == nil || err != repositories.ErrInsufficientStock {
		t.Errorf("Expected ErrInsufficientStock, got %v", err)
	}

	var currentStock int
	err = db.QueryRow("SELECT stock FROM products WHERE sku = 'TEST-SKU-1'").Scan(&currentStock)
	if err != nil {
		t.Fatalf("Failed to query product stock: %v", err)
	}

	if currentStock != 5 {
		t.Errorf("Expected stock to remain 5, but got %d (atomic transaction failed!)", currentStock)
	}
}

func TestCreateOrder_ConcurrencyRace(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	prodRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	_, err := prodRepo.CreateProduct(models.CreateProductReq{
		SKU:   "RACE-SKU",
		Name:  "Race Item",
		Price: 5000,
	})
	if err != nil {
		t.Fatalf("Failed to seed product: %v", err)
	}

	_, err = prodRepo.AdjustStock(1, models.AdjustStockReq{
		Type: "in",
		Qty:  10,
		Note: "Initial stock for race test",
	})
	if err != nil {
		t.Fatalf("Failed to adjust initial stock: %v", err)
	}

	var wg sync.WaitGroup
	concurrencyCount := 10
	successCount := 0
	failCount := 0
	var mu sync.Mutex

	for i := 0; i < concurrencyCount; i++ {
		wg.Add(1)
		go func(customerID int) {
			defer wg.Done()
			req := models.CreateOrderReq{
				Customer: fmt.Sprintf("Customer-%d", customerID),
				Items: []models.CreateOrderItemReq{
					{SKU: "RACE-SKU", Qty: 2},
				},
			}
			_, _, _, err := orderRepo.CreateOrder(req)

			mu.Lock()
			defer mu.Unlock()
			switch err {
			case nil:
				successCount++
			case repositories.ErrInsufficientStock:
				failCount++
			}
		}(i)
	}

	wg.Wait()

	var finalStock int
	err = db.QueryRow("SELECT stock FROM products WHERE sku = 'RACE-SKU'").Scan(&finalStock)
	if err != nil {
		t.Fatalf("Failed to query final stock: %v", err)
	}

	if finalStock < 0 {
		t.Errorf("Race condition detected! Final stock went negative: %d", finalStock)
	}

	if successCount > 5 {
		t.Errorf("Expected at most 5 successful orders, got %d", successCount)
	}

	t.Logf("Concurrency Test Results -> Success Orders: %d, Rejected (Insufficient): %d, Final Stock: %d", successCount, failCount, finalStock)
}
