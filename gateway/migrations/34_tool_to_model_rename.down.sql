BEGIN;

ALTER TABLE models RENAME TO tools;

ALTER TABLE jobs RENAME COLUMN model_id TO tool_id;
ALTER TABLE tools RENAME model_json TO tool_json;
ALTER TABLE tools RENAME COLUMN default_model TO default_tool;

DROP INDEX idx_jobs_model_id;
CREATE INDEX idx_jobs_tool_id ON jobs USING btree (tool_id);

ALTER TABLE tools RENAME CONSTRAINT models_name_key TO tools_name_key;
ALTER TABLE tools RENAME CONSTRAINT models_pkey TO tools_pkey;

COMMIT;