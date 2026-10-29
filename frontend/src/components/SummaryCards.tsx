import React from 'react';
import type { Product } from '../types';

interface SummaryCardsProps {
  products: Product[];
}

export const SummaryCards: React.FC<SummaryCardsProps> = ({ products }) => {
  const totalProducts = products.length;
  const activeProducts = products.filter((p) => p.is_active).length;
  const totalStockValue = products.reduce(
    (acc, p) => acc + p.price * p.stock,
    0,
  );
  const lowStockCount = products.filter(
    (p) => p.stock <= 10 && p.is_active,
  ).length;

  const formatIDR = (val: number) => `Rp ${val.toLocaleString('id-ID')}`;

  return (
    <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
      <div className="bg-white p-4 rounded-lg shadow border border-gray-100">
        <p className="text-sm text-gray-500 font-medium">Total Products</p>
        <p className="text-2xl font-bold text-gray-800">{totalProducts}</p>
      </div>
      <div className="bg-white p-4 rounded-lg shadow border border-gray-100">
        <p className="text-sm text-gray-500 font-medium">Active Products</p>
        <p className="text-2xl font-bold text-green-600">{activeProducts}</p>
      </div>
      <div className="bg-white p-4 rounded-lg shadow border border-gray-100">
        <p className="text-sm text-gray-500 font-medium">Total Stock Value</p>
        <p className="text-2xl font-bold text-blue-600">
          {formatIDR(totalStockValue)}
        </p>
      </div>
      <div className="bg-white p-4 rounded-lg shadow border border-gray-100">
        <p className="text-sm text-gray-500 font-medium">Low Stock (≤ 10)</p>
        <p className="text-2xl font-bold text-orange-600">{lowStockCount}</p>
      </div>
    </div>
  );
};
