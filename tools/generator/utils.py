import pandas as pd
import os
import glob
import json
import numpy as np
import mdtraj as md
import subprocess
from Bio.PDB import PDBParser
from Bio.SeqUtils import seq1
import random
import sys
import logging
import boto3

from omegaconf import DictConfig, OmegaConf

def compute_log_likelihood(sequence, LLmatrix): # TD: move into the scorer module or even to sequence-transformer

    sequence = squeeze_seq(sequence)
    if len(sequence) != LLmatrix.shape[1]:
        raise ValueError("Length of sequence must match the number of columns in LLmatrix.")
    
    amino_acid_code = ''.join('LAGVSERTIDPKQNFYMHWC')

    # Initialize total log likelihood
    total_log_likelihood = 0

    # Compute the total log likelihood of sequence
    for i, aa in enumerate(sequence):
        # Find the row index for this amino acid
        row_index = amino_acid_code.index(aa)
        
        # Add the log likelihood from the corresponding cell in LLmatrix
        total_log_likelihood += LLmatrix[row_index, i]

    return total_log_likelihood

def squeeze_seq(new_sequence):
    return ''.join(filter(lambda x: x != '-', new_sequence))

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
                    # extract affinity value from the output
                    affinity = float(lines[0].split(" ")[-1].split("/")[0])
                    return affinity
                else:
                    logging.info(f"No output from prodigy for {file_path}")
                    return None  # no output from prodigy
        except subprocess.CalledProcessError:
            logging.info(f"Warning: Prodigy command failed for {file_path}. This is not an error per se and most likely due to the binder not being closely positioned against the target.")
            return None  # Prodigy command failed
    else:
        logging.info(f"Invalid file path")
        return None  # Invalid file path provided

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
    logging.info(f"median PAE at the interface {median_pae_interface}")

    return median_pae_interface

def write_af2_update(df, directory, json_pattern):
    # Loop over JSON files that match the given pattern for the current iteration
    for json_file in glob.glob(os.path.join(directory, f"{json_pattern}_scores_rank_001*.json")):
        with open(json_file, 'r') as file:
            data = json.load(file)

        # compute metric values
        avg_plddt = sum(data['plddt']) / len(data['plddt'])
        max_pae = data['max_pae']

        # Find corresponding PDB file
        pdb_file = None
        rank_str = json_file.split('rank')[1].split('.')[0]
        for pdb in glob.glob(os.path.join(directory, f"{json_pattern}_unrelaxed_rank_001*.pdb")):
            pdb_file = pdb
            break

        if pdb_file:
            pae_matrix = np.array(data['pae'])
            i_pae = compute_ipae(os.path.abspath(pdb_file), pae_matrix)

            sequence = extract_sequence_from_pdb(os.path.abspath(pdb_file))

            pdb_file_path = os.path.abspath(pdb_file)

            # Add new columns to the DataFrame if they don't exist
            if 'sequence of complex' not in df:
                df['sequence of complex'] = None
            if 'mean plddt' not in df:
                df['mean plddt'] = None
            if 'max pae' not in df:
                df['max pae'] = None
            if 'i_pae' not in df:
                df['i_pae'] = None
            if 'absolute json path' not in df:
                df['absolute json path'] = None
            if 'absolute pdb path' not in df:
                df['absolute pdb path'] = None

            # Update the DataFrame with new values
            df.at[0, 'sequence of complex'] = sequence
            df.at[0, 'mean plddt'] = avg_plddt
            df.at[0, 'max pae'] = max_pae
            df.at[0, 'i_pae'] = i_pae
            df.at[0, 'absolute json path'] = os.path.abspath(json_file)
            df.at[0, 'absolute pdb path'] = pdb_file_path

    return df

def determine_acceptance(threshold):
    # Generate a boolean value based on the threshold probability
    accept_flag = random.choices([True, False], weights=[threshold, 1-threshold], k=1)[0]

    return accept_flag

def concatenate_to_df(t, df, df_main):
    if t == 0:
        # Find the row in df_main with 't' == 0 and acceptance_flag == True
        target_row = df_main[(df_main['t'] == 0) & (df_main['acceptance_flag'] == True)].index
        if not target_row.empty:
            target_row_index = target_row[0]
            # Ensure all columns in df are in df_main, if not, add them with the values from df
            for col in df.columns:
                if col not in df_main.columns:
                    # Add the column to df_main and fill all previous rows with pd.NA
                    df_main[col] = pd.NA
                # Write the values into the row of df_main which has 't'-column value 0 and acceptance_flag value True
                df_main.at[target_row_index, col] = df.at[df.index[0], col]
    else:
        # Ensure all columns in df are in df_main, if not, add them with the values from df
        for col in df.columns:
            if col not in df_main.columns:
                # Add the column to df_main and fill all previous rows with pd.NA
                df_main[col] = pd.NA
                # Fill the last row with the value from df
                df_main.at[df_main.index[-1], col] = df.at[df.index[0], col]
            else:
                # If the column exists, just append the new value
                df_main.at[df_main.index[-1], col] = df.at[df.index[0], col]

    return df_main

def read_second_line_of_fasta(file_path):
    with open(file_path, 'r') as file:
        lines = file.readlines()
        if len(lines) >= 2:
            return lines[1].strip()
    return None

