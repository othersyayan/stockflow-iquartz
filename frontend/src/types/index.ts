export interface Product {
  id: number;
  sku: string;
  name: string;
  price: number;
  stock: number;
  is_active: boolean;
  created_at: string;
  last_movement_at: string | null;
}

export interface ProductParams {
  search: string;
  page: number;
  limit: number;
}

export interface CreateProductReq {
  sku: string;
  name: string;
  price: number;
}

export interface AdjustStockReq {
  id: number;
  data: { type: 'in' | 'out'; qty: number; note: string };
}

export interface CreateOrderReq {
  customer: string;
  items: { sku: string; qty: number }[];
}

export interface StockMovement {
  id: number;
  product_id: number;
  type: 'in' | 'out';
  qty: number;
  note: string;
  created_at: string;
}

export interface SalesOrderItem {
  id?: number;
  order_id?: number;
  sku: string;
  qty: number;
  unit_price: number;
}

export interface SalesOrder {
  id: number;
  customer: string;
  status: 'CREATED' | 'CANCELLED';
  total: number;
  created_at: string;
  cancelled_at?: string | null;
  items: SalesOrderItem[];
}

export interface ApiErrorResponse {
  error: {
    code: string;
    message?: string;
    sku?: string;
    available?: number;
  };
}
