import React from 'react';
import type { Product } from '../types';

interface ProductTableProps {
  products: Product[];
  loading: boolean;
  error: string | null;
  onRetry: () => void;
  search: string;
  onSearchChange: (val: string) => void;
  filter: string;
  onFilterChange: (val: string) => void;
  onOpenAdjustModal: (product: Product) => void;
  onDeleteProduct: (id: number) => void;
}

export const ProductTable: React.FC<ProductTableProps> = ({
  products,
  loading,
  error,
  onRetry,
  search,
  onSearchChange,
  filter,
  onFilterChange,
  onOpenAdjustModal,
  onDeleteProduct,
}) => {
  const formatIDR = (val: number) => `Rp ${val.toLocaleString('id-ID')}`;

  const filteredProducts = products.filter((p) => {
    if (filter === 'active') return p.is_active;
    if (filter === 'inactive') return !p.is_active;
    if (filter === 'low') return p.stock <= 10;
    return true;
  });

  return (
    <div className="bg-white rounded-lg shadow border border-gray-100 overflow-hidden">
      <div className="p-4 border-b border-gray-100 flex flex-col md:flex-row justify-between gap-4">
        <input
          type="text"
          placeholder="Search by SKU or Name..."
          value={search}
          onChange={(e) => onSearchChange(e.target.value)}
          className="px-4 py-2 border rounded-md w-full md:w-80 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <div className="flex gap-2">
          {['all', 'active', 'inactive', 'low'].map((f) => (
            <button
              key={f}
              onClick={() => onFilterChange(f)}
              className={`px-3 py-1.5 rounded-md text-sm font-medium capitalize cursor-pointer ${
                filter === f
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              {f === 'low' ? 'Low Stock' : f}
            </button>
          ))}
        </div>
      </div>

      {loading ? (
        <div className="p-8 text-center text-gray-500">Loading products...</div>
      ) : error ? (
        <div className="p-8 text-center text-red-500">
          <p>{error}</p>
          <button
            onClick={onRetry}
            className="mt-2 px-4 py-2 bg-red-600 text-white rounded-md text-sm"
          >
            Retry
          </button>
        </div>
      ) : filteredProducts.length === 0 ? (
        <div className="p-8 text-center text-gray-500">No products found.</div>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-gray-50 text-gray-600 text-xs uppercase tracking-wider">
                <th className="p-4">SKU / Name</th>
                <th className="p-4">Price</th>
                <th className="p-4">Stock</th>
                <th className="p-4">Status</th>
                <th className="p-4">Last Movement</th>
                <th className="p-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 text-sm">
              {filteredProducts.map((p) => (
                <tr key={p.id} className="hover:bg-gray-50">
                  <td className="p-4">
                    <div className="font-semibold text-gray-800">{p.name}</div>
                    <div className="text-xs text-gray-500 font-mono">
                      {p.sku}
                    </div>
                  </td>
                  <td className="p-4 font-medium text-gray-700">
                    {formatIDR(p.price)}
                  </td>
                  <td className="p-4">
                    <span
                      className={`font-bold ${p.stock <= 10 ? 'text-orange-600' : 'text-gray-800'}`}
                    >
                      {p.stock}
                    </span>
                  </td>
                  <td className="p-4">
                    <span
                      className={`px-2 py-1 rounded-full text-xs font-semibold ${
                        p.is_active
                          ? 'bg-green-100 text-green-700'
                          : 'bg-red-100 text-red-700'
                      }`}
                    >
                      {p.is_active ? 'Active' : 'Inactive'}
                    </span>
                  </td>
                  <td className="p-4 text-gray-500 text-xs">
                    {p.last_movement_at
                      ? new Date(p.last_movement_at).toLocaleString()
                      : 'No history'}
                  </td>
                  <td className="p-4 text-right space-x-2">
                    <button
                      onClick={() => onOpenAdjustModal(p)}
                      className="px-3.5 py-1.5 bg-blue-50 text-blue-600 hover:bg-blue-100 rounded-xl text-xs font-medium cursor-pointer"
                    >
                      Adjust Stock
                    </button>
                    <button
                      onClick={() => onDeleteProduct(p.id)}
                      className="px-3.5 py-1.5 bg-red-50 text-red-600 hover:bg-red-100 rounded-xl text-xs font-medium cursor-pointer"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};
