-- rename ray service url to ray service endpoint
ALTER TABLE tools RENAME COLUMN ray_service_url TO ray_service_endpoint;