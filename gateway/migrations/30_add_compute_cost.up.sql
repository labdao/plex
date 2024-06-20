-- Add the 'compute_cost' column to the 'tools' table
ALTER TABLE tools ADD COLUMN compute_cost INTEGER NOT NULL DEFAULT 0;