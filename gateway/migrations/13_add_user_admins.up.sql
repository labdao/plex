-- Up migration to add the 'admin' column to the 'users' table
ALTER TABLE users ADD COLUMN admin BOOLEAN DEFAULT false;