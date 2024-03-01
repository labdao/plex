-- Down migration to remove the 'admin' column from the 'users' table
ALTER TABLE users DROP COLUMN admin;