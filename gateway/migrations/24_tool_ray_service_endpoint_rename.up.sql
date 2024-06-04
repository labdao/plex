-- rename ray service url to ray service endpoint
ALTER TABLE tools ADD COLUMN ray_service_endpoint VARCHAR(255);