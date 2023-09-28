-- Add the 'name' column without NOT NULL constraint
ALTER TABLE tools ADD COLUMN name VARCHAR(255);

-- Extract and set the 'name' value from 'ToolJSON' for all rows
UPDATE tools SET name = (SELECT value FROM json_each_text(tool_json) WHERE key = 'name');

-- Set the NOT NULL constraint on 'name'
ALTER TABLE tools ALTER COLUMN name SET NOT NULL;

-- Drop the 'ToolJSON' column
ALTER TABLE tools DROP COLUMN tool_json;