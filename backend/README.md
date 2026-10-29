# StockFlow Backend API

Backend RESTful API service for **StockFlow** (Inventory & Sales Mini ERP Module), built with **Go (Golang)** and **PostgreSQL**.

---

## Architecture & Tech Stack

- **Language:** Go (1.25 / 1.22+)
- **Router:** Go 1.22+ Standard Library `net/http` (`http.NewServeMux`) with path-value routing (`/api/products/{id}`)
- **Database:** PostgreSQL 16
- **Database Driver:** `github.com/lib/pq`
- **Config Management:** `github.com/spf13/viper` (supports `.env` and system environment variables)
- **Containerization:** Multi-stage Docker build (Golang Alpine builder -> minimal Alpine runtime)

### Layered Architecture

```
backend/
├── database/          # Database connection pool setup
├── handlers/          # HTTP controllers (request parsing, response serializing)
├── models/            # Domain entities, DTOs, and standardized error responses
├── repositories/      # Database queries, atomic transactions, row locking
├── services/          # Business logic and domain validations
├── migrations/        # SQL DDL schemas
├── seed.sql           # Initial demo dataset
├── main.go            # Dependency injection, middleware, routes, server bootstrap
└── Dockerfile         # Multi-stage production container image
```

---

## How to Run (Step by Step)

### Option 1: Run with Docker Compose (Recommended)

From the project root directory:

```bash
docker compose up --build backend postgres
```

The backend will be available at `http://localhost:3000`.

---

### Option 2: Run Locally (Bare Metal)

#### 1. Prerequisites

- **Go 1.22+** installed (`go version`)
- **PostgreSQL 16+** running locally

#### 2. Create Database & Run Schema

Connect to your local PostgreSQL instance:

```bash
psql -U postgres
```

Create the database:

```sql
CREATE DATABASE stockflow;
```

Execute the schema migration and initial seed data:

```bash
psql -U postgres -d stockflow -f migrations/000001_init_schema.sql
psql -U postgres -d stockflow -f seed.sql
```

#### 3. Configure Environment Variables

Copy `.env.example` to `.env`:

```bash
cp .env.example .env
```

Ensure `.env` matches your local PostgreSQL credentials:

```env
DATABASE_URL=postgresql://[user]:[password]@localhost:5432/[db_name]?sslmode=disable
PORT=3000
```

#### 4. Install Dependencies

```bash
go mod download
```

#### 5. Run the Server

```bash
go run main.go
```

Expected output:

```
Connected to database successfully
Server running di localhost:3000
```

---

## Running Tests & Concurrency Validation

The repository includes unit tests and an automated **concurrency race-condition test** that runs 10 parallel goroutines attempting to purchase the same limited stock item to ensure that stock never goes below zero:

1. Create a dedicated test database (or ensure test DB is accessible):

```sql
CREATE DATABASE stockflow_test;
```

Apply the schema:

```bash
psql -U postgres -d stockflow_test -f migrations/000001_init_schema.sql
```

2. Run tests:

```bash
go test -v ./repositories
```

---

## API Reference & cURL Examples

Base URL: `http://localhost:3000`

### 1. Health Check

```bash
curl -X GET http://localhost:3000/health
```

**Response (200 OK):**

```json
{
  "message": "API Running",
  "status": "OK"
}
```

---

### 2. Products API

#### A. List Products (with Search & Pagination)

```bash
# Get page 1 with limit 10
curl -X GET "http://localhost:3000/api/products?page=1&limit=10"

# Search by SKU or Name
curl -X GET "http://localhost:3000/api/products?search=kopi&page=1&limit=10"
```

**Response (200 OK):**

```json
{
  "data": [
    {
      "id": 1,
      "sku": "SKU-001",
      "name": "Kopi Arabika 250g",
      "price": 85000,
      "stock": 50,
      "is_active": true,
      "created_at": "2026-10-28T14:41:00Z",
      "last_movement_at": "2026-10-29T10:00:00Z"
    }
  ],
  "limit": 10,
  "page": 1,
  "total": 12
}
```

