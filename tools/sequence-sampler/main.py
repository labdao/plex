# import glob
import os
import time
import pandas as pd

# from AF2_module import AF2Runner
import hydra
from hydra import compose, initialize
from hydra.core.hydra_config import HydraConfig
from omegaconf import DictConfig, OmegaConf

from sampler import Sampler
# from oracle import Oracle

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
                sequences.append({'t': 0, 'seed_number': seq_num, 'seed': '', 'modified_seq': '', 'permissibility_seed': '', 'permissibility_modified_seq': ''})
                seq_num += 1
            else:
                # Add sequence data to the most recently added sequence entry
                sequences[-1]['seed'] += line.strip()
                sequences[-1]['modified_seq'] += sequences[-1]['seed']
                sequences[-1]['permissibility_seed'] += cfg.params.basic_settings.init_permissibility_vec
                sequences[-1]['permissibility_modified_seq'] += cfg.params.basic_settings.init_permissibility_vec
    
    # After the file reading loop
    for sequence in sequences:
        # Convert the string to a list of characters and update the dictionary
        # sequence['permissibility_vectors'] = [vec for vec in sequence['permissibility_vectors']]
        sequence['permissibility_seed'] = list(sequence['permissibility_seed'])
        sequence['permissibility_modified_seq'] = list(sequence['permissibility_modified_seq'])

    return pd.DataFrame(sequences)


@hydra.main(version_base=None, config_path="conf", config_name="config")
def my_app(cfg: DictConfig) -> None:

    print(OmegaConf.to_yaml(cfg))
    print(f"Working directory : {os.getcwd()}")

    # defining output directory
    if cfg.outputs.directory is None:
        outputs_directory = hydra.core.hydra_config.HydraConfig.get().runtime.output_dir
    else:
        outputs_directory = cfg.outputs.directory
    print(f"Output directory : {outputs_directory}")

    fasta_file = find_fasta_file(cfg.inputs.directory) # load fasta with inital sequences and convert to data frame
    df_0 = load_initial_data(fasta_file, cfg)
    seed = df_0[-1]['seed']
    permissibility_seed = df_0[-1]['permissibility_seed']
    print('df0', df_0)

    start_time = time.time()
    print("sequence to structure complete...")

    for t in range(cfg.params.basic_settings.number_of_evo_cycles):
        print("starting evolution step", t)
        print('seed', seed)

        sampler = Sampler(t+1, seed, permissibility_seed, cfg)
        mod_seq, modified_permissibility_seq = sampler.apply_policy()
        seed = mod_seq

        print('mod seq', mod_seq)
        print('modified_permissibility_seq', modified_permissibility_seq)

        # df.to_csv(f"{outputs_directory}/summary.csv", index=False)

    print("sequence to structure complete...")
    end_time = time.time()
    duration = end_time - start_time
    print(f"executed in {duration:.2f} seconds.")

if __name__ == "__main__":
    my_app()


### OLD CODE ###
# ...