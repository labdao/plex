BEGIN;

-- Add the new ID column to the jobs table
ALTER TABLE jobs ADD COLUMN id SERIAL UNIQUE;

-- Update join tables to include the new job_id column
ALTER TABLE job_inputs ADD COLUMN job_id INT;
ALTER TABLE job_outputs ADD COLUMN job_id INT;

-- Populate the new job_id columns in join tables
UPDATE job_inputs SET job_id = j.id FROM jobs j WHERE job_inputs.job_bacalhau_job_id = j.bacalhau_job_id;
UPDATE job_outputs SET job_id = j.id FROM jobs j WHERE job_outputs.job_bacalhau_job_id = j.bacalhau_job_id;

-- Remove old foreign key constraints
ALTER TABLE job_inputs DROP CONSTRAINT job_inputs_job_bacalhau_job_id_fkey;
ALTER TABLE job_outputs DROP CONSTRAINT job_outputs_job_bacalhau_job_id_fkey;

-- Make 'id' the primary key
ALTER TABLE jobs DROP CONSTRAINT jobs_pkey;
ALTER TABLE jobs ADD PRIMARY KEY (id);

-- Add the new foreign key constraints referencing the new job_id
ALTER TABLE job_inputs ADD CONSTRAINT job_inputs_jobid_fkey FOREIGN KEY (job_id) REFERENCES jobs(id);
ALTER TABLE job_outputs ADD CONSTRAINT job_outputs_jobid_fkey FOREIGN KEY (job_id) REFERENCES jobs(id);

-- Renaming join tables
ALTER TABLE job_inputs RENAME TO job_input_files;
ALTER TABLE job_outputs RENAME TO job_output_files;

-- Other jobs table modifications
ALTER TABLE jobs ALTER COLUMN state TYPE VARCHAR(255);
ALTER TABLE jobs ADD COLUMN inputs JSON;
ALTER TABLE jobs ADD COLUMN queue VARCHAR(255) NOT NULL;
ALTER TABLE jobs ADD COLUMN created_at TIMESTAMP;
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

COMMIT;
