BEGIN;

-- Rename table
ALTER TABLE models RENAME TO tools;

-- Rename columns
ALTER TABLE jobs RENAME COLUMN model_id TO tool_id;
ALTER TABLE tools RENAME model_json TO tool_json;

-- Update constraints and keys
ALTER TABLE tools RENAME CONSTRAINT models_name_key TO tools_name_key;
ALTER TABLE tools RENAME CONSTRAINT models_pkey TO tools_pkey;
ALTER TABLE jobs RENAME CONSTRAINT fk_jobs_model TO fk_jobs_tool;
ALTER TABLE jobs RENAME CONSTRAINT jobs_modelid_fkey TO jobs_toolid_fkey;

-- Update indexes
DROP INDEX idx_jobs_model_id;
CREATE INDEX idx_jobs_tool_id ON jobs USING btree (tool_id);

COMMIT;