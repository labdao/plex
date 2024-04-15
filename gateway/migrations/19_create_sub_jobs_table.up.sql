-- CREATE TABLE sub_jobs (
--     id SERIAL PRIMARY KEY,
--     job_id INT NOT NULL,
--     binder_index INT NOT NULL,
--     status VARCHAR(255) NOT NULL DEFAULT 'PENDING',
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     FOREIGN KEY (job_id) REFERENCES jobs(id)
-- );

-- linking jobs and job_output_files through sub_jobs
ALTER TABLE job_output_files DROP CONSTRAINT IF EXISTS job_output_files_job_id_fkey;
ALTER TABLE job_output_files RENAME COLUMN job_id TO sub_job_id;
ALTER TABLE job_output_files ADD CONSTRAINT fk_job_output_files_sub_job FOREIGN KEY (sub_job_id) REFERENCES sub_jobs(id) ON DELETE CASCADE;

-- populate sub_jobs table with existing jobs for backward compatibility
INSERT INTO sub_jobs (job_id, binder_index, status, created_at)
SELECT id, 0, 'completed', NOW()
FROM jobs;

-- Get mapping of job_id to the newly created sub_job_id
CREATE TEMP TABLE tmp_job_sub_job_mapping AS
SELECT id AS job_id, currval(pg_get_serial_sequence('sub_jobs', 'id')) AS sub_job_id
FROM jobs;

-- Update job_output_files to link to the new sub_jobs
UPDATE job_output_files
SET sub_job_id = tmp.sub_job_id
FROM tmp_job_sub_job_mapping tmp
WHERE job_output_files.sub_job_id = tmp.job_id;

-- Drop the temporary table after use
DROP TABLE tmp_job_sub_job_mapping;
