ALTER TABLE jobs RENAME COLUMN job_id TO ray_job_id;
ALTER TABLE jobs DROP COLUMN job_uuid;