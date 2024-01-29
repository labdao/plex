import pandas as pd
import os
import glob
import json
import numpy as np
import mdtraj as md
import subprocess
from Bio.PDB import PDBParser
from Bio.SeqUtils import seq1

def squeeze_seq(new_sequence):
    return ''.join(filter(lambda x: x != '-', new_sequence))

def generate_contig(action_mask, target, starting_target_residue=None, end_target_residue=None):
    if starting_target_residue is None:
        starting_target_residue = 1
    if end_target_residue is None:
        end_target_residue = len(target)
    
    # Initialize variables
    action_mask_contig = ''
    current_group = ''
    alphabet = 'LAGVSERTIDPKQNFYMHWC'
    position = 0  # Position within the action_mask
    
    # Iterate over the squeezed_action_mask to form groups
    for char in action_mask:
        if char in alphabet:
            if current_group == '' or current_group[-1] in alphabet:
                current_group += char  # Continue the current alphabet group
            elif current_group[-1]=='X':
                action_mask_contig += f'{len(current_group)}/'
                current_group = char
        elif char=='X':  # char is 'X'
            if current_group == '' or current_group[-1] == 'X':
                current_group += char  # Continue the current X group
            elif current_group[-1] in alphabet:
                action_mask_contig += f'B{position-len(current_group)+1}-{position}/'
                current_group = char
        elif char=='-':
            if current_group!='' and current_group[-1] in alphabet:
                action_mask_contig += f'B{position-len(current_group)+1}-{position}/'
                current_group = ''
            elif current_group!='' and current_group[-1]=='X':
                action_mask_contig += f'{len(current_group)}/'
                current_group = ''

        position += 1
    
    # Add the last group to the contig
    if current_group:
        if current_group[-1] == 'X':
            action_mask_contig += f'{len(current_group)}/'  # X group
        else:
            action_mask_contig += f'B{position-len(current_group)+1}-{position}/'  # Alphabet group
    
    # Remove the trailing '/' if it exists
    if action_mask_contig.endswith('/'):
        action_mask_contig = action_mask_contig[:-1]
    
    # Construct the final contig string
    contig = f'A{starting_target_residue}-{end_target_residue}/0 {action_mask_contig}'
    
    return contig

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
        print("Invalid file path")
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
    print("median PAE at the interface:", median_pae_interface)

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
        print('pdb file', pdb_file)
        print('rank str', rank_str)

        if pdb_file:
            pae_matrix = np.array(data['pae'])
            i_pae = compute_ipae(os.path.abspath(pdb_file), pae_matrix)

            sequence = extract_sequence_from_pdb(os.path.abspath(pdb_file))

            pdb_file_path = os.path.abspath(pdb_file)

            # Add new columns to the DataFrame if they don't exist
            if 'sequence' not in df:
                df['sequence'] = None
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
            df.at[0, 'sequence'] = sequence
            df.at[0, 'mean plddt'] = avg_plddt
            df.at[0, 'max pae'] = max_pae
            df.at[0, 'i_pae'] = i_pae
            df.at[0, 'absolute json path'] = os.path.abspath(json_file)
            df.at[0, 'absolute pdb path'] = pdb_file_path

    return df

def concatenate_to_df(t, df, df_main):
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

def reinsert_deletions(modified_seq, action_mask):
    if len(modified_seq) != len(action_mask.replace('-', '')): # check if lengths match when '-' is removed from action mask
        raise ValueError("Length of modified_seq does not match the length of action_mask without '-' characters.")

    seq_with_deletions = ''
    modified_seq_index = 0

    # Iterate over the action_mask and construct seq_with_deletions
    for char in action_mask:
        if char == '-':
            seq_with_deletions += '-'
        else:
            seq_with_deletions += modified_seq[modified_seq_index]
            modified_seq_index += 1

    return seq_with_deletions

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

    # Ensure permissibility_seed is long enough and initialize with ':'
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
                elif type_char == '+':
                    permissibility_seed[i-1] = '+'

    # Join the list into a string and return
    return ''.join(permissibility_seed)
