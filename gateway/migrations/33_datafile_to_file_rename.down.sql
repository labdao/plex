BEGIN;

ALTER TABLE files RENAME TO data_files;
ALTER TABLE file_tags RENAME TO datafile_tags;
ALTER TABLE user_files RENAME TO user_datafiles;

ALTER TABLE datafile_tags RENAME file_c_id TO data_file_c_id;
ALTER TABLE user_datafiles RENAME file_c_id TO data_file_c_id;

ALTER TABLE job_input_files RENAME COLUMN file_c_id TO data_file_c_id;
ALTER TABLE job_output_files RENAME COLUMN file_c_id TO data_file_c_id;

ALTER TABLE data_files RENAME CONSTRAINT files_pkey TO data_files_pkey;

ALTER TABLE datafile_tags RENAME CONSTRAINT file_tags_pkey TO datafile_tags_pkey;
ALTER TABLE datafile_tags RENAME CONSTRAINT file_tags_file_c_id_fkey TO datafile_tags_data_file_c_id_fkey;
ALTER TABLE datafile_tags RENAME CONSTRAINT file_tags_tag_name_fkey TO datafile_tags_tag_name_fkey;

ALTER TABLE user_datafiles RENAME CONSTRAINT pk_user_files TO pk_user_datafiles;
ALTER TABLE user_datafiles RENAME CONSTRAINT fk_user_files_file TO fk_user_datafiles_data_file;
ALTER TABLE user_datafiles RENAME CONSTRAINT fk_user_files_wallet_address TO fk_user_datafiles_wallet_address;

ALTER TABLE job_input_files RENAME CONSTRAINT job_input_files_file_c_id_fkey TO job_input_files_data_file_c_id_fkey;
ALTER TABLE job_output_files RENAME CONSTRAINT job_output_files_file_c_id_fkey TO job_output_files_data_file_c_id_fkey;

COMMIT;