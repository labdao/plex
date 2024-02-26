-- -- Enable UUID extension, if not already present
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- -- Add the UUID column to the 'data_files' table
-- ALTER TABLE data_files ADD COLUMN id UUID NOT NULL DEFAULT uuid_generate_v4();

-- -- Create a unique index on the new UUID column to ensure uniqueness
-- CREATE UNIQUE INDEX IF NOT EXISTS idx_data_files_id ON data_files(id);

-- -- Add the 'data_file_id' UUID column to 'job_input_files'
-- ALTER TABLE job_input_files ADD COLUMN data_file_id UUID;

-- -- Populate 'data_file_id' in 'job_input_files' from 'data_files'
-- UPDATE job_input_files
-- SET data_file_id = df.id
-- FROM data_files AS df
-- WHERE job_input_files.data_file_c_id = df.cid;

-- -- Optional: Add a foreign key constraint to 'job_input_files'
-- ALTER TABLE job_input_files ADD CONSTRAINT fk_job_input_files_data_file_id FOREIGN KEY (data_file_id) REFERENCES data_files(id);

-- -- Add the 'data_file_id' UUID column to 'job_output_files'
-- ALTER TABLE job_output_files ADD COLUMN data_file_id UUID;

-- -- Populate 'data_file_id' in 'job_output_files' from 'data_files'
-- UPDATE job_output_files
-- SET data_file_id = df.id
-- FROM data_files AS df
-- WHERE job_output_files.data_file_c_id = df.cid;

-- -- Optional: Add a foreign key constraint to 'job_output_files'
-- ALTER TABLE job_output_files ADD CONSTRAINT fk_job_output_files_data_file_id FOREIGN KEY (data_file_id) REFERENCES data_files(id);
