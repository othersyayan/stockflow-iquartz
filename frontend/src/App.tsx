import { useEffect, useState, useCallback } from 'react';
import type { Product, ApiErrorResponse } from './types';

import { api } from './api/client';

import { SummaryCards } from './components/SummaryCards';
import { ProductTable } from './components/ProductTable';
import { AdjustModalStock } from './components/AdjustStockModal';

export default function App() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState('');
  const [filter, setFilter] = useState('all');

  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
  const [adjustType, setAdjustType] = useState<'in' | 'out'>('in');
  const [adjustQty, setAdjustQty] = useState<number>(1);
  const [adjustNote, setAdjustNote] = useState('');
  const [adjustError, setAdjustError] = useState<{
    message: string;
    available?: number;
  } | null>(null);
  const [isSaving, setIsSaving] = useState(false);

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

  const handleAdjustType = (type: 'in' | 'out') => {
    setAdjustType(type);
    setAdjustError(null);
  };

  const handleAdjustQty = (qty: number) => {
    setAdjustQty(qty);
    setAdjustError(null);
  };

  const handleAdjustNote = (note: string) => {
    setAdjustNote(note);
    setAdjustError(null);
  };

  const handleAdjustStockSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedProduct) return;
    setIsSaving(true);
    setAdjustError(null);

    try {
      await api.adjustStock({
        id: selectedProduct.id,
        data: {
          type: adjustType,
          qty: Number(adjustQty),
          note: adjustNote,
        },
      });
      setSelectedProduct(null);
      setAdjustQty(1);
      setAdjustNote('');
    } catch (err: unknown) {
      const apiErr = err as ApiErrorResponse;
      if (apiErr.error && apiErr.error.code === 'INSUFFICIENT_STOCK') {
        setAdjustError({
          message: `Insufficient stock! Available quantity:`,
          available: apiErr.error.available,
        });
      } else {
        setAdjustError({ message: 'Failed to adjust stock.' });
      }
    } finally {
      setIsSaving(false);
      loadProducts();
    }
  };

  const handleDeleteProduct = async (id: number) => {
    if (!confirm('Are you sure you want to delete this product?')) return;
    try {
      await api.deleteProduct(id);
    } catch (err: unknown) {
      const apiErr = err as ApiErrorResponse;
      if (apiErr.error) {
        if (apiErr.error.code === 'PRODUCT_HAS_HISTORY') {
          alert(
            `${apiErr.error.message}! \n Please delete the movement history first.`,
          );
        } else {
          alert(`${apiErr.error.message}`);
        }
      } else {
        alert('Product already deleted!');
      }
    } finally {
      loadProducts();
    }
  };

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

        <ProductTable
          products={products}
          loading={loading}
          error={error}
          onRetry={loadProducts}
          search={search}
          onSearchChange={setSearch}
          filter={filter}
          onFilterChange={setFilter}
          onOpenAdjustModal={(p) => {
            setSelectedProduct(p);
            setAdjustError(null);
          }}
          onDeleteProduct={handleDeleteProduct}
        />

        {selectedProduct && (
          <AdjustModalStock
            isSaving={isSaving}
            adjustNote={adjustNote}
            adjustQty={adjustQty}
            adjustType={adjustType}
            product={selectedProduct}
            adjustError={adjustError}
            handleAdjustQty={handleAdjustQty}
            handleAdjustType={handleAdjustType}
            handleAdjustNote={handleAdjustNote}
            handleSelectedProduct={setSelectedProduct}
            handleAdjustStockSubmit={handleAdjustStockSubmit}
          />
        )}
      </div>
    </div>
  );
}
