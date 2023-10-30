-- Create data_files table
CREATE TABLE data_files (
    cid VARCHAR(255) NOT NULL PRIMARY KEY,
    wallet_address VARCHAR(42) NOT NULL,
    filename VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP
);

-- Create flows table
CREATE TABLE flows (
    cid VARCHAR(255) NOT NULL PRIMARY KEY,
    name VARCHAR(255),
    wallet_address VARCHAR(42) NOT NULL
);

-- Create jobs table
CREATE TABLE jobs (
    bacalhau_job_id VARCHAR(255) NOT NULL PRIMARY KEY,
    state VARCHAR(255) DEFAULT 'processing',
    error TEXT DEFAULT '',
    wallet_address VARCHAR(255),
    tool_id VARCHAR(255) NOT NULL,
    flow_id VARCHAR(255) NOT NULL
);

CREATE INDEX idx_jobs_tool_id ON jobs(tool_id);
CREATE INDEX idx_jobs_flow_id ON jobs(flow_id);

-- Create tools table
CREATE TABLE tools (
    cid VARCHAR(255) NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    wallet_address VARCHAR(42) NOT NULL,
    tool_json JSON
);

-- Create users table
CREATE TABLE users (
    wallet_address VARCHAR(42) NOT NULL PRIMARY KEY
);

-- Create many-to-many relation tables
CREATE TABLE job_inputs (
    job_bacalhau_job_id VARCHAR(255) NOT NULL,
    data_file_c_id   VARCHAR(255) NOT NULL,
    PRIMARY KEY (bacalhau_job_id, cid),
    FOREIGN KEY (bacalhau_job_id) REFERENCES jobs(bacalhau_job_id),
    FOREIGN KEY (cid) REFERENCES data_files(cid)
);

CREATE TABLE job_outputs (
    job_bacalhau_job_id VARCHAR(255) NOT NULL,
    data_file_c_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (bacalhau_job_id, cid),
    FOREIGN KEY (bacalhau_job_id) REFERENCES jobs(bacalhau_job_id),
    FOREIGN KEY (cid) REFERENCES data_files(cid)
);
