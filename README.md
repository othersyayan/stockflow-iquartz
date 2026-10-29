# StockFlow — Inventory & Sales Mini ERP

A robust, full-stack **Inventory & Sales Mini ERP** application built with **Go (Golang)**, **PostgreSQL**, **React 19**, **TypeScript**, and **Docker Compose**.

StockFlow is engineered with strict transactional integrity, concurrency-safe inventory deductions, audit trail tracking, and a clean, responsive modern dashboard.

---

## Repository Structure

```
.
├── backend/                  # Golang API backend
│   ├── database/             # PostgreSQL connection & pool setup
│   ├── handlers/             # HTTP controller layer
│   ├── models/               # Domain models, requests, responses
│   ├── repositories/         # Database persistence & transactional queries
│   │   └── order_test.go     # Unit and concurrency race-condition tests
│   ├── services/             # Business logic orchestration
│   ├── migrations/           # Database DDL schema files
│   ├── seed.sql              # Initial demo dataset
│   ├── Dockerfile            # Multi-stage production Go image
│   └── main.go               # Server bootstrap & router registration
├── frontend/                 # React 19 + TypeScript frontend
│   ├── src/
│   │   ├── api/              # Typed API client
│   │   ├── components/       # SummaryCards, ProductTable, AdjustStockModal
│   │   ├── types/            # TypeScript data contracts
│   │   └── App.tsx           # Dashboard orchestrator
│   ├── Dockerfile            # Multi-stage Nginx + Vite image
│   └── vite.config.ts        # Vite build & local reverse proxy
├── docker-compose.yml        # Orchestration for Database, Backend, and Frontend
└── README.md                 # Project root documentation
```

---

## How to Run (Step by Step)

### Option 1: Docker Compose (Quickest & Recommended)

Run the entire stack with a single command:

```bash
docker compose up --build
```

