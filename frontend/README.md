# StockFlow Frontend Web App

Modern, responsive web dashboard for **StockFlow** (Inventory & Sales Mini ERP Module), built with **React 19**, **TypeScript**, **Vite 8**, and **Tailwind CSS v4**.

---

## Architecture & Tech Stack

- **UI Framework:** React 19 (`react`, `react-dom`)
- **Language:** TypeScript 5.5+
- **Bundler & Tooling:** Vite 8 with Hot Module Replacement (HMR)
- **Styling:** Tailwind CSS v4 (using official `@tailwindcss/vite` plugin)
- **Reverse Proxy / Production Server:** Nginx Alpine (in Docker container)

### Component Hierarchy

```
frontend/
├── src/
│   ├── api/
│   │   └── client.ts            # Typed HTTP client for all backend endpoints
│   ├── components/
│   │   ├── AdjustStockModal.tsx # Stock In / Out adjustment modal dialog
│   │   ├── ProductTable.tsx     # Filterable, searchable product inventory table
│   │   └── SummaryCards.tsx     # KPI metrics (Total products, active, valuation, low stock)
│   ├── types/
│   │   └── index.ts             # Domain models, request payloads, API response interfaces
│   ├── App.tsx                  # Root page orchestrator with global state & data fetching
│   ├── index.css                # Tailwind CSS v4 directives
│   └── main.tsx                 # React DOM mount point
├── Dockerfile                   # Multi-stage build (Node 20 builder -> Nginx server)
├── package.json                 # Dependency definitions and scripts
└── vite.config.ts               # Vite configuration with `/api` reverse proxy
```

---

## How to Run (Step by Step)

### Option 1: Run with Docker Compose

From the project root:

```bash
docker compose up --build frontend
```

The dashboard will be accessible at: `http://localhost:5173`.

---

### Option 2: Run Locally (Development Server)

#### 1. Prerequisites

- **Node.js 20+** installed (`node -v`)
- **npm 10+** installed (`npm -v`)
- Backend running on `http://localhost:3000` (or configure via proxy)

#### 2. Install Dependencies

```bash
npm install
```

#### 3. Start Development Server

```bash
npm run dev
```

The application will launch at `http://localhost:5173`. Requests to `/api/*` are automatically forwarded to `http://localhost:3000` via Vite's proxy configuration in `vite.config.ts`.

#### 4. Build for Production

```bash
npm run build
```

Generates production artifacts inside the `dist/` directory.

#### 5. Run Linter

```bash
npm run lint
```

---

## User Interface & Key Features

1. **Inventory Summary Cards:**
   - **Total Products:** Total product count in the catalog.
   - **Active Products:** Count of products currently marked as active.
   - **Total Stock Value:** Aggregated valuation (`SUM(price * stock)`), formatted in Indonesian Rupiah (`Rp`).
   - **Low Stock Indicator:** Products with stock quantity $\le 10$ flagged for restocking attention.

2. **Interactive Product Table:**
   - Real-time search by Product Name or SKU.
   - Filter chips: `All`, `Active`, `Inactive`, and `Low Stock`.
   - Visual stock alerts (bold orange badge for low stock).
   - Relative/Formatted timestamps for last recorded stock movement.

3. **Stock Adjustment Dialog:**
   - Seamlessly perform **Stock In** (restock) or **Stock Out** (usage / write-off).
   - Real-time validation preventing invalid quantities.
   - Clear error callout highlighting available stock when insufficient quantity is requested.

4. **Guarded Product Deletion:**
   - Safety checks preventing accidental deletion of products that have movement history records.

---

## API Integration & cURL Examples

The frontend interacts with the backend through `/api/*`. Below are the primary API calls the UI communicates with:

### 1. Fetching Inventory List

```bash
curl -X GET "http://localhost:3000/api/products?search=&page=1&limit=100"
```

**UI Usage:** Populates the `SummaryCards` and `ProductTable`.

### 2. Performing Stock In / Out Adjustment

```bash
curl -X POST "http://localhost:3000/api/products/1/stock" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "in",
    "qty": 15,
    "note": "Restocked via frontend modal"
  }'
```

**UI Usage:** Triggered by `AdjustModalStock` on submit.

### 3. Deleting Product

```bash
curl -X DELETE "http://localhost:3000/api/products/1"
```

**UI Usage:** Triggered by the "Delete" action button in the table row.

---

## Design Decisions

1. **Native Fetch with Centralized Error Handling:**
   - Used a lightweight typed API client (`src/api/client.ts`) wrapping native browser `fetch`. Avoided heavyweight HTTP libraries (like Axios) to minimize bundle footprint while preserving full TypeScript safety.
2. **Tailwind CSS v4 with `@tailwindcss/vite`:**
   - Implemented Tailwind CSS 4 for blazing-fast JIT compile speeds and zero-runtime CSS footprint.
3. **Optimistic vs. Server-Authoritative Refresh:**
   - Chose authoritative re-fetching after mutations (`loadProducts()`) rather than purely optimistic client updates to ensure the calculated summary metrics (like Total Stock Value and Low Stock count) are always 100% synchronized with the database state.
4. **Proxy Setup in Vite & Nginx:**
   - In dev mode, `vite.config.ts` proxies `/api` to `localhost:3000`. In Docker production mode, an Nginx reverse proxy routes `/api/` directly to the `backend:3000` service container, avoiding CORS headaches across environments.

---

## What I Would Improve

1. **State Management & Caching:** Migrate from raw `useEffect` to **TanStack Query (React Query)** for automatic background revalidation, cache invalidation, and request deduplication.
2. **Sales Order Management UI:** Build a dedicated Order Creation modal and Order History tab to exercise the backend's `/api/orders` endpoints directly from the browser.
3. **Form Validation with Zod:** Integrate `react-hook-form` and `zod` for declarative schema validation on modals.
4. **Toast Notification System:** Replace standard browser `alert()` and `confirm()` with custom, accessible toast notifications (e.g. `sonner`).
5. **Pagination UI:** Implement server-side pagination controls (Next/Previous/Page numbers) instead of fetching up to 100 items at once.

---

## Honest Hours Spent

- **Project Scaffolding & Tailwind v4 Setup:** 1.5 hour (Vite, TS configuration, asset structure)
- **API Client & TypeScript Types Definition:** 1.5 hour (Contract alignment with backend)
- **Summary Cards & Inventory Table Components:** 2.0 hours (Filtering, search, currency formatting)
- **Stock Adjustment Modal & Error Handling:** 2.0 hours (In/out logic, modal states, insufficient stock feedback)
- **Total Frontend Effort:** **~7 hours**
