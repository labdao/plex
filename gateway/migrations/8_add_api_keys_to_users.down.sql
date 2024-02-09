-- Remove the DID column from the users table
ALTER TABLE users
DROP COLUMN did;

-- Remove the foreign key constraint from the api_keys table
ALTER TABLE api_keys
DROP CONSTRAINT fk_user;