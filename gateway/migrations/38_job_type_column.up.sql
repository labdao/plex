ALTER TABLE jobs ADD COLUMN IF NOT EXISTS job_type VARCHAR;
ALTER TABLE experiments ADD COLUMN IF NOT EXISTS experiment_uuid VARCHAR;