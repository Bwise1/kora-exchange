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

-- Transactions table - stores all transaction types (DEPOSIT, SWAP, TRANSFER)
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    transaction_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'COMPLETED',
    wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipient_wallet_id UUID REFERENCES wallets(id) ON DELETE SET NULL,
    from_currency VARCHAR(10) NOT NULL,
    from_amount NUMERIC(20, 8) NOT NULL,
    to_currency VARCHAR(10),
    to_amount NUMERIC(20, 8),
    exchange_rate NUMERIC(20, 8),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT check_transaction_type CHECK (transaction_type IN ('DEPOSIT', 'SWAP', 'TRANSFER', 'WITHDRAW')),
    CONSTRAINT check_status CHECK (status IN ('PENDING', 'COMPLETED', 'FAILED'))
);

-- Indexes for transactions table
CREATE INDEX IF NOT EXISTS idx_transactions_wallet_id ON transactions(wallet_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_recipient_wallet_id ON transactions(recipient_wallet_id);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(transaction_type);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at DESC);

-- Audit logs table - stores all API requests for security and compliance
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    operation VARCHAR(50) NOT NULL,
    client_ip VARCHAR(100) NOT NULL,
    user_agent TEXT,
    request_method VARCHAR(10),
    request_path TEXT,
    request_body TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_audit_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Indexes for audit_logs table
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_operation ON audit_logs(operation);
CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_client_ip ON audit_logs(client_ip);
