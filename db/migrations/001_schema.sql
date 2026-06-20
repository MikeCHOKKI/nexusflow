-- =============================================================
-- NexusFlow — Schema Initial
-- Schémas isolés par service (same DB, separate concerns)
-- =============================================================

-- ─── Users Schema ────────────────────────────────────────────
CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE users.users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       VARCHAR(255) UNIQUE NOT NULL,
    password    VARCHAR(255) NOT NULL,   -- bcrypt hash
    name        VARCHAR(255) NOT NULL,
    role        VARCHAR(50) NOT NULL DEFAULT 'customer',
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users.users(email);
CREATE INDEX idx_users_role ON users.users(role);

-- ─── Catalog Schema ──────────────────────────────────────────
CREATE SCHEMA IF NOT EXISTS catalog;

CREATE TABLE catalog.products (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    price       DECIMAL(12,2) NOT NULL,
    currency    VARCHAR(3) NOT NULL DEFAULT 'XOF',
    category    VARCHAR(100) NOT NULL,
    stock       INTEGER NOT NULL DEFAULT 0,
    image_url   VARCHAR(500),
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_products_category ON catalog.products(category);
CREATE INDEX idx_products_active ON catalog.products(is_active);
CREATE INDEX idx_products_name ON catalog.products(name);

-- ─── Orders Schema ───────────────────────────────────────────
CREATE SCHEMA IF NOT EXISTS orders;

CREATE TYPE orders.order_status AS ENUM (
    'pending', 'confirmed', 'processing', 'shipped',
    'delivered', 'cancelled', 'refunded'
);

CREATE TABLE orders.orders (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID NOT NULL,
    total_amount      DECIMAL(14,2) NOT NULL,
    currency          VARCHAR(3) NOT NULL DEFAULT 'XOF',
    status            orders.order_status NOT NULL DEFAULT 'pending',
    shipping_address  TEXT,
    payment_id        UUID,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_user ON orders.orders(user_id);
CREATE INDEX idx_orders_status ON orders.orders(status);
CREATE INDEX idx_orders_payment ON orders.orders(payment_id);

CREATE TABLE orders.order_items (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id      UUID NOT NULL REFERENCES orders.orders(id) ON DELETE CASCADE,
    product_id    UUID NOT NULL,
    product_name  VARCHAR(255) NOT NULL,
    quantity      INTEGER NOT NULL,
    unit_price    DECIMAL(12,2) NOT NULL,
    total_price   DECIMAL(12,2) NOT NULL
);

CREATE INDEX idx_order_items_order ON orders.order_items(order_id);

-- ─── Payments Schema ─────────────────────────────────────────
CREATE SCHEMA IF NOT EXISTS payments;

CREATE TYPE payments.payment_method AS ENUM (
    'orange_money', 'mtn_momo', 'wave', 'visa', 'mastercard'
);

CREATE TYPE payments.payment_status AS ENUM (
    'initiated', 'processing', 'completed', 'failed', 'refunded'
);

CREATE TABLE payments.payments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id        UUID NOT NULL,
    user_id         UUID NOT NULL,
    amount          DECIMAL(14,2) NOT NULL,
    currency        VARCHAR(3) NOT NULL DEFAULT 'XOF',
    method          payments.payment_method NOT NULL,
    status          payments.payment_status NOT NULL DEFAULT 'initiated',
    transaction_ref VARCHAR(255),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_order ON payments.payments(order_id);
CREATE INDEX idx_payments_user ON payments.payments(user_id);
CREATE INDEX idx_payments_status ON payments.payments(status);