- **Frontend Dashboard:** [http://localhost:5173](http://localhost:5173)
- **Backend API:** [http://localhost:3000](http://localhost:3000)
- **PostgreSQL:** `localhost:5432` (`user: root`, `password: postgres`, `db: stockflow`)

To stop the services:

```bash
docker compose down
```

---

### Option 2: Local Manual Setup (Without Docker)

#### Step 1: Database Setup (PostgreSQL)

1. Ensure PostgreSQL is running on your machine:
   ```bash
   psql -U postgres
   ```
2. Create the application database:
   ```sql
   CREATE DATABASE stockflow;
   ```
3. Run the schema migrations and seed data:
   ```bash
   psql -U postgres -d stockflow -f backend/migrations/000001_init_schema.sql
   psql -U postgres -d stockflow -f backend/seed.sql
   ```

#### Step 2: Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Create your `.env` configuration:
   ```bash
   cp .env.example .env
   ```
   Adjust the credentials in `.env` if necessary:
   ```env
   DATABASE_URL=postgresql://[user]:[password]@localhost:5432/[db_name]?sslmode=disable
   PORT=3000
   ```
3. Download Go modules and start the API:
   ```bash
   go mod download
   go run main.go
   ```
   The backend will start on `http://localhost:3000`.

#### Step 3: Frontend Setup

1. Open a new terminal and navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install npm dependencies:
   ```bash
   npm install
   ```
3. Start the Vite development server:
   ```bash
   npm run dev
   ```
   Open [http://localhost:5173](http://localhost:5173) in your browser.

---

## Concurrency & Race-Condition Tests

To verify that StockFlow guarantees inventory safety during high concurrency traffic:

1. Create a dedicated test database:
   ```sql
   CREATE DATABASE stockflow_test;
   ```
2. Apply the schema:
   ```bash
   psql -U postgres -d stockflow_test -f backend/migrations/000001_init_schema.sql
   ```
3. Run the Go test suite:
   ```bash
   cd backend && go test -v ./repositories
   ```

**What this test validates:**

- **Race Condition Resistance:** 10 simultaneous goroutines concurrently order from limited stock (10 units).
- **Exact Deduction:** Exactly 5 orders succeed (each requesting 2 units), the remaining 5 fail cleanly with `INSUFFICIENT_STOCK`.
- **Zero Negative Stock:** Stock never drops below 0.
- **Rollback Consistency:** Failed orders leave stock intact.

---

## API Reference & cURL Examples

### 1. Health Check

```bash
curl -i -X GET http://localhost:3000/health
```

```json
{
  "status": "OK",
  "message": "API Running"
}
```

---

### 2. Products

#### Get Products (with Search & Pagination)

```bash
# Get page 1 with limit 10
curl -i -X GET "http://localhost:3000/api/products?page=1&limit=10"

# Search by SKU or Name
curl -i -X GET "http://localhost:3000/api/products?search=kopi"
```

#### Create New Product

```bash
curl -i -X POST http://localhost:3000/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "sku": "SKU-020",
    "name": "Kopi Susu Gula Aren 1L",
    "price": 65000
  }'
```

#### Update Product

```bash
curl -i -X PUT http://localhost:3000/api/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Kopi Arabika 250g (Specialty Grade)",
    "price": 95000,
    "is_active": true
  }'
```

#### Adjust Stock (Stock In / Stock Out)

```bash
# Stock In (Add inventory)
curl -i -X POST http://localhost:3000/api/products/1/stock \
  -H "Content-Type: application/json" \
  -d '{
    "type": "in",
    "qty": 25,
    "note": "Restock from roastery"
  }'

# Stock Out (Deduct inventory)
curl -i -X POST http://localhost:3000/api/products/1/stock \
  -H "Content-Type: application/json" \
  -d '{
    "type": "out",
    "qty": 5,
    "note": "Sampling & quality test"
  }'
```

#### Get Stock Movements History

```bash
curl -i -X GET http://localhost:3000/api/products/1/movements
```

#### Delete Product

```bash
curl -i -X DELETE http://localhost:3000/api/products/12
```

_Note: If the product already has stock movement records, deletion is blocked (`409 Conflict: PRODUCT_HAS_HISTORY`) to protect historical audit data._

---

### 3. Sales Orders

#### Create Sales Order (Atomic Multi-Item Checkout)

```bash
curl -i -X POST http://localhost:3000/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer": "Cafe Senja",
    "items": [
      {
        "sku": "SKU-001",
        "qty": 2
      },
      {
        "sku": "SKU-002",
        "qty": 3
      }
    ]
  }'
```

**Response (201 Created):**

```json
{
  "id": 1,
  "customer": "Cafe Senja",
  "status": "CREATED",
  "total": 224000,
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
      "qty": 3,
      "unit_price": 18000
    }
  ]
}
```

#### Get Sales Order by ID

```bash
curl -i -X GET http://localhost:3000/api/orders/1
```

#### Cancel Sales Order (Restores Stock Atomically)

```bash
curl -i -X POST http://localhost:3000/api/orders/1/cancel
```

_Restores deducted product quantities back into inventory and records compensatory `in` stock movements with reference to the cancelled order._

---

## Design Decisions

1. **Pessimistic Concurrency Control (`SELECT ... FOR UPDATE`):**
   - High-throughput flash sales and concurrent ordering cannot rely on optimistic checks without retry loops. Acquiring row locks inside a transactional block ensures sequential, deterministic stock deduction and avoids phantom inventory overselling.
2. **Database-Level Integrity Constraints:**
   - In addition to application validation, PostgreSQL `CHECK (stock >= 0)` and foreign key cascades serve as immutable defense-in-depth guarantees.
3. **Audit Trail Immutability (`stock_movements`):**
   - Stock quantities are never edited in isolation; every state change leaves a traceable ledger record with direction, timestamp, and context.
4. **Standard Library Routing (Go 1.22+):**
   - Eliminated third-party router frameworks (e.g., Gin, Fiber) in favor of the enhanced standard `http.ServeMux`, resulting in smaller binary sizes, zero maintenance dependencies, and faster startup.
5. **Modern Frontend with Tailwind CSS v4 & Native Fetch:**
   - Used React 19 and Tailwind CSS 4 for instant compilation and low overhead. Kept the API client native and type-safe without bundling bulky HTTP libraries.
6. **Containerization & Reverse Proxying:**
   - Configured Docker multi-stage builds for both Go and React. Used Nginx to serve the static frontend bundle while proxying `/api` calls directly to the backend container, eliminating CORS discrepancies across development and production.

---

## What I Would Improve

If given more time to advance StockFlow toward enterprise production readiness:

1. **Authentication & Authorization:**
   - Implement JWT authentication with role-based permissions (`Admin`, `Warehouse Operator`, `Cashier`).
2. **Idempotency Keys:**
   - Add `Idempotency-Key` headers on `POST /api/orders` to safeguard against duplicate charges during network reconnects.
3. **Database Migration Manager:**
   - Integrate `golang-migrate` or `goose` CLI to manage automated rolling migrations and rollbacks in CI/CD pipelines.
4. **Interactive Order Management UI:**
   - Expand the frontend dashboard with a dedicated "Point of Sale (POS) / Order Creation" cart flow and order history listing.
5. **Caching & Async Events:**
   - Add Redis for caching high-read catalog queries and an event bus (e.g. RabbitMQ or Kafka) to broadcast inventory threshold alerts asynchronously.
6. **OpenAPI / Swagger Specification:**
   - Embed auto-generated Swagger UI (`swag`) for interactive API exploration.

---

## Honest Hours Spent

| Phase                                   | Description                                                                                       | Hours      |
| --------------------------------------- | ------------------------------------------------------------------------------------------------- | ---------- |
| **1. System Design & DB Setup**         | PostgreSQL schema design, constraints, migration script, seed data                                | 1.5 hrs    |
| **2. Backend Core Modules**             | Go layered structure, Product CRUD, pagination, search, stock movements                           | 3.5 hrs    |
| **3. Order Logic & Concurrency**        | Transactional checkout, `SELECT ... FOR UPDATE` row locks, order cancellation & stock recovery    | 3.5 hrs    |
| **4. Testing & Concurrency Validation** | Automated parallel goroutine race tests (`order_test.go`), boundary validations                   | 2.5 hrs    |
| **5. Frontend Web Application**         | React 19 + TypeScript + Tailwind v4 setup, API client, Summary Cards, Product Table, Adjust Modal | 7.0 hrs    |
| **6. Dockerization & DevOps**           | Multi-stage Dockerfiles (Go & Nginx), Docker Compose orchestration, healthchecks                  | 1.5 hrs    |
| **7. Documentation & Polishing**        | Root, backend, and frontend READMEs, cURL test documentation, code cleanup                        | 1.5 hrs    |
| **Total**                               |                                                                                                   | **21 hrs** |
