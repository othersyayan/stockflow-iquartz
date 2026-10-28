INSERT INTO products (sku, name, price, stock, is_active) VALUES
('SKU-001', 'Kopi Arabika 250g', 85000, 50, true),
('SKU-002', 'Susu UHT 1L', 18000, 100, true),
('SKU-003', 'Gula Pasir 1kg', 16500, 8, true),
('SKU-004', 'Minyak Goreng 2L', 34000, 5, true),
('SKU-005', 'Beras Premium 5kg', 72000, 30, true),
('SKU-006', 'Teh Celup Pack', 12000, 4, true),
('SKU-007', 'Mie Instan Dus', 115000, 15, true),
('SKU-008', 'Tepung Terigu 1kg', 13000, 60, true),
('SKU-009', 'Sabun Cuci Piring', 15000, 20, true),
('SKU-010', 'Deterjen Bubuk 1kg', 24000, 9, true),
('SKU-011', 'Air Mineral 600ml Dus', 48000, 2, true),
('SKU-012', 'Tissue Facial 250s', 10000, 0, false)
ON CONFLICT (sku) DO NOTHING;