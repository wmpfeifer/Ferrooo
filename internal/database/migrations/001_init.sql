CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    correlation_id UUID UNIQUE NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    processed_by VARCHAR(10) NOT NULL,
    requested_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payments_correlation_id ON payments(correlation_id);
CREATE INDEX IF NOT EXISTS idx_payments_processed_by ON payments(processed_by);
CREATE INDEX IF NOT EXISTS idx_payments_requested_at ON payments(requested_at);