BEGIN;

CREATE TABLE IF NOT EXISTS organizations (
    id SERIAL PRIMARY KEY,
    name VARCHAR,
    description VARCHAR
);

ALTER TABLE users ADD COLUMN id SERIAL;

ALTER TABLE api_keys DROP CONSTRAINT IF EXISTS fk_user;
ALTER TABLE user_files DROP CONSTRAINT IF EXISTS fk_user_files_wallet_address;
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS transactions_user_id_fkey;

ALTER TABLE users DROP CONSTRAINT users_pkey;
ALTER TABLE users ADD PRIMARY KEY (id);

ALTER TABLE users ADD COLUMN organization_id INT;
ALTER TABLE users ADD FOREIGN KEY (organization_id) REFERENCES organizations(id);

ALTER TABLE experiments DROP COLUMN IF EXISTS experiment_uuid;
ALTER TABLE experiments DROP COLUMN IF EXISTS wallet_address;
ALTER TABLE experiments DROP COLUMN IF EXISTS start_time;
ALTER TABLE experiments DROP COLUMN IF EXISTS end_time;

ALTER TABLE experiments ADD COLUMN IF NOT EXISTS user_id INT;
ALTER TABLE experiments ADD COLUMN IF NOT EXISTS last_modified_at TIMESTAMP;
ALTER TABLE experiments ADD FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE jobs DROP COLUMN IF EXISTS wallet_address;
ALTER TABLE jobs DROP COLUMN IF EXISTS queue;
ALTER TABLE jobs DROP COLUMN IF EXISTS job_type;
ALTER TABLE jobs DROP COLUMN IF EXISTS result_json;
ALTER TABLE jobs DROP COLUMN IF EXISTS annotations;

--rename state to job_status
ALTER TABLE jobs RENAME COLUMN state TO job_status;

ALTER TABLE jobs ADD COLUMN IF NOT EXISTS last_modified_at TIMESTAMP;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS user_id INT;
ALTER TABLE jobs ADD FOREIGN KEY (user_id) REFERENCES users(id);

--cid is a primary key
ALTER TABLE models DROP CONSTRAINT models_pkey;
ALTER TABLE models DROP COLUMN IF EXISTS cid;
ALTER TABLE models DROP COLUMN IF EXISTS wallet_address;
ALTER TABLE models DROP COLUMN IF EXISTS model_type;
ALTER TABLE models DROP COLUMN IF EXISTS s3_uri;
ALTER TABLE models DROP COLUMN IF EXISTS container;
ALTER TABLE models DROP COLUMN IF EXISTS memory;
ALTER TABLE models DROP COLUMN IF EXISTS cpu;
ALTER TABLE models DROP COLUMN IF EXISTS gpu;
ALTER TABLE models DROP COLUMN IF EXISTS network;
ALTER TABLE models DROP COLUMN IF EXISTS timestamp;

ALTER TABLE models ADD COLUMN IF NOT EXISTS created_at TIMESTAMP;
ALTER TABLE models ADD COLUMN id SERIAL PRIMARY KEY;
ALTER TABLE models ADD COLUMN user_id INT;
ALTER TABLE models ADD FOREIGN KEY (user_id) REFERENCES users(id);

--files pkey is referenced by user_files, job_output_files, and job_input_files and file_tags drop them first before droping pkey
ALTER TABLE file_tags DROP CONSTRAINT file_tags_file_c_id_fkey;
ALTER TABLE job_input_files DROP CONSTRAINT job_input_files_file_c_id_fkey;
ALTER TABLE user_files DROP CONSTRAINT fk_user_files_file;
ALTER TABLE job_output_files DROP CONSTRAINT job_output_files_file_c_id_fkey;

ALTER TABLE job_input_files DROP CONSTRAINT IF EXISTS  fk_job_input_files_file;
ALTER TABLE job_output_files DROP CONSTRAINT IF EXISTS fk_job_output_files_file;

ALTER TABLE job_input_files RENAME COLUMN file_c_id to file_id;
ALTER TABLE job_output_files RENAME COLUMN file_c_id TO file_id;

ALTER TABLE job_input_files
ALTER COLUMN file_id TYPE int USING file_id::int;

ALTER TABLE job_output_files
ALTER COLUMN file_id TYPE int USING file_id::int;

ALTER TABLE files DROP CONSTRAINT files_pkey;
ALTER TABLE files ALTER COLUMN cid DROP NOT NULL;
ALTER TABLE files DROP COLUMN cid;
ALTER TABLE files DROP COLUMN IF EXISTS timestamp;
ALTER TABLE files DROP COLUMN IF EXISTS wallet_address;

ALTER TABLE files ADD COLUMN id SERIAL PRIMARY KEY;
ALTER TABLE files ADD COLUMN file_hash VARCHAR;
ALTER TABLE files ADD COLUMN created_at TIMESTAMP;
ALTER TABLE files ADD COLUMN last_modified_at TIMESTAMP;
ALTER TABLE files ADD COLUMN user_id INT;
ALTER TABLE files ADD FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS revoked_at TIMESTAMP;
ALTER TABLE api_keys DROP CONSTRAINT IF EXISTS api_keys_user_id_fkey;
ALTER TABLE api_keys DROP COLUMN IF EXISTS user_id;
ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS user_id INT;
ALTER TABLE api_keys ADD FOREIGN KEY (user_id) REFERENCES users(id);

CREATE TABLE IF NOT EXISTS designs (
    id SERIAL PRIMARY KEY,
    job_id INT,
    x_axis_value VARCHAR,
    y_axis_value VARCHAR,
    checkpoint_pdb_id INT,
    additional_details JSON,
    FOREIGN KEY (job_id) REFERENCES jobs(id),
    FOREIGN KEY (checkpoint_pdb_id) REFERENCES files(id)
);

DROP TABLE IF EXISTS request_trackers;

CREATE TABLE IF NOT EXISTS inference_events (
    id SERIAL PRIMARY KEY,
    job_id INT,
    ray_job_id VARCHAR,
    input_json JSON,
    output_json JSON,
    retry_count INT,
    job_status VARCHAR,
    response_code INT,
    event_time TIMESTAMP,
    error_message VARCHAR,
    event_type VARCHAR,
    FOREIGN KEY (job_id) REFERENCES jobs(id)
);

CREATE TABLE IF NOT EXISTS user_events (
    id SERIAL PRIMARY KEY,
    user_id INT,
    api_key_id INT,
    event_time TIMESTAMP,
    event_type VARCHAR,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (api_key_id) REFERENCES api_keys(id)
);

CREATE TABLE IF NOT EXISTS file_events (
    id SERIAL PRIMARY KEY,
    file_id INT,
    user_id INT,
    event_time TIMESTAMP,
    event_type VARCHAR,
    FOREIGN KEY (file_id) REFERENCES files(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);


COMMIT;
