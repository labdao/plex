# import glob
import os
import time
import pandas as pd

import hydra
from hydra import compose, initialize
from hydra.core.hydra_config import HydraConfig
from omegaconf import DictConfig, OmegaConf

from sampler import Sampler
from utils import slash_to_convexity_notation

import json

import logging
import random

def get_plex_job_inputs():
    # Retrieve the environment variable
    json_str = os.getenv("PLEX_JOB_INPUTS")

    # Check if the environment variable is set
    if json_str is None:
        raise ValueError("PLEX_JOB_INPUTS environment variable is missing.")

    # Convert the JSON string to a Python dictionary
    try:
        data = json.loads(json_str)
        return data
    except json.JSONDecodeError:
        # Handle the case where the string is not valid JSON
        raise ValueError("PLEX_JOB_INPUTS is not a valid JSON string.")

def squeeze_seq(new_sequence):
    return ''.join(filter(lambda x: x != '-', new_sequence))

def find_fasta_file(directory_path):
    for root, dirs, files in os.walk(directory_path):
        for file in files:
            if file.endswith(".fasta"):
                return os.path.abspath(os.path.join(root, file))
    return None  # Return None if no .fasta file is found in the directory

def apply_initial_permissibility_vector(seed, permissibility_seed, cfg):

    mod_sequence = []
    seed_list = list(seed)

    for i, char in enumerate(permissibility_seed):
        if char == 'X' or char == '+' or char in cfg.params.basic_settings.alphabet:
            mod_sequence.append(seed_list[i])
        # elif char == '-':
        #     mod_sequence.append('-')      
    
    mod_sequence = ''.join(mod_sequence)

    return mod_sequence

def load_initial_data(fasta_file, cfg):
    sequences = []
    with open(fasta_file, 'r') as file:
        seq_num = 1
        for line in file:
            if line.startswith('>'):
                # Add an entry for a new sequence, including the 'step' column set to 0
                sequences.append({
                    't': 0,
                    'sample_number': 0,
                    'seed': '',
                    'permissibility_seed': '',
                    '(levenshtein-distance, mask)': 'none',
                    'modified_seq': '',
                    'permissibility_modified_seq': '',
                    'acceptance_flag': False}
                )
                seq_num += 1
            else:
                # Add sequence data to the most recently added sequence entry
                sequences[-1]['seed'] += line.strip()
                contig_in_convexity_notation = slash_to_convexity_notation(sequences[-1]['seed'], cfg.params.basic_settings.init_permissibility_vec)
                sequences[-1]['modified_seq'] += apply_initial_permissibility_vector(sequences[-1]['seed'], contig_in_convexity_notation, cfg)
                sequences[-1]['permissibility_seed'] += contig_in_convexity_notation
                sequences[-1]['permissibility_modified_seq'] += contig_in_convexity_notation

    sequences[0]['acceptance_flag'] = True # manual selection of starting sequence

    return pd.DataFrame(sequences)


@hydra.main(version_base=None, config_path="conf", config_name="config")
def my_app(cfg: DictConfig) -> None:

    # defining output directory
    if cfg.outputs.directory is None:
        outputs_directory = hydra.core.hydra_config.HydraConfig.get().runtime.output_dir
    else:
        outputs_directory = cfg.outputs.directory

    ## plex user inputs # some of these are currently not used!
    user_inputs = get_plex_job_inputs()
    permissibility_seed = user_inputs["init_permissibility_vec"]
    logging.info(f"user inputs from plex: {user_inputs}")

    # # Override Hydra default params with user supplied params
    OmegaConf.update(cfg, "params.basic_settings.experiment_name", user_inputs["experiment_name"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.number_of_evo_cycles", user_inputs["number_of_evo_cycles"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.policy_flag", user_inputs["policy_flag"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.init_permissibility_vec", user_inputs["init_permissibility_vec"], merge=False)
    # OmegaConf.update(cfg, "params.basic_settings.temperature", user_inputs["temperature"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.max_levenshtein_step_size", user_inputs["max_levenshtein_step_size"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.alphabet", user_inputs["alphabet"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.scorers", user_inputs["scorers"], merge=False)
    # OmegaConf.update(cfg, "params.basic_settings.scoring_metrics", user_inputs["scoring_metrics"], merge=False)
    # OmegaConf.update(cfg, "params.basic_settings.scoring_weights", user_inputs["scoring_weights"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.bouncer_flag", user_inputs["bouncer_flag"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.generators", user_inputs["generators"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.target_template_complex", user_inputs["target_template_complex"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.target_chain", user_inputs["target_chain"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.binder_chain", user_inputs["binder_chain"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.target_seq", user_inputs["target_seq"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.target_pdb", user_inputs["target_pdb"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.binder_template_sequence", user_inputs["binder_template_sequence"], merge=False)

    logging.info(f"{OmegaConf.to_yaml(cfg)}")
    logging.info(f"Working directory : {os.getcwd()}")

    logging.info(f"inputs directory: {cfg.inputs.directory}")
    # fasta_file = find_fasta_file(cfg.inputs.directory) # load fasta with inital sequences and convert to data frame
    fasta_file = os.path.abspath(os.path.join(cfg.inputs.directory, cfg.params.basic_settings.binder_template_sequence))
    logging.info(f"fasta file {fasta_file}")

    df = load_initial_data(fasta_file, cfg)

    seed_row = df[(df['t']==0) & (df['acceptance_flag'] == True)]
    seed = seed_row['modified_seq'].values[0]
    logging.info(f"initial seed sequence {seed}")
    permissibility_seed = seed_row['permissibility_modified_seq'].values[0]

    logging.info("initial sequence to structure complete...")


    start_time = time.time()
    for t in range(cfg.params.basic_settings.number_of_evo_cycles):

        logging.info(f"starting evolution step, {t+1}")
        logging.info(f"seed, {seed}")

        sampler = Sampler(t+1, seed, permissibility_seed, cfg, outputs_directory, df)
        mod_seq, modified_permissibility_seq, action, levenshtein_step_size, action_mask, df = sampler.apply_policy()

        logging.info(f"mod seq, {mod_seq}")
        # logging.info(f"modified_permissibility_seq, {modified_permissibility_seq}")

        df.to_csv(f"{outputs_directory}/summary.csv", index=False)

        # update seed and permissibility seed
        seed = mod_seq
        permissibility_seed = modified_permissibility_seq

        print('\n')

    logging.info("sequence to structure complete...")
    end_time = time.time()
    duration = end_time - start_time
    logging.info(f"executed in {duration:.2f} seconds.")

if __name__ == "__main__":
    my_app()
