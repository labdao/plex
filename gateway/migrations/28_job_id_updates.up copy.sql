ALTER TABLE jobs RENAME COLUMN ray_job_id TO job_id;
ALTER TABLE jobs ADD COLUMN job_uuid UUID;