# import glob
import os
import time
import pandas as pd
import json

import hydra
from hydra import compose, initialize
from hydra.core.hydra_config import HydraConfig
from omegaconf import DictConfig, OmegaConf

from Bio import PDB
import subprocess

import logging

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

def find_chain_residue_range(pdb_path, chain_id):
    """
    Finds the start and end residue sequence indices for a given chain in a PDB file.
    """
    parser = PDB.PDBParser(QUIET=True)
    structure = parser.get_structure('protein', pdb_path)
    for model in structure:
        for chain in model:
            if chain.id == chain_id:
                residues = list(chain)
                if residues:
                    start_residue = residues[0].id[1]
                    end_residue = residues[-1].id[1]
                    return start_residue, end_residue
    return None, None

def get_files_from_directory(root_dir, extension, max_depth=3):
    pdb_files = []

    for root, dirs, files in os.walk(root_dir):
        depth = root[len(root_dir) :].count(os.path.sep)

        if depth <= max_depth:
            for f in files:
                if f.endswith(extension):
                    pdb_files.append(os.path.join(root, f))

            # Prune the directory list if we are at max_depth
            if depth == max_depth:
                del dirs[:]
    print(
        "Found {} files with extension {} in directory {}".format(
            len(pdb_files), extension, root_dir
        )
    )
    return pdb_files

@hydra.main(version_base=None, config_path="conf", config_name="config")
def my_app(cfg: DictConfig) -> None:

    user_inputs = get_plex_job_inputs()
    print(f"user inputs from plex: {user_inputs}")

    print(type(user_inputs["mpnn_sampling_temp"]))

    # Override Hydra default params with user supplied params
    OmegaConf.update(cfg, "params.expert_settings.num_seqs", user_inputs["num_seqs"], merge=False)
    OmegaConf.update(cfg, "params.expert_settings.rm_aa", user_inputs["rm_aa"], merge=False)
    OmegaConf.update(cfg, "params.expert_settings.mpnn_sampling_temp", user_inputs["mpnn_sampling_temp"], merge=False)
    OmegaConf.update(cfg, "params.expert_settings.use_solubleMPNN", user_inputs["use_solubleMPNN"], merge=False)
    OmegaConf.update(cfg, "params.expert_settings.initial_guess", user_inputs["initial_guess"], merge=False)
    OmegaConf.update(cfg, "params.expert_settings.chains_to_design", user_inputs["chains_to_design"], merge=False)

    print(OmegaConf.to_yaml(cfg))
    print(f"Working directory : {os.getcwd()}")

    # defining output directory
    if cfg.outputs.directory is None:
        outputs_directory = hydra.core.hydra_config.HydraConfig.get().runtime.output_dir
    else:
        outputs_directory = cfg.outputs.directory
    print(f"Output directory : {outputs_directory}")

    # defining input files
    if user_inputs.get("protein_complex"):
        input_target_path = user_inputs["protein_complex"]
    else:
        input_target_path = get_files_from_directory(cfg.inputs.target_directory, ".pdb")

        if cfg.inputs.target_pattern is not None:
            input_target_path = [
                file for file in input_target_path if cfg.inputs.target_pattern in file
            ]

    # if not isinstance(input_target_path, list):
    #     input_target_path = [input_target_path]
    print("Identified complex: ", input_target_path)


    num_seqs = str(cfg.params.expert_settings.num_seqs)
    print('here', type(num_seqs))
    rm_aa = cfg.params.expert_settings.rm_aa
    mpnn_sampling_temp = str(cfg.params.expert_settings.mpnn_sampling_temp)
    use_solubleMPNN = cfg.params.expert_settings.use_solubleMPNN
    initial_guess = cfg.params.expert_settings.initial_guess
    chains_to_design = cfg.params.expert_settings.chains_to_design

    start_time = time.time()

    logging.info(f'Running MPNN')

    # Activate the conda environment 'mlfold'
    subprocess.run(['conda', 'activate', 'mlfold'], shell=True)

    print("pdb path", input_target_path)
    print("output directory", outputs_directory)

    # Define the command and arguments
    command = [
        'python', 'ProteinMPNN/protein_mpnn_run.py',
        '--pdb_path', input_target_path,
        '--pdb_path_chains', chains_to_design,
        '--out_folder', outputs_directory,
        '--num_seq_per_target', num_seqs,
        '--sampling_temp', mpnn_sampling_temp,
        '--seed', '37',
        '--batch_size', '1'
    ]

    # Run the command
    result = subprocess.run(command, capture_output=True, text=True)

    # Print the output
    print(result.stdout)
    print(result.stderr)  

    logging.info("MPNN complete...")
    end_time = time.time()
    duration = end_time - start_time
    logging.info(f"executed in {duration:.2f} seconds.")

if __name__ == "__main__":
    my_app()


    # # Define the command and arguments
    # command = [
    #         'python', 'ProteinMPNN/protein_mpnn_run.py',
    #         '--pdb_path', input_target_path,
    #         '--pdb_path_chains', chains_to_design,
    #         '--out_folder', outputs_directory,
    #         '--num_seq_per_target', num_seqs,
    #         '--sampling_temp', mpnn_sampling_temp,
    #         '--seed', '37',
    #         '--batch_size', '1'
    #     ]

    # # Run the command
    # subprocess.run(command, capture_output=True, text=True)  
