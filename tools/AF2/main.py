import os
import time
import pandas as pd
import mdtraj as md
import numpy as np
from AF2_module import AF2Runner
import csv
from Bio.PDB import PDBParser
from Bio.SeqUtils import seq1

import hydra
from hydra import compose, initialize
from hydra.core.hydra_config import HydraConfig
from omegaconf import DictConfig, OmegaConf
import glob
import json
import subprocess

def compute_affinity(file_path):
    if pd.notna(file_path):
        try:
            # Run Prodigy and capture the output in temp.txt
            subprocess.run(
                ["prodigy", "-q", file_path], stdout=open("temp.txt", "w"), check=True
            )
            # Read the output from temp.txt
            with open("temp.txt", "r") as f:
                lines = f.readlines()
                if lines:  # Check if lines is not empty
                    # Extract the affinity value from the output
                    affinity = float(lines[0].split(" ")[-1].split("/")[0])
                    return affinity
                else:
                    print(f"No output from prodigy for {file_path}")
                    return None  # No output from Prodigy
        except subprocess.CalledProcessError:
            print(f"Prodigy command failed for {file_path}")
            return None  # Prodigy command failed
    else:
        print("Invalid file path")
        return None  # Invalid file path provided

def extract_sequence_from_pdb(pdb_file):
    parser = PDBParser()
    structure = parser.get_structure('structure', pdb_file)
    sequences = []

    for model in structure:
        for chain in model:
            seq = []
            for residue in chain:
                if residue.id[0] == ' ':  # This checks for hetero/water residues
                    seq.append(seq1(residue.resname.strip()))
            sequences.append(''.join(seq))

    return ':'.join(sequences)

def compute_ipae(pdb_file, pae_matrix):

    traj = md.load(pdb_file)

    distance_cutoff = 0.35
    contacts = md.compute_contacts(traj, contacts='all', scheme='closest-heavy')

    # Get the contact pairs and distances
    contact_distances = contacts[0].reshape(-1)  # flatten the distances array
    contact_pairs = contacts[1]

    # Filter contact pairs based on the distance cutoff
    contact_indices = contact_pairs[contact_distances < distance_cutoff]

    # extract interface values
    interface_pae_values = pae_matrix[contact_indices[:, 0], contact_indices[:, 1]]

    # Compute the median PAE for the interface contacts
    median_pae_interface = np.median(interface_pae_values)

    # Output the median PAE value
    print("median PAE at the interface:", median_pae_interface)

    return median_pae_interface

def create_summary(directory, json_pattern):
    summary_file = os.path.join(directory, 'updated_summary.csv')

    # Initialize an empty DataFrame
    df = pd.DataFrame()

    # loop over JSON files that match the given pattern for the current iteration
    for json_file in glob.glob(os.path.join(directory, f"{json_pattern}_scores*.json")):
        with open(json_file, 'r') as file:
            data = json.load(file)

        # Compute average plddt
        avg_plddt = sum(data['plddt']) / len(data['plddt'])
        # Get max_pae value
        max_pae = data['max_pae']

        # Find corresponding PDB file
        pdb_file = None
        rank_str = json_file.split('rank')[1].split('.')[0]
        for pdb in glob.glob(os.path.join(directory, f"{json_pattern}_unrelaxed_rank{rank_str}*.pdb")):
            pdb_file = pdb
            break

        pae_matrix = np.array(data['pae'])
        i_pae = compute_ipae(os.path.abspath(pdb_file), pae_matrix)

        sequence = extract_sequence_from_pdb(os.path.abspath(pdb_file))

        # Usage example
        pdb_file_path = os.path.abspath(pdb_file)
        affinity = compute_affinity(pdb_file_path)
        if affinity is not None:
            print(f"The affinity for the file {pdb_file_path} is {affinity}")

        # Prepare new row
        new_row = {
            'sequence': sequence,
            'mean plddt': avg_plddt,
            'max pae': max_pae,
            'i_pae': i_pae,
            'affinity': affinity,
            'absolute json path': os.path.abspath(json_file),
            'absolute pdb path': os.path.abspath(pdb_file) if pdb_file else None
        }

        # Concatenate new row to DataFrame
        df = pd.concat([df, pd.DataFrame([new_row])], ignore_index=True)

    if not df.empty:
        if os.path.exists(summary_file):
            df.to_csv(summary_file, mode='a', header=False, index=False)
        else:
            df.to_csv(summary_file, index=False)

    return df

def write_dataframe_to_fastas(dataframe, cfg):
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
        file_path = os.path.join(input_dir, f"seq_{row['sequence_number']}.fasta")
        with open(file_path, 'w') as file:
            file.write(f">{row['sequence_number']}\n{row['seq']}\n")
    return os.path.abspath(input_dir)


def find_fasta_file(directory_path):
    for root, dirs, files in os.walk(directory_path):
        for file in files:
            if file.endswith(".fasta"):
                return os.path.abspath(os.path.join(root, file))
    return None  # Return None if no .fasta file is found in the directory


def load_fasta_to_dataframe(fasta_file):
    sequences = []
    with open(fasta_file, 'r') as file:
        seq_num = 1
        for line in file:
            if line.startswith('>'):
                sequences.append({'sequence_number': seq_num, 'seq': ''})
                seq_num += 1
            else:
                sequences[-1]['seq'] += line.strip()

    return pd.DataFrame(sequences)

def seq2struc(df, outputs_directory, cfg):

    seq_input_dir = write_dataframe_to_fastas(df, cfg)

    af2_runner = AF2Runner(seq_input_dir, outputs_directory)
    af2_runner.run()
        
    return None

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
    df = load_fasta_to_dataframe(fasta_file)

    start_time = time.time()

    seq2struc(df, outputs_directory, cfg)

    # create and write a csv file with sequence and metric information for each output structure
    create_summary(outputs_directory, json_pattern='seq_1')

    print("Sequence to structure complete...")
    end_time = time.time()
    duration = end_time - start_time
    print(f"executed in {duration:.2f} seconds.")

if __name__ == "__main__":
    my_app()