CREATE TABLE user_datafiles (
    wallet_address varchar(42) NOT NULL,
    c_id varchar(255) NOT NULL,
    created_at timestamp,
    CONSTRAINT pk_user_datafiles PRIMARY KEY (wallet_address, c_id),
    CONSTRAINT fk_user_datafiles_wallet_address FOREIGN KEY (wallet_address) REFERENCES users(wallet_address),
    CONSTRAINT fk_user_datafiles_data_file FOREIGN KEY (c_id) REFERENCES data_files(cid)
);

INSERT INTO user_datafiles (wallet_address, c_id, created_at)
SELECT wallet_address, cid, COALESCE(timestamp, CURRENT_TIMESTAMP)
FROM data_files;

-- below steps will be done as a separate migration after we test and make sure the above information has been copied over correctly
-- ALTER TABLE data_files
-- DROP COLUMN wallet_address,
-- DROP COLUMN timestamp;
