import React from 'react';
import type { Product } from '../types';

interface AdjustStockModalProps {
  product: Product | null;
  adjustType: 'in' | 'out';
  adjustQty: number;
  adjustNote: string;
  isSaving: boolean;
  adjustError: {
    message: string;
    available?: number;
  } | null;
  handleAdjustQty: (qty: number) => void;
  handleAdjustNote: (note: string) => void;
  handleAdjustType: (adjustType: 'in' | 'out') => void;
  handleSelectedProduct: (product: Product | null) => void;
  handleAdjustStockSubmit: (e: React.FormEvent<Element>) => Promise<void>;
}

export const AdjustModalStock: React.FC<AdjustStockModalProps> = ({
  product,
  isSaving,
  adjustQty,
  adjustNote,
  adjustType,
  adjustError,
  handleAdjustQty,
  handleAdjustNote,
  handleAdjustType,
  handleSelectedProduct,
  handleAdjustStockSubmit,
}) => {
  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-lg p-6 max-w-md w-full shadow-lg">
        <h3 className="text-lg font-bold mb-2">Adjust Stock: {product.name}</h3>
        <p className="text-sm text-gray-500 mb-4 font-mono">
          SKU: {product.sku} | Current Stock: {product.stock}
        </p>

        {adjustError && (
          <div className="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 rounded text-sm">
            <p>
              {adjustError.message}{' '}
              {adjustError.available !== undefined && (
                <strong className="text-red-800">
                  {adjustError.available}
                </strong>
              )}
            </p>
          </div>
        )}

        <form onSubmit={handleAdjustStockSubmit} className="space-y-4">
          <div>
            <label className="block text-xs font-semibold uppercase text-gray-600 mb-1">
              Direction
            </label>
            <select
              value={adjustType}
              onChange={(e) => handleAdjustType(e.target.value as 'in' | 'out')}
              className="w-full border rounded p-2 text-sm"
            >
              <option value="in">Stock In (Tambah)</option>
              <option value="out">Stock Out (Kurang)</option>
            </select>
          </div>

          <div>
            <label className="block text-xs font-semibold uppercase text-gray-600 mb-1">
              Quantity
            </label>
            <input
              type="number"
              min="1"
              value={adjustQty}
              onChange={(e) => handleAdjustQty(Number(e.target.value))}
              className="w-full border rounded p-2 text-sm"
              required
            />
          </div>

          <div>
            <label className="block text-xs font-semibold uppercase text-gray-600 mb-1">
              Note
            </label>
            <input
              type="text"
              value={adjustNote}
              onChange={(e) => handleAdjustNote(e.target.value)}
              placeholder="Reason for adjustment..."
              className="w-full border rounded p-2 text-sm"
            />
          </div>

          <div className="flex justify-end gap-2 pt-2">
            <button
              type="button"
              onClick={() => handleSelectedProduct(null)}
              className="px-4 py-2 border rounded text-sm text-gray-600 hover:bg-gray-100 cursor-pointer"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isSaving}
              className="px-4 py-2 bg-blue-600 text-white rounded text-sm hover:bg-blue-700 disabled:opacity-50 cursor-pointer"
            >
              {isSaving ? 'Saving...' : 'Save Adjustment'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
