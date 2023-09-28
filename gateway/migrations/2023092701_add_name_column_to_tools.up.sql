-- Add the 'name' column without NOT NULL constraint
ALTER TABLE tools ADD COLUMN name VARCHAR(255);

-- Extract and set the 'name' value from 'ToolJSON' for all rows
UPDATE tools
SET name = sub.name
FROM (
    SELECT tool_id, value as name
    FROM tools, json_each_text(tools.tool_json)
    WHERE key = 'name'
) sub
WHERE tools.tool_id = sub.tool_id;

-- Set the NOT NULL constraint on 'name'
ALTER TABLE tools ALTER COLUMN name SET NOT NULL;

-- Drop the 'ToolJSON' column
ALTER TABLE tools DROP COLUMN tool_json;