BEGIN;

ALTER TABLE experiments RENAME TO flows;

ALTER TABLE flows RENAME COLUMN experiment_uuid TO flow_uuid;

ALTER TABLE jobs RENAME COLUMN experiment_id TO flow_id;

DROP INDEX idx_jobs_experiment_id;
CREATE INDEX idx_jobs_flow_id ON jobs USING btree (flow_id);

ALTER SEQUENCE experiments_id_seq RENAME TO flows_id_seq;

ALTER TABLE flows RENAME CONSTRAINT experiments_pkey TO flows_pkey;
ALTER TABLE jobs RENAME CONSTRAINT jobs_experimentid_fkey TO jobs_flowid_fkey;

COMMIT;