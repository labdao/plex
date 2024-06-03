-- Add the 'tier' column to the 'users' table
ALTER TABLE users ADD COLUMN tier INTEGER NOT NULL DEFAULT 0;

-- Add the 'compute_credits' column to the 'users' table
ALTER TABLE users ADD COLUMN compute_credits INTEGER NOT NULL DEFAULT 0;