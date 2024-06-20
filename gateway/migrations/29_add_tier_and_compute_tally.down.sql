-- Remove the 'compute_tally' column from the 'users' table
ALTER TABLE users DROP COLUMN compute_tally;

-- Remove the 'tier' column from the 'users' table
ALTER TABLE users DROP COLUMN tier;