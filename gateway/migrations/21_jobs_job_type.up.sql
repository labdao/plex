ALTER TABLE jobs RENAME COLUMN bacalhau_job_id TO job_id;
ALTER TABLE jobs ADD COLUMN job_type VARCHAR(255) NOT NULL DEFAULT 'bacalhau';