-- This just resets the DB to a clean state
-- Clear your schema_migrations (truncate schema_migrations;) table if you want to run this again
-- This should be removed before we do stable stage/production deploys

DROP TABLE IF EXISTS tools;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS jobs_id_seq;
DROP TABLE IF EXISTS tool_entities;
DROP TABLE IF EXISTS users_id_seq;
DROP TABLE IF EXISTS data_files;
DROP TABLE IF EXISTS graphs;
