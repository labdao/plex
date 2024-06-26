BEGIN;

ALTER TABLE tools RENAME TO models;

ALTER TABLE jobs RENAME COLUMN tool_id TO model_id;
ALTER TABLE models RENAME tool_json TO model_json;
ALTER TABLE models RENAME COLUMN default_tool TO default_model;

DROP INDEX idx_jobs_tool_id;
CREATE INDEX idx_jobs_model_id ON jobs USING btree (model_id);

ALTER TABLE models RENAME CONSTRAINT tools_name_key TO models_name_key;
ALTER TABLE models RENAME CONSTRAINT tools_pkey TO models_pkey;

COMMIT;