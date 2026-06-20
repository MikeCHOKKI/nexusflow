-- =============================================================
-- NexusFlow — Migration 002: Add customer_email to orders
-- Permet de notifier le client lors des changements de statut
-- et d'inclure les articles dans l'email de confirmation.
-- =============================================================

ALTER TABLE orders.orders
    ADD COLUMN customer_email VARCHAR(255),
    ADD COLUMN customer_name  VARCHAR(255) NOT NULL DEFAULT 'Client';
