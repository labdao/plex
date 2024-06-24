-- RENAME TABLE TOOL TO MODEL
ALTER TABLE tools RENAME TO models;

ALTER TABLE models RENAME COLUMN default_tool TO default_model;
ALTER TABLE models RENAME COLUMN tool_json TO model_json;
ALTER TABLE models RENAME COLUMN tool_type TO model_type;

--rename tool referenced id column in jobs table. also the foreign key constraint
ALTER TABLE jobs RENAME COLUMN tool_id TO model_id;
ALTER TABLE jobs DROP CONSTRAINT fk_jobs_tool;
ALTER TABLE jobs ADD CONSTRAINT fk_jobs_model FOREIGN KEY (model_id) REFERENCES models(cid) ON DELETE CASCADE;

--same with flow to experiment
ALTER TABLE flows RENAME TO experiments;

ALTER TABLE experiments RENAME COLUMN flow_uuid TO experiment_uuid;

--in jobs
ALTER TABLE jobs RENAME COLUMN flow_id TO experiment_id;
ALTER TABLE jobs DROP CONSTRAINT IF EXISTS jobs_flowid_fkey;
ALTER TABLE jobs DROP CONSTRAINT IF EXISTS fk_flows_jobs;
ALTER TABLE jobs ADD CONSTRAINT jobs_experimentid_fkey FOREIGN KEY (experiment_id) REFERENCES experiments(id) ON DELETE CASCADE;