import os
from AF2_module import AF2Runner

def write_dataframe_to_fastas(t, dataframe, cfg):
    input_dir = os.path.join(cfg.inputs.directory, 'current_sequences')
    if os.path.exists(input_dir):
        # If the folder already exists, empty the folder of all files
        for file_name in os.listdir(input_dir):
            file_path = os.path.join(input_dir, file_name)
            if os.path.isfile(file_path):
                os.remove(file_path)
    else:
        os.makedirs(input_dir, exist_ok=True)

    for index, row in dataframe.iterrows():
        file_path = os.path.join(input_dir, f"seq_{row['sequence_number']}_t_{t}.fasta")
        with open(file_path, 'w') as file:
            file.write(f">{row['sequence_number']}\n{row['seq']}\n")
    return os.path.abspath(input_dir)

import os
import json
import pandas as pd
import numpy as np
import glob

def supplement_dataframe(t, df, directory_path):
    # Ensure DataFrame has a proper index
    df.reset_index(drop=True, inplace=True)

    # Adding new columns only if they don't already exist
    for col in ['top rank json', 'top rank pdb', 'plddt', 'max_pae', 'ptm']:
        if col not in df.columns:
            df[col] = np.nan

    for index, row in df.iterrows():
        if row['t'] == t:
            sequence_number = row['sequence_number']

            # Search for the matching json and pdb files
            json_pattern = os.path.join(directory_path, f'seq_{sequence_number}_t_{t}_scores_rank_001*.json')
            pdb_pattern = os.path.join(directory_path, f'seq_{sequence_number}_t_{t}_unrelaxed_rank_001*.pdb')

            json_files = glob.glob(json_pattern)
            pdb_files = glob.glob(pdb_pattern)

            # Assuming the first match is the desired file
            if json_files:
                df.at[index, 'top rank json'] = json_files[0]
                with open(json_files[0], 'r') as file:
                    data = json.load(file)
                    df.at[index, 'plddt'] = np.mean(data['plddt'])
                    df.at[index, 'max_pae'] = data['max_pae']
                    df.at[index, 'ptm'] = data['ptm']

            if pdb_files:
                df.at[index, 'top rank pdb'] = pdb_files[0]

    return df

# def supplement_dataframe(t, df, directory_path):
#     # Adding new columns only if they don't already exist
#     for col in ['top rank json', 'top rank pdb', 'plddt', 'max_pae', 'ptm']:
#         if col not in df.columns:
#             df[col] = np.nan

#     for index, row in df.iterrows():
#         if row['t'] == t:
#             sequence_number = row['sequence_number']

#             # Search for the matching json and pdb files
#             json_pattern = os.path.join(directory_path, f'seq_{sequence_number}_t_{t}_scores_rank_001*.json')
#             pdb_pattern = os.path.join(directory_path, f'seq_{sequence_number}_t_{t}_unrelaxed_rank_001*.pdb')

#             json_files = glob.glob(json_pattern)
#             pdb_files = glob.glob(pdb_pattern)

#             # Assuming the first match is the desired file
#             if json_files:
#                 df.at[index, 'top rank json'] = json_files[0]
#                 with open(json_files[0], 'r') as file:
#                     data = json.load(file)
#                     df.at[index, 'plddt'] = np.mean(data['plddt'])
#                     df.at[index, 'max_pae'] = data['max_pae']
#                     df.at[index, 'ptm'] = data['ptm']

#             if pdb_files:
#                 df.at[index, 'top rank pdb'] = pdb_files[0]

#     return df

class Oracle:
    def __init__(self, t, df, df_action, outputs_directory, cfg):

        self.t = t
        self.df = df
        self.df_action = df_action
        self.outputs_directory = outputs_directory
        self.cfg = cfg

    def run(self):

        # prepare input sequences as fastas and run AF2 K-times
        seq_input_dir = write_dataframe_to_fastas(self.t, self.df_action, self.cfg)

        K = self.cfg.params.basic_settings.AF2_repeats_per_seq
        for n in range(K):
            print("starting repeat number ", n)
            af2_runner = AF2Runner(seq_input_dir, self.outputs_directory)
            af2_runner.run()

        # complete df data frame with info
        supplemented_dataframe = supplement_dataframe(self.t, self.df, self.outputs_directory)

        return supplemented_dataframe