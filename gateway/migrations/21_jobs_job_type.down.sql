ALTER TABLE jobs DROP COLUMN job_type;
ALTER TABLE jobs RENAME COLUMN job_id TO bacalhau_job_id;