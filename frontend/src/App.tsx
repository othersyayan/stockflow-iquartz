import { useEffect, useState, useCallback } from 'react';
import type { Product } from './types';

import { api } from './api/client';

import { SummaryCards } from './components/SummaryCards';

export default function App() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState('');
  const [filter, setFilter] = useState('all');

  const loadProducts = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const res = await api.getProducts({
        search,
        page: 1,
        limit: 100,
      });

      setProducts(res.data || []);
    } catch (err: unknown) {
      setError('Failed to load products. Please check backend connection.');
    } finally {
      setLoading(false);
    }
  }, [search]);

  useEffect(() => {
    loadProducts();
  }, [loadProducts]);

  return (
    <div className="min-h-screen bg-gray-50 text-gray-900 p-6 md:p-10 font-sans">
      <div className="max-w-7xl mx-auto">
        <header className="mb-8 flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold tracking-tight text-gray-800">
              StockFlow
            </h1>
            <p className="text-sm text-gray-500">
              Inventory & Sales Mini ERP Module
            </p>
          </div>
        </header>

        <SummaryCards products={products} />
      </div>
    </div>
  );
}
