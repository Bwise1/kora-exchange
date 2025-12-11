-- Initial database schema for Operation Borderless
-- Creates tables for wallets, transactions, fx_rates, and audit_logs
-- Uses NUMERIC for monetary amounts and rates for fintech-grade precision
-- Includes cascading behavior for foreign keys to ensure data integrity


-- Creating users table to store users information
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,    
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

-- Wallets table - stores multi-currency balances in JSONB
CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    wallet_address VARCHAR(100) UNIQUE NOT NULL,
    balances JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for wallets table
CREATE INDEX IF NOT EXISTS idx_wallets_balances ON wallets USING GIN (balances);
CREATE INDEX IF NOT EXISTS idx_wallets_user_id ON wallets(user_id);
CREATE INDEX IF NOT EXISTS idx_wallets_address ON wallets(wallet_address);