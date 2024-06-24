-- RENAME TABLE TOOL TO MODEL
ALTER TABLE tool RENAME TO model;

ALTER TABLE model RENAME COLUMN default_tool TO default_model;
ALTER TABLE model RENAME COLUMN tool_json TO model_json;
ALTER TABLE model RENAME COLUMN tool_type TO model_type;

--rename tool referenced id column in jobs table. also the foreign key constraint
ALTER TABLE jobs RENAME COLUMN tool_id TO model_id;
ALTER TABLE jobs DROP CONSTRAINT jobs_tool_id_fkey;
ALTER TABLE jobs ADD CONSTRAINT jobs_model_id_fkey FOREIGN KEY (model_id) REFERENCES model(id) ON DELETE CASCADE;

--same with flow to experiment
ALTER TABLE flow RENAME TO experiment;

ALTER TABLE experiment RENAME COLUMN flow_uuid TO experiment_uuid;

--in jobs
ALTER TABLE jobs RENAME COLUMN flow_id TO experiment_id;
ALTER TABLE jobs DROP CONSTRAINT jobs_flow_id_fkey;
ALTER TABLE jobs ADD CONSTRAINT jobs_experiment_id_fkey FOREIGN KEY (experiment_id) REFERENCES experiment(id) ON DELETE CASCADE;