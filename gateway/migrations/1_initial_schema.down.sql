-- Drop many-to-many relation tables
DROP TABLE IF EXISTS job_inputs;
DROP TABLE IF EXISTS job_outputs;

-- Drop users table
DROP TABLE IF EXISTS users;

-- Drop tools table
DROP TABLE IF EXISTS tools;

-- Drop jobs table along with its indexes
DROP INDEX IF EXISTS idx_jobs_tool_id;
DROP INDEX IF EXISTS idx_jobs_flow_id;
DROP TABLE IF EXISTS jobs;

-- Drop flows table
DROP TABLE IF EXISTS flows;

-- Drop data_files table
DROP TABLE IF EXISTS data_files;