#### B. Create Product

```bash
curl -X POST http://localhost:3000/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "sku": "SKU-099",
    "name": "Kopi Robusta Lampung 250g",
    "price": 45000
  }'
```

**Response (201 Created):**

```json
{
  "id": 13,
  "sku": "SKU-099",
  "name": "Kopi Robusta Lampung 250g",
  "price": 45000,
  "stock": 0,
  "is_active": true,
  "created_at": "2026-10-29T21:00:00Z",
  "last_movement_at": null
}
```

#### C. Update Product

```bash
curl -X PUT http://localhost:3000/api/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Kopi Arabika Gayo Super 250g",
    "price": 90000,
    "is_active": true
  }'
```

**Response (200 OK):**

```json
{
  "id": 1,
  "sku": "SKU-001",
  "name": "Kopi Arabika Gayo Super 250g",
  "price": 90000,
  "stock": 50,
  "is_active": true,
  "created_at": "2026-10-28T14:41:00Z",
  "last_movement_at": null
}
```

#### D. Adjust Stock (Stock In / Stock Out)

```bash
# Stock In (Increase stock)
curl -X POST http://localhost:3000/api/products/1/stock \
  -H "Content-Type: application/json" \
  -d '{
    "type": "in",
    "qty": 20,
    "note": "Restock from distributor"
  }'

# Stock Out (Decrease stock)
curl -X POST http://localhost:3000/api/products/1/stock \
  -H "Content-Type: application/json" \
  -d '{
    "type": "out",
    "qty": 5,
    "note": "Damaged goods write-off"
  }'
```

**Response (200 OK):**

```json
{
  "id": 1,
  "sku": "SKU-001",
  "name": "Kopi Arabika 250g",
  "price": 85000,
  "stock": 65,
  "is_active": true,
  "created_at": "2026-10-28T14:41:00Z",
  "last_movement_at": null
}
```

**Error Response when Stock Insufficient (409 Conflict):**

```json
{
  "error": {
    "code": "INSUFFICIENT_STOCK",
    "available": 10
  }
}
```

#### E. Get Product Stock Movement History

```bash
curl -X GET http://localhost:3000/api/products/1/movements
```

**Response (200 OK):**

```json
{
  "data": [
    {
      "id": 3,
      "product_id": 1,
      "type": "out",
      "qty": 5,
      "note": "Damaged goods write-off",
      "created_at": "2026-10-29T21:10:00Z"
    },
    {
      "id": 2,
      "product_id": 1,
      "type": "in",
      "qty": 20,
      "note": "Restock from distributor",
      "created_at": "2026-10-29T21:05:00Z"
    }
  ]
}
```

#### F. Delete Product

```bash
curl -X DELETE http://localhost:3000/api/products/13
```

- Returns `200 OK` if deleted successfully.
- Returns `409 Conflict` if the product has existing stock movements:

```json
{
  "error": {
    "code": "PRODUCT_HAS_HISTORY",
    "message": "Product already has history movement"
  }
}
```

---

### 3. Sales Orders API

#### A. Create Sales Order (Atomic Transaction)

```bash
curl -X POST http://localhost:3000/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer": "Toko Berkah Sejahtera",
    "items": [
      {
        "sku": "SKU-001",
        "qty": 2
      },
      {
        "sku": "SKU-002",
        "qty": 5
      }
    ]
  }'
```

**Response (201 Created):**

```json
{
  "id": 1,
  "customer": "Toko Berkah Sejahtera",
  "status": "CREATED",
  "total": 260000,
  "created_at": "2026-10-29T21:30:00Z",
  "items": [
    {
      "id": 1,
      "order_id": 1,
      "sku": "SKU-001",
      "qty": 2,
      "unit_price": 85000
    },
    {
      "id": 2,
      "order_id": 1,
      "sku": "SKU-002",
      "qty": 5,
      "unit_price": 18000
    }
  ]
}
```

