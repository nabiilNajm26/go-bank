CREATE TABLE idempotency_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    idempotency_key VARCHAR(255) UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    request_path VARCHAR(500) NOT NULL,
    request_body TEXT,
    response_status INTEGER,
    response_body TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_idempotency_keys_key ON idempotency_keys(idempotency_key);
CREATE INDEX idx_idempotency_keys_user_id ON idempotency_keys(user_id);
CREATE INDEX idx_idempotency_keys_expires_at ON idempotency_keys(expires_at);