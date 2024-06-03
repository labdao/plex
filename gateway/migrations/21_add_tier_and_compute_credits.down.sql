-- Remove the 'compute_credits' column from the 'users' table
ALTER TABLE users DROP COLUMN compute_credits;

-- Remove the 'tier' column from the 'users' table
ALTER TABLE users DROP COLUMN tier;