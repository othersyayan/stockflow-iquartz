import type {
  Product,
  SalesOrder,
  ProductParams,
  StockMovement,
  AdjustStockReq,
  CreateProductReq,
  ApiErrorResponse,
  CreateOrderReq,
} from '../types';

async function handleResponse<T>(res: Response): Promise<T> {
  if (!res.ok) {
    const errorBody: ApiErrorResponse = await res.json().catch(() => ({
      error: { code: 'UNKNOWN_ERROR', message: res.statusText },
    }));
    throw errorBody;
  }
  if (res.status === 204) {
    return {} as T;
  }
  return res.json();
}

export const api = {
  getProducts: async ({ search, page, limit }: ProductParams) => {
    const res = await fetch(
      `/api/products?search=${encodeURIComponent(search)}&page=${page}&limit=${limit}`,
    );

    return handleResponse<{
      data: Product[];
      page: number;
      limit: number;
      total: number;
    }>(res);
  },

  createProduct: async (req: CreateProductReq) => {
    const res = await fetch('/api/products', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(req),
    });

    return handleResponse<Product>(res);
  },

  deleteProduct: async (id: number) => {
    const res = await fetch(`/api/products/${id}`, { method: 'DELETE' });
    return handleResponse<void>(res);
  },

  adjustStock: async ({ id, data }: AdjustStockReq) => {
    const res = await fetch(`/api/products/${id}/stock`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });

    return handleResponse<Product>(res);
  },

  getMovements: async (id: number) => {
    const res = await fetch(`/api/products/${id}/movements`);
    return handleResponse<{ data: StockMovement[] }>(res);
  },

  createOrder: async ({ customer, items }: CreateOrderReq) => {
    const res = await fetch('/api/orders', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ customer, items }),
    });

    return handleResponse<SalesOrder>(res);
  },

  getOrder: async (id: number) => {
    const res = await fetch(`/api/orders/${id}`);
    return handleResponse<SalesOrder>(res);
  },

  cancelOrder: async (id: number) => {
    const res = await fetch(`/api/orders/${id}/cancel`, { method: 'POST' });
    return handleResponse<void>(res);
  },
};
