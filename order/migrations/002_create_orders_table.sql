-- +goose Up
CREATE TABLE orders (
    order_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid UUID NOT NULL,
    part_uuids UUID[] NOT NULL,
    total_price DECIMAL(12,2) NOT NULL CHECK (total_price >= 0),
    transaction_uuid UUID,
    payment_method payment_method NOT NULL DEFAULT 'UNKNOWN',
    order_status order_status,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);