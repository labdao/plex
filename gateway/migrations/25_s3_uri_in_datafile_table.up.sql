BEGIN;

ALTER TABLE data_files ADD COLUMN s3_uri VARCHAR(255);

UPDATE data_files
SET s3_uri = CONCAT(
  's3://', 
  s3_bucket, 
  CASE
    WHEN LEFT(s3_location, 1) != '/' THEN '/'
    ELSE ''
  END,
  s3_location
)
WHERE s3_bucket IS NOT NULL AND s3_location IS NOT NULL;

ALTER TABLE data_files DROP COLUMN s3_bucket;
ALTER TABLE data_files DROP COLUMN s3_location;

COMMIT;
