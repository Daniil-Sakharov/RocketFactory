-- +goose Up
CREATE TYPE payment_method AS ENUM (
    'UNKNOWN',
    'CARD',
    'SBP',
    'CREDIT_CARD',
    'INVESTOR_MONEY'
);

CREATE TYPE order_status AS ENUM (
    'PENDING_PAYMENT',
    'PAID',
    'CANCELLED'
);
