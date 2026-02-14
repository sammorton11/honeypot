CREATE TABLE IF NOT EXISTS attempts (
    id SERIAL PRIMARY KEY,
    address TEXT NOT NULL,
    network TEXT NOT NULL,
    message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_attempts_created_at ON attempts(created_at);
CREATE INDEX IF NOT EXISTS idx_attempts_address ON attempts(address);
