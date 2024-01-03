# import glob
import os
import time
import pandas as pd

from AF2_module import AF2Runner
import hydra
from hydra import compose, initialize
from hydra.core.hydra_config import HydraConfig
from omegaconf import DictConfig, OmegaConf

from agent import Agent
from oracle import Oracle


def find_fasta_file(directory_path):
    for root, dirs, files in os.walk(directory_path):
        for file in files:
            if file.endswith(".fasta"):
                return os.path.abspath(os.path.join(root, file))
    return None  # Return None if no .fasta file is found in the directory

def load_fasta_to_dataframe(fasta_file, cfg):
    sequences = []
    with open(fasta_file, 'r') as file:
        seq_num = 1
        for line in file:
            if line.startswith('>'):
                # Add an entry for a new sequence, including the 'step' column set to 0
                sequences.append({'t': 0, 'sequence_number': seq_num, 'seq': ''}) # TD: consider changing name seq to original_seq here already
                seq_num += 1
            else:
                # Add sequence data to the most recently added sequence entry
                sequences[-1]['seq'] += line.strip()

    return pd.DataFrame(sequences)

def load_initial_data(fasta_file, cfg):
    sequences = []
    with open(fasta_file, 'r') as file:
        seq_num = 1
        for line in file:
            if line.startswith('>'):
                # Add an entry for a new sequence, including the 'step' column set to 0
                sequences.append({'t': 0, 'sequence_number': seq_num, 'seq': '', 'permissibility_vectors': ''}) # TD: consider changing name seq to original_seq here already
                seq_num += 1
            else:
                # Add sequence data to the most recently added sequence entry
                sequences[-1]['seq'] += line.strip()
                sequences[-1]['permissibility_vectors'] += cfg.params.basic_settings.init_permissibility_vec
    
    # After the file reading loop
    for sequence in sequences:
        # Convert the string to a list of characters and update the dictionary
        # sequence['permissibility_vectors'] = [vec for vec in sequence['permissibility_vectors']]
        sequence['permissibility_vectors'] = list(sequence['permissibility_vectors'])

    return pd.DataFrame(sequences)

def step(t, df, df_action, outputs_directory, cfg):

    # # run oracle
    oracle_runner = Oracle(t, df, df_action, outputs_directory, cfg)
    df = oracle_runner.run()

    # # run reward
    reward = 0
    
    return df, reward


@hydra.main(version_base=None, config_path="conf", config_name="config_sequence-optimizer")
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
    # df_0 = load_fasta_to_dataframe(fasta_file, cfg)
    df_0 = load_initial_data(fasta_file, cfg)

    start_time = time.time()
    print("sequence to structure complete...")

    reward = 0
    df, reward_step = step(0, df_0, df_0, outputs_directory, cfg)
    for t in range(cfg.params.basic_settings.number_of_evo_cycles):
        print("starting iteration number ", t)

        agent = Agent(t+1, df, reward, cfg)
        df, df_action = agent.apply_policy()
        df, reward_step = step(t+1, df, df_action, outputs_directory, cfg)

        reward = reward_step
        df.to_csv(f"{outputs_directory}/summary.csv", index=False)

    # t_rows = df[df['t'] == cfg.params.basic_settings.number_of_evo_cycles]
    # best_
    # print('df', df)
    # df.to_csv(f"{outputs_directory}/summary.csv", index=False)

    print("sequence to structure complete...")
    end_time = time.time()
    duration = end_time - start_time
    print(f"executed in {duration:.2f} seconds.")

if __name__ == "__main__":
    my_app()