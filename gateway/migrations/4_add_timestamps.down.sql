-- Remove StartTime and EndTime columns from flows table
ALTER TABLE flows
DROP COLUMN IF EXISTS start_time,
DROP COLUMN IF EXISTS end_time;

-- Remove Timestamp column from tools table
ALTER TABLE tools
DROP COLUMN IF EXISTS timestamp;

-- Remove CreatedAt column from users table
ALTER TABLE users
DROP COLUMN IF EXISTS created_at;
