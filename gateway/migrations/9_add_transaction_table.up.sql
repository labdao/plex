CREATE TABLE transactions (
    id VARCHAR(255) PRIMARY KEY,
    amount FLOAT NOT NULL,
    is_debit BOOLEAN NOT NULL,
    user_id VARCHAR(42) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    FOREIGN KEY (user_id) REFERENCES users(wallet_address) ON DELETE CASCADE
);

