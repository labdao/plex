import os
import subprocess

from helpers.settings import PROJECT_ROOT


DEFAULT_DATA_CIDS = ['Qme6wmyrLQHurhrWBGD1jVFspMzfLLY1vgjiSa4U5kWJqV']


def load_cids_to_inputs(cids):
    for cid in cids:
        subprocess.run(f'ipfs get /ipfs/{cid} -o {os.path.join(PROJECT_ROOT, "inputs")}', shell=True)
