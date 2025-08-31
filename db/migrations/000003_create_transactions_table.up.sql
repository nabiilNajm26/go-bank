CREATE TYPE transaction_type AS ENUM ('transfer', 'deposit', 'withdrawal', 'payment');
CREATE TYPE transaction_status AS ENUM ('pending', 'completed', 'failed', 'reversed');

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    from_account_id UUID REFERENCES accounts(id) ON DELETE SET NULL,
    to_account_id UUID REFERENCES accounts(id) ON DELETE SET NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    type transaction_type NOT NULL,
    status transaction_status NOT NULL DEFAULT 'pending',
    reference VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT check_transfer_accounts CHECK (
        (type = 'transfer' AND from_account_id IS NOT NULL AND to_account_id IS NOT NULL AND from_account_id != to_account_id) OR
        (type = 'deposit' AND from_account_id IS NULL AND to_account_id IS NOT NULL) OR
        (type = 'withdrawal' AND from_account_id IS NOT NULL AND to_account_id IS NULL) OR
        (type = 'payment' AND from_account_id IS NOT NULL)
    )
);

CREATE INDEX idx_transactions_from_account ON transactions(from_account_id);
CREATE INDEX idx_transactions_to_account ON transactions(to_account_id);
CREATE INDEX idx_transactions_reference ON transactions(reference);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_created_at ON transactions(created_at DESC);