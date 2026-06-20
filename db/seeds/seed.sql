-- =============================================================
-- NexusFlow — Seed Data
-- =============================================================

-- ─── Users ───────────────────────────────────────────────────
-- Password: "password123" (bcrypt hash)
INSERT INTO users.users (id, email, password, name, role) VALUES
    ('a0000000-0000-0000-0000-000000000001', 'admin@nexusflow.local',
     '$2b$12$sddEUBAmmovyl24OnWi9.OvFYln2LjBRgWNwBFH.7JuJ8hJb2YO/C',
     'Admin Nexus', 'admin'),
    ('a0000000-0000-0000-0000-000000000002', 'alice@example.com',
     '$2b$12$sddEUBAmmovyl24OnWi9.OvFYln2LjBRgWNwBFH.7JuJ8hJb2YO/C',
     'Alice Koné', 'customer'),
    ('a0000000-0000-0000-0000-000000000003', 'bob@example.com',
     '$2b$12$sddEUBAmmovyl24OnWi9.OvFYln2LjBRgWNwBFH.7JuJ8hJb2YO/C',
     'Bob Traoré', 'customer')
ON CONFLICT (id) DO NOTHING;

-- ─── Products ────────────────────────────────────────────────
INSERT INTO catalog.products (id, name, description, price, currency, category, stock) VALUES
    ('b0000000-0000-0000-0000-000000000001', 'Smartphone Pro X',
     'Smartphone 6.7" 256Go, caméra 108MP, batterie 5000mAh', 450000, 'XOF', 'Électronique', 50),
    ('b0000000-0000-0000-0000-000000000002', 'Laptop UltraBook',
     'Intel i7, 16Go RAM, 512Go SSD, écran 15.6" FHD', 850000, 'XOF', 'Électronique', 30),
    ('b0000000-0000-0000-0000-000000000003', 'Casque Audio BT',
     'Casque sans fil, réduction bruit active, 30h autonomie', 75000, 'XOF', 'Audio', 100),
    ('b0000000-0000-0000-0000-000000000004', 'Montre Connectée',
     'Fitness tracker, GPS, cardio, 14 jours autonomie', 95000, 'XOF', 'Accessoires', 75),
    ('b0000000-0000-0000-0000-000000000005', 'Enceinte Portable',
     'Bluetooth 5.3, 360°, IPX7, 20h batterie', 55000, 'XOF', 'Audio', 60)
ON CONFLICT (id) DO NOTHING;

-- ─── Orders ──────────────────────────────────────────────────
INSERT INTO orders.orders (id, user_id, total_amount, currency, status, shipping_address)
VALUES
    ('c0000000-0000-0000-0000-000000000001',
     'a0000000-0000-0000-0000-000000000002',
     525000, 'XOF', 'delivered', 'Cotonou, Bénin — Quartier Haie Vive'),
    ('c0000000-0000-0000-0000-000000000002',
     'a0000000-0000-0000-0000-000000000003',
     850000, 'XOF', 'processing', 'Dakar, Sénégal — Sacré Cœur 3')
ON CONFLICT (id) DO NOTHING;

INSERT INTO orders.order_items (order_id, product_id, product_name, quantity, unit_price, total_price)
VALUES
    ('c0000000-0000-0000-0000-000000000001',
     'b0000000-0000-0000-0000-000000000001', 'Smartphone Pro X', 1, 450000, 450000),
    ('c0000000-0000-0000-0000-000000000001',
     'b0000000-0000-0000-0000-000000000003', 'Casque Audio BT', 1, 75000, 75000),
    ('c0000000-0000-0000-0000-000000000002',
     'b0000000-0000-0000-0000-000000000002', 'Laptop UltraBook', 1, 850000, 850000)
ON CONFLICT DO NOTHING;

-- ─── Payments ────────────────────────────────────────────────
INSERT INTO payments.payments (id, order_id, user_id, amount, currency, method, status, transaction_ref)
VALUES
    ('d0000000-0000-0000-0000-000000000001',
     'c0000000-0000-0000-0000-000000000001',
     'a0000000-0000-0000-0000-000000000002',
     525000, 'XOF', 'wave', 'completed', 'WAVE-TX-2024-001'),
    ('d0000000-0000-0000-0000-000000000002',
     'c0000000-0000-0000-0000-000000000002',
     'a0000000-0000-0000-0000-000000000003',
     850000, 'XOF', 'orange_money', 'processing', 'OM-TX-2024-002')
ON CONFLICT (id) DO NOTHING;