def slash_to_convexity_notation(sequence, slash_contig):

    # Find the maximum index required
    max_index = 0
    parts = slash_contig.split('/')
    for part in parts:
        if ':' in part:
            _, end = map(int, part[1:].split(':'))
            max_index = max(max_index, end)
        elif part:
            max_index = max(max_index, int(part[1:]))

    # Ensure permissibility_seed is long enough and initialize with '-'
    permissibility_seed = ['-'] * max(max_index, len(sequence))

    # Process each part of the slash_contig
    for part in parts:
        if part:
            type_char = part[0]
            if ':' in part:
                start, end = map(int, part[1:].split(':'))
            else:
                start = end = int(part[1:])

            for i in range(start, end + 1):
                if type_char == 'B':
                    permissibility_seed[i-1] = sequence[i-1] if i-1 < len(sequence) else '-'
                elif type_char == 'x':
                    permissibility_seed[i-1] = 'X'
                elif type_char == '*':
                    permissibility_seed[i-1] = '*'

    # Join the list into a string and return
    return ''.join(permissibility_seed)

def user_input_parsing(cfg: DictConfig, user_inputs: dict) -> DictConfig:

    OmegaConf.update(cfg, "params.basic_settings.number_of_binders", user_inputs["number_of_binders"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.sequence_input", user_inputs["binder_protein_sequence"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.target_seq", user_inputs["target_protein_sequence"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.init_permissibility_vec", user_inputs["init_permissibility_vec"], merge=False)

    hotspots = user_inputs["hotspots"]
    user_inputs["hotspots"] = '[' + hotspots.replace(' ', '') + ']'
    OmegaConf.update(cfg, "params.RFdiffusion_settings.hotspots", user_inputs["hotspots"], merge=False)

    OmegaConf.update(cfg, "params.basic_settings.high_fidelity", user_inputs["high_fidelity"], merge=False)

    # OmegaConf.update(cfg, "params.basic_settings.target_pdb", user_inputs["target_pdb"], merge=False)
    # OmegaConf.update(cfg, "params.basic_settings.target_chain", user_inputs["target_chain"], merge=False)
    
    return cfg

def replace_invalid_characters(seed, alphabet):
    return ''.join(['X' if c not in alphabet and c not in ['*', 'x'] else c for c in seed])

def check_gpu_availability():
    logging.basicConfig(level=logging.INFO)
    
    # Check for nvidia-smi
    try:
        subprocess.check_output(['nvidia-smi'])
    except subprocess.CalledProcessError as e:
        logging.info("nvidia-smi GPU check failed. Ending job.")
        sys.exit(1)
    except FileNotFoundError:
        logging.info("nvidia-smi is not installed.")
        sys.exit(1)
    
    # Check for PyTorch GPU
    try:
        import torch
        if not torch.cuda.is_available():
            raise RuntimeError("PyTorch cannot find a GPU. Ending job.")
    except ImportError:
        logging.info("PyTorch is not installed.")
        sys.exit(1)
    except RuntimeError as e:
        logging.info(str(e))
        sys.exit(1)

def check_string_against_alphabet(input_string, alphabet='LAGVSERTIDPKQNFYMHWC'):
    for char in input_string:
        if char not in alphabet:
            logging.info(f"Warning: not enough input information given to predict an initial binder sequence. Please try to provide the amino acid identity for more of the residues.")
            sys.exit(1)

def expand_and_clean_sequence(input_sequence, alphabet):
    import re

    # Remove all spaces
    sequence = input_sequence.replace(" ", "").upper()

    # Convert to uppercase and replace invalid characters with 'X'
    valid_chars = set(alphabet.upper() + 'X*') | set('0123456789')
    sequence = ''.join(char if char.upper() in valid_chars else 'X' for char in sequence)

    # Find all occurrences of character followed by '*' and a number
    pattern = re.compile(r'([A-Z])(\*\d+)')
    
    # Function to replace the pattern with the character repeated
    def replacer(match):
        char = match.group(1)
        count = int(match.group(2)[1:])  # Extract the number and convert to int
        return char * count

    # Replace all occurrences using the replacer function
    expanded_sequence = pattern.sub(replacer, sequence)

    return expanded_sequence.upper()

def upload_to_s3(file_name, bucket_name, object_name=None):
    if object_name is None:
        object_name = file_name

    s3_client = boto3.client('s3')
    try:
        s3_client.upload_file(file_name, bucket_name, object_name)
        print(f"Successfully uploaded {file_name} to {bucket_name}/{object_name}")
    except Exception as e:
        print(f"Failed to upload {file_name} to {bucket_name}/{object_name}: {e}")
        raise e

def create_and_upload_checkpoints(df_results, result_csv_path, index):
    checkpoint_csv_path = os.path.dirname(result_csv_path)
    # Get the last row of the DataFrame
    row = df_results.iloc[-1]
    pdb_path = row['absolute pdb path']
    pdb_file_name = os.path.basename(pdb_path)

    new_df = pd.DataFrame({
        'cycle': row['t'],
        'proposal': row['sample_number'],
        'plddt': row['mean plddt'],
        'i_pae': row['i_pae'],
        'dim1': row['max pae'],
        'dim2': row['affinity'],
        'pdbFileName': pdb_file_name
    }, index=[0])

    new_df.fillna(0, inplace=True)

    event_csv_filename = f"checkpoint_{index}_summary.csv"
    event_csv_filepath = f"{checkpoint_csv_path}/{event_csv_filename}"
    new_df.to_csv(event_csv_filepath, index=False)

    bucket_name = "app-checkpoint-bucket"
    job_uuid = os.getenv("JOB_UUID")
    object_name = f"checkpoints/{job_uuid}/checkpoint_{index}"
    upload_to_s3(event_csv_filepath, bucket_name, f"{object_name}/{event_csv_filename}")
    upload_to_s3(pdb_path, bucket_name, f"{object_name}/{pdb_file_name}")
    os.remove(event_csv_filepath)
    print(f"Checkpoint {index} event CSV and PDB uploaded.")

    return