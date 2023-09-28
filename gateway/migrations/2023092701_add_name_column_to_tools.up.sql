-- Add the 'name' column without NOT NULL constraint
ALTER TABLE tools ADD COLUMN name VARCHAR(255);

-- Extract and set the 'name' value from 'ToolJSON' for all rows
UPDATE tools
SET name = sub.name
FROM (
    SELECT t.cid, j.value as name
    FROM tools t, LATERAL json_each_text(t.tool_json::json) as j(key, value)
    WHERE j.key = 'name'
) sub
WHERE tools.cid = sub.cid;

-- Set the NOT NULL constraint on 'name'
ALTER TABLE tools ALTER COLUMN name SET NOT NULL;

-- Drop the 'ToolJSON' column
ALTER TABLE tools DROP COLUMN tool_json;