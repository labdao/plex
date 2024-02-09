CREATE TABLE api_keys (
    id SERIAL PRIMARY KEY,
    key VARCHAR(255) NOT NULL UNIQUE,
    scope VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    user_id VARCHAR(42) NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(wallet_address)
);
