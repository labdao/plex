BEGIN;

-- Rename tables
ALTER TABLE experiments RENAME TO flows;

-- Rename columns in related tables
ALTER TABLE jobs RENAME COLUMN experiment_id TO flow_id;

-- Rename primary key and other constraints
ALTER TABLE flows RENAME CONSTRAINT experiments_pkey TO flows_pkey;

-- Update foreign keys
ALTER TABLE jobs RENAME CONSTRAINT fk_experiments_jobs TO fk_flows_jobs;
ALTER TABLE jobs RENAME CONSTRAINT jobs_experimentid_fkey TO jobs_flowid_fkey;

-- Update indexes
DROP INDEX idx_jobs_experiment_id;
CREATE INDEX idx_jobs_flow_id ON jobs USING btree (flow_id);

COMMIT;