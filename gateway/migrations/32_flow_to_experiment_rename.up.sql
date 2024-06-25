-- Rename tables
ALTER TABLE flows RENAME TO experiments;

-- Rename columns in related tables
ALTER TABLE jobs RENAME COLUMN flow_id TO experiment_id;

-- Rename primary key and other constraints
ALTER TABLE experiments RENAME CONSTRAINT flows_pkey TO experiments_pkey;

-- Update foreign keys
ALTER TABLE jobs RENAME CONSTRAINT fk_flows_jobs TO fk_experiments_jobs;
ALTER TABLE jobs RENAME CONSTRAINT jobs_flowid_fkey TO jobs_experimentid_fkey;

-- Update indexes
DROP INDEX idx_jobs_flow_id;
CREATE INDEX idx_jobs_experiment_id ON jobs USING btree (experiment_id);
