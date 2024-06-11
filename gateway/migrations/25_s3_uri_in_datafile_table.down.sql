BEGIN;

ALTER TABLE data_files ADD COLUMN s3_bucket VARCHAR(255);
ALTER TABLE data_files ADD COLUMN s3_location VARCHAR(255);

UPDATE data_files
SET s3_bucket = SPLIT_PART(SUBSTRING(s3_uri FROM 6), '/', 1)
WHERE s3_uri LIKE 's3://%/%';

UPDATE data_files
SET s3_location = SUBSTRING(s3_uri FROM LENGTH(SPLIT_PART(SUBSTRING(s3_uri FROM 6), '/', 1)) + 7)
WHERE s3_uri LIKE 's3://%/%';

ALTER TABLE data_files DROP COLUMN s3_uri;

COMMIT;
