WITH mismatched_wallets AS (
    SELECT df.cid, j.wallet_address
    FROM data_files df
    INNER JOIN job_output_files jof ON df.cid = jof.data_file_c_id
    INNER JOIN jobs j ON j.id = jof.job_id
    WHERE df.wallet_address != j.wallet_address
    and df.wallet_address = ''
)
UPDATE data_files df
SET wallet_address = mw.wallet_address
FROM mismatched_wallets mw
WHERE df.cid = mw.cid;