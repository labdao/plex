-- Add IsMember column to users table
ALTER TABLE users
ADD COLUMN is_member BOOLEAN NOT NULL DEFAULT FALSE;