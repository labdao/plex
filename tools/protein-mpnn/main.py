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

def run_diffusion(
    contigs,
    path,
    pdb=None,
    iterations=50,
    symmetry="none",
    order=1,
    hotspot=None,
    chains=None,
    add_potential=False,
    num_designs=1,
    use_beta_model=False,
    visual="none",
    outputs_directory="outputs",
):
    print("Running diffusion with contigs:", contigs, "and path:", path)
    full_path = f"{outputs_directory}/{path}"
    os.makedirs(full_path, exist_ok=True)

    # # Set up the environment for the subprocess - required so that RFdiffussion can find its proper packages
    env = os.environ.copy()
    env['PYTHONPATH'] = "/app/RFdiffusion:" + env.get('PYTHONPATH', '')

    command = [
        'python', 'RFdiffusion/scripts/run_inference.py',
        f'inference.output_prefix={os.path.join(outputs_directory, f"design")}',
        'inference.model_directory_path=RFdiffusion/models',
        f"inference.input_pdb={pdb}",
        'inference.num_designs=2',
        f'contigmap.contigs={[contigs]}'
    ]

    print('command', command)

    logging.info(f'starting inference')
    # Run the command
    result = subprocess.run(command, capture_output=True, text=True, env=env)

    # Check if the command was successful
    if result.returncode == 0:
        print("Inference script ran successfully")
        print(result.stdout)
    else:
        print("Error running inference script")
        print(result.stderr)

    return result


@hydra.main(version_base=None, config_path="conf", config_name="config")
def my_app(cfg: DictConfig) -> None:

    user_inputs = get_plex_job_inputs()
    print(f"user inputs from plex: {user_inputs}")

    # Override Hydra default params with user supplied params
    # OmegaConf.update(cfg, "params.basic_settings.pdb", user_inputs["protein_complex"], merge=False)
    OmegaConf.update(cfg, "params.expert_settings.num_seqs", user_inputs["num_seqs"], merge=False)
    OmegaConf.update(cfg, "params.expert_settings.rm_aa", user_inputs["rm_aa"], merge=False)
    OmegaConf.update(cfg, "params.expert_settings.mpnn_sampling_temp", user_inputs["mpnn_sampling_temp"], merge=False)
    OmegaConf.update(cfg, "params.expert_settings.use_solubleMPNN", user_inputs["cuse_solubleMPNN"], merge=False)
    OmegaConf.update(cfg, "params.expert_settings.initial_guess", user_inputs["initial_guess"], merge=False)

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

    if not isinstance(input_target_path, list):
        input_target_path = [input_target_path]
    print("Identified Targets : ", input_target_path)

    # first_residue_target_chain, last_residue_target_chain = find_chain_residue_range(input_target_path[0], user_inputs["target_chain"])
    # if isinstance(user_inputs["target_start_residue"], int) and user_inputs["target_start_residue"] > 0:
    #     OmegaConf.update(cfg, "params.advanced_settings.pdb_start_residue", user_inputs["target_start_residue"], merge=False)
    # else: OmegaConf.update(cfg, "params.advanced_settings.pdb_start_residue", first_residue_target_chain, merge=False)
    
    # if isinstance(user_inputs["target_end_residue"], int) and user_inputs["target_end_residue"] > 0:
    #     OmegaConf.update(cfg, "params.advanced_settings.pdb_end_residue", user_inputs["target_end_residue"], merge=False)
    # else: OmegaConf.update(cfg, "params.advanced_settings.pdb_end_residue", last_residue_target_chain, merge=False)


    # running ProteinMPNN


        flags = {
            "contigs": contigs,
            "pdb": pdb,
            "order": order,
            "iterations": iterations,
            "symmetry": symmetry,
            "hotspot": hotspot,
            "path": path,
            "chains": chains,
            "add_potential": add_potential,
            "num_designs": num_designs,
            "use_beta_model": use_beta_model,
            "visual": visual,
            "outputs_directory": outputs_directory
        }

        print("outputs_directory", outputs_directory)

        for k, v in flags.items():
            if isinstance(v, str):
                flags[k] = v.replace("'", "").replace('"', "")

        run_diffusion(**flags)

        logging.info("design complete...")
        end_time = time.time()
        duration = end_time - start_time
        logging.info(f"executed in {duration:.2f} seconds.")

if __name__ == "__main__":
    my_app()
