-- Rename table
ALTER TABLE tools RENAME TO models;

-- Rename columns
ALTER TABLE jobs RENAME COLUMN tool_id TO model_id;
ALTER TABLE models RENAME tool_json TO model_json;
ALTER TABLE models RENAME tool_type TO model_type;

-- Update constraints and keys
ALTER TABLE models RENAME CONSTRAINT tools_name_key TO models_name_key;
ALTER TABLE models RENAME CONSTRAINT tools_pkey TO models_pkey;
ALTER TABLE jobs RENAME CONSTRAINT fk_jobs_tool TO fk_jobs_model;
ALTER TABLE jobs RENAME CONSTRAINT jobs_toolid_fkey TO jobs_modelid_fkey;

-- Update indexes
DROP INDEX idx_jobs_tool_id;
CREATE INDEX idx_jobs_model_id ON jobs USING btree (model_id);
