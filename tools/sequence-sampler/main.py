# import glob
import os
import time
import pandas as pd

import hydra
from hydra import compose, initialize
from hydra.core.hydra_config import HydraConfig
from omegaconf import DictConfig, OmegaConf

from sampler import Sampler

# def get_plex_job_inputs():
#     # Retrieve the environment variable
#     json_str = os.getenv("PLEX_JOB_INPUTS")

#     # Check if the environment variable is set
#     if json_str is None:
#         raise ValueError("PLEX_JOB_INPUTS environment variable is missing.")

#     # Convert the JSON string to a Python dictionary
#     try:
#         data = json.loads(json_str)
#         return data
#     except json.JSONDecodeError:
#         # Handle the case where the string is not valid JSON
#         raise ValueError("PLEX_JOB_INPUTS is not a valid JSON string.")

def squeeze_seq(new_sequence):
    return ''.join(filter(lambda x: x != '-', new_sequence))

def find_fasta_file(directory_path):
    for root, dirs, files in os.walk(directory_path):
        for file in files:
            if file.endswith(".fasta"):
                return os.path.abspath(os.path.join(root, file))
    return None  # Return None if no .fasta file is found in the directory

def load_initial_data(fasta_file, cfg):
    sequences = []
    with open(fasta_file, 'r') as file:
        seq_num = 1
        for line in file:
            if line.startswith('>'):
                # Add an entry for a new sequence, including the 'step' column set to 0
                sequences.append({
                    't': 0,
                    'seed': '',
                    'permissibility_seed': '',
                    '(levenshtein-distance, mask)': 'none',
                    'modified_seq': '',
                    'permissibility_modified_seq': ''}
                )
                seq_num += 1
            else:
                # Add sequence data to the most recently added sequence entry
                sequences[-1]['seed'] += line.strip()
                sequences[-1]['modified_seq'] += sequences[-1]['seed']
                sequences[-1]['permissibility_seed'] += cfg.params.basic_settings.init_permissibility_vec
                sequences[-1]['permissibility_modified_seq'] += cfg.params.basic_settings.init_permissibility_vec

    return pd.DataFrame(sequences)


@hydra.main(version_base=None, config_path="conf", config_name="config")
def my_app(cfg: DictConfig) -> None:
    # user_inputs = get_plex_job_inputs()
    # print(f"user inputs from plex: {user_inputs}")

    print(OmegaConf.to_yaml(cfg))
    print(f"Working directory : {os.getcwd()}")

    # defining output directory
    if cfg.outputs.directory is None:
        outputs_directory = hydra.core.hydra_config.HydraConfig.get().runtime.output_dir
    else:
        outputs_directory = cfg.outputs.directory
    print(f"Output directory : {outputs_directory}")

    fasta_file = find_fasta_file(cfg.inputs.directory) # load fasta with inital sequences and convert to data frame
    df = load_initial_data(fasta_file, cfg)
    print('df', df)
    seed = df.iloc[-1]['seed']
    permissibility_seed = df.iloc[-1]['permissibility_seed']

    start_time = time.time()
    print("sequence to structure complete...")

    for t in range(cfg.params.basic_settings.number_of_evo_cycles):
        print("starting evolution step", t)
        print('seed', seed)

        sampler = Sampler(t+1, seed, permissibility_seed, cfg, outputs_directory)
        mod_seq, modified_permissibility_seq, action, levenshtein_step_size, action_mask = sampler.apply_policy()

        print('mod seq', mod_seq)
        print('modified_permissibility_seq', modified_permissibility_seq)

        new_row = {
            't': t+1,
            'seed': squeeze_seq(seed),
            'permissibility_seed': ''.join(permissibility_seed),
            '(levenshtein-distance, mask)': action,
            'modified_seq': mod_seq,
            'permissibility_modified_seq': ''.join(modified_permissibility_seq)
        }

        # Append the new row to the DataFrame
        df = pd.concat([df, pd.DataFrame([new_row])], ignore_index=True)
        df.to_csv(f"{outputs_directory}/summary.csv", index=False)

        seed = mod_seq
        permissibility_seed = modified_permissibility_seq

        print('')

    print("sequence to structure complete...")
    end_time = time.time()
    duration = end_time - start_time
    print(f"executed in {duration:.2f} seconds.")

if __name__ == "__main__":
    my_app()