#### B. Get Sales Order by ID

```bash
curl -X GET http://localhost:3000/api/orders/1
```

**Response (200 OK):**

```json
{
  "id": 1,
  "customer": "Toko Berkah Sejahtera",
  "status": "CREATED",
  "total": 260000,
  "created_at": "2026-10-29T21:30:00Z",
  "items": [
    {
      "id": 1,
      "order_id": 1,
      "sku": "SKU-001",
      "qty": 2,
      "unit_price": 85000
    },
    {
      "id": 2,
      "order_id": 1,
      "sku": "SKU-002",
      "qty": 5,
      "unit_price": 18000
    }
  ]
}
```

#### C. Cancel Sales Order (Rollback & Restore Stock)

```bash
curl -X POST http://localhost:3000/api/orders/1/cancel
```

**Response:** `200 OK` (Stock is restored and stock movements are logged as `in`).

If attempted again:
**Response (409 Conflict):**

```json
{
  "error": {
    "code": "ALREADY_CANCELLED",
    "message": "Order 1 is already cancelled"
  }
}
```

---

## Design Decisions

1. **Standard Library `net/http` router (Go 1.22+):**
   - Kept the service lightweight and zero-dependency bloat by leveraging Go 1.22's enhanced router supporting methods and path parameters (`/api/products/{id}`).
2. **Pessimistic Row-Level Locking (`SELECT ... FOR UPDATE`):**
   - During order checkout and stock decrement, we acquire exclusive row locks on the affected product rows inside an active SQL transaction. This completely prevents race conditions and phantom reads where two concurrent shoppers could simultaneously buy the last available item.
3. **Database Constraints as the Ultimate Guardrail:**
   - The PostgreSQL schema enforces `CHECK (stock >= 0)`, `CHECK (price >= 0)`, and `CHECK (qty > 0)`. Even if an unexpected code bug occurred, the database itself guarantees negative stock can never exist.
4. **Append-Only Audit Trail (`stock_movements`):**
   - Stock values are not just mutated blindly; every mutation (manual adjustment, order placement, or order cancellation) logs a corresponding record in `stock_movements` with direction (`in`/`out`), quantity, and reason/note.
5. **Standardized Error Envelope:**
   - All errors follow a predictable JSON schema (`{"error": {"code": "...", "message": "...", "available": ...}}`) enabling frontends to display user-friendly contextual alerts (e.g. showing remaining available quantity).

---

## What I Would Improve

1. **Database Migration Tooling:** Integrate `golang-migrate` or `goose` instead of raw SQL scripts to track migration version history automatically.
2. **Structured Logging & Observability:** Replace standard `fmt` / `log` with `zerolog` or `zap` with correlation IDs and Prometheus metrics endpoints.
3. **Authentication & Authorization:** Add JWT middleware with role-based access control (RBAC), distinguishing between Warehouse Staff (stock adjustments) and Sales Reps (order placement).
4. **Idempotency Keys:** Support `Idempotency-Key` header on `/api/orders` to safely handle network retries from clients without duplicate orders.
5. **Swagger / OpenAPI Documentation:** Generate interactive documentation using `swag` / OpenAPI 3.0 specs.

---

## Honest Hours Spent

- **Database Design & Initial Setup:** 1.5 hours (Schema, migrations, constraints, seed dataset)
- **Product Module (CRUD & Stock Adjustments):** 3.5 hours (Repositories, validation, stock movements)
- **Sales Order Module & Concurrency Logic:** 3.5 hours (Multi-item transactions, row locking, rollback & cancellation)
- **Unit & Concurrency Tests:** 3.5 hours (Parallel goroutine testing for race condition validation)
- **Total Backend Effort:** **~12 hours**
