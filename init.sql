-- ==========================================
-- 1. MIGRATION (Skema Database)
-- ==========================================
CREATE TABLE products (
    id         SERIAL PRIMARY KEY,
    sku        TEXT NOT NULL UNIQUE,
    name       TEXT NOT NULL,
    price      INTEGER NOT NULL CHECK (price >= 0),
    stock      INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE stock_movements (
    id         SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id),
    type       TEXT NOT NULL CHECK (type IN ('in', 'out')),
    qty        INTEGER NOT NULL CHECK (qty > 0),
    note       TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE sales_orders (
    id           SERIAL PRIMARY KEY,
    customer     TEXT NOT NULL,
    status       TEXT NOT NULL DEFAULT 'CREATED' CHECK (status IN ('CREATED','CANCELLED')),
    total        INTEGER NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    cancelled_at TIMESTAMPTZ
);

CREATE TABLE sales_order_items (
    id         SERIAL PRIMARY KEY,
    order_id   INTEGER NOT NULL REFERENCES sales_orders(id),
    sku        TEXT NOT NULL REFERENCES products(sku),
    qty        INTEGER NOT NULL CHECK (qty > 0),
    unit_price INTEGER NOT NULL
);

-- ==========================================
-- 2. SEED DATA (Minimal 12 Produk & Contoh Order)
-- ==========================================
INSERT INTO products (sku, name, price, stock, is_active) VALUES
('SKU-001', 'Laptop ThinkPad T14', 18500000, 15, true),
('SKU-002', 'Mechanical Keyboard Keychron K2', 1200000, 45, true),
('SKU-003', 'Logitech MX Master 3S', 1650000, 30, true),
('SKU-004', 'Monitor Dell UltraSharp 27', 6500000, 12, true),
('SKU-005', 'Kabel HDMI 2.1 Baseus 2M', 150000, 100, true),
('SKU-006', 'Standing Desk FlexiSpot', 4200000, 5, true),
('SKU-007', 'Kursi Ergonomis Herman Miller', 15000000, 2, true),
('SKU-008', 'Webcam Logitech C920', 1100000, 8, true),
('SKU-009', 'Microphone Shure MV7', 4500000, 0, true),   -- Contoh stok habis
('SKU-010', 'MacBook Pro M3 14-inch', 32000000, 10, true),
('SKU-011', 'SSD Samsung 980 PRO 1TB', 1800000, 25, true),
('SKU-012', 'Mousepad Rexus Kvlar', 85000, 150, false); -- Contoh produk tidak aktif

-- Seed Pergerakan Stok Awal (Opsional untuk melengkapi data)
INSERT INTO stock_movements (product_id, type, qty, note) VALUES
(1, 'in', 15, 'Initial Stock'),
(2, 'in', 45, 'Initial Stock'),
(3, 'in', 30, 'Initial Stock');

-- Seed Contoh Pesanan (Sales Order)
INSERT INTO sales_orders (customer, total, status) VALUES 
('PT Teknologi Maju', 19700000, 'CREATED'),
('Budi Santoso', 1100000, 'CANCELLED');

INSERT INTO sales_order_items (order_id, sku, qty, unit_price) VALUES
(1, 'SKU-001', 1, 18500000),
(1, 'SKU-002', 1, 1200000),
(2, 'SKU-008', 1, 1100000);