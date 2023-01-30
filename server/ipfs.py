import subprocess


DEFAULT_DATA_CIDS = ['QmRJrFNe6jfmiiUzjhEH4UWTivaHmpwR2Ryenu7k9aeLTM']


def load_cids_to_inputs(cids):
    for cid in cids:
        subprocess.run(f'ipfs get /ipfs/{cid}', shell=True)
