-- Remove IsMember column from users table
ALTER TABLE users
DROP COLUMN IF EXISTS is_member;