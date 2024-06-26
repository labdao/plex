BEGIN;

ALTER TABLE data_files RENAME TO files;
ALTER TABLE datafile_tags RENAME TO file_tags;
ALTER TABLE user_datafiles RENAME TO user_files;

ALTER TABLE file_tags RENAME data_file_c_id TO file_c_id;
ALTER TABLE user_files RENAME data_file_c_id TO file_c_id;

ALTER TABLE job_input_files RENAME COLUMN data_file_c_id TO file_c_id;
ALTER TABLE job_output_files RENAME COLUMN data_file_c_id TO file_c_id;

ALTER TABLE files RENAME CONSTRAINT data_files_pkey TO files_pkey;

ALTER TABLE file_tags RENAME CONSTRAINT datafile_tags_pkey TO file_tags_pkey;
ALTER TABLE file_tags RENAME CONSTRAINT datafile_tags_data_file_c_id_fkey TO file_tags_file_c_id_fkey;
ALTER TABLE file_tags RENAME CONSTRAINT datafile_tags_tag_name_fkey TO file_tags_tag_name_fkey;

ALTER TABLE user_files RENAME CONSTRAINT pk_user_datafiles TO pk_user_files;
ALTER TABLE user_files RENAME CONSTRAINT fk_user_datafiles_data_file TO fk_user_files_file;
ALTER TABLE user_files RENAME CONSTRAINT fk_user_datafiles_wallet_address TO fk_user_files_wallet_address;

ALTER TABLE job_input_files RENAME CONSTRAINT job_input_files_data_file_c_id_fkey TO job_input_files_file_c_id_fkey;
ALTER TABLE job_output_files RENAME CONSTRAINT job_output_files_data_file_c_id_fkey TO job_output_files_file_c_id_fkey;

COMMIT;