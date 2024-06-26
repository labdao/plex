BEGIN;

ALTER TABLE flows RENAME TO experiments;

ALTER TABLE experiments RENAME COLUMN flow_uuid TO experiment_uuid;

ALTER TABLE jobs RENAME COLUMN flow_id TO experiment_id;

DROP INDEX idx_jobs_flow_id;
CREATE INDEX idx_jobs_experiment_id ON jobs USING btree (experiment_id);

ALTER SEQUENCE flows_id_seq RENAME TO experiments_id_seq;

ALTER TABLE experiments RENAME CONSTRAINT flows_pkey TO experiments_pkey;
ALTER TABLE jobs RENAME CONSTRAINT jobs_flowid_fkey TO jobs_experimentid_fkey;

COMMIT;