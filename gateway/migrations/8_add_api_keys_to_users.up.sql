-- Add the DID column to the users table
ALTER TABLE users
ADD COLUMN did VARCHAR(255) UNIQUE;

-- Ensure the foreign key constraint from the api_keys table to the users table is correct
ALTER TABLE api_keys
DROP CONSTRAINT IF EXISTS fk_user,
ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(wallet_address);