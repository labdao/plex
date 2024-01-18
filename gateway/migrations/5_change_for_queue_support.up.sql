BEGIN;

-- Add the new ID column to the jobs table
ALTER TABLE jobs ADD COLUMN id SERIAL UNIQUE;

-- Create new many-to-many relation tables
CREATE TABLE job_input_files (
    job_id INT NOT NULL,
    data_file_c_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (job_id, data_file_c_id),
    FOREIGN KEY (job_id) REFERENCES jobs(id),
    FOREIGN KEY (data_file_c_id) REFERENCES data_files(cid)
);

CREATE TABLE job_output_files (
    job_id INT NOT NULL,
    data_file_c_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (job_id, data_file_c_id),
    FOREIGN KEY (job_id) REFERENCES jobs(id),
    FOREIGN KEY (data_file_c_id) REFERENCES data_files(cid)
);

-- Move over data from previous join tables
INSERT INTO job_input_files (job_id, data_file_c_id)
SELECT j.id, ji.data_file_c_id
FROM job_inputs ji
JOIN jobs j ON ji.job_bacalhau_job_id = j.bacalhau_job_id;

INSERT INTO job_output_files (job_id, data_file_c_id)
SELECT j.id, jo.data_file_c_id
FROM job_outputs jo
JOIN jobs j ON jo.job_bacalhau_job_id = j.bacalhau_job_id;

-- Remove old join tables
DROP TABLE job_inputs;
DROP TABLE job_outputs;

-- Make 'id' instead of 'bacalhau_job-id' the primary key for jobs
ALTER TABLE jobs DROP CONSTRAINT jobs_pkey;
ALTER TABLE jobs ADD PRIMARY KEY (id);

-- Other jobs table modifications
ALTER TABLE jobs ALTER COLUMN state SET DEFAULT 'queued';
ALTER TABLE jobs ADD COLUMN inputs JSON;
ALTER TABLE jobs ADD COLUMN queue VARCHAR(255);
ALTER TABLE jobs ADD COLUMN annotations VARCHAR(255);

-- Tool model changes
ALTER TABLE tools ADD COLUMN container TEXT;
ALTER TABLE tools ADD COLUMN memory INT;
ALTER TABLE tools ADD COLUMN cpu FLOAT;
ALTER TABLE tools ADD COLUMN gpu INT;
ALTER TABLE tools ADD COLUMN network BOOLEAN;

-- Update tools based on ToolJson data
UPDATE tools SET
    container = COALESCE(NULLIF((tool_json ->> 'dockerPull')::TEXT, ''), 'unknown'), 
    memory = (tool_json ->> 'memoryGB')::INT,
    cpu = (tool_json ->> 'cpu')::FLOAT,
    gpu = CASE WHEN (tool_json ->> 'gpuBool')::BOOLEAN THEN 1 ELSE 0 END,
    network = (tool_json ->> 'networkBool')::BOOLEAN
WHERE tool_json IS NOT NULL;

-- Flow model changes
ALTER TABLE flows DROP CONSTRAINT flows_pkey;
ALTER TABLE flows ADD COLUMN id SERIAL PRIMARY KEY;
ALTER TABLE flows ALTER COLUMN cid SET NOT NULL;

-- Add a temporary column to jobs for the new foreign key
ALTER TABLE jobs ADD COLUMN temp_flow_id INT;

-- Update the temp_flow_id with the corresponding flow's new id
UPDATE jobs
SET temp_flow_id = f.id
FROM flows f
WHERE jobs.flow_id = f.cid;

-- Drop the old FlowID column and rename the temp column
ALTER TABLE jobs DROP COLUMN flow_id;
ALTER TABLE jobs RENAME COLUMN temp_flow_id TO flow_id;

-- Add foreign key constraint
ALTER TABLE jobs ADD CONSTRAINT fk_jobs_flows FOREIGN KEY (flow_id) REFERENCES flows(id);

-- Drop the old cid index if it exists
DROP INDEX IF EXISTS idx_jobs_flow_id;

CREATE INDEX IF NOT EXISTS idx_jobs_flow_id ON jobs (flow_id);

COMMIT;
