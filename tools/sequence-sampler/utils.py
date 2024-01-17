import pandas as pd
import os
import glob
import json
import numpy as np
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
            # i_pae = compute_ipae(os.path.abspath(pdb_file), pae_matrix)

            sequence = extract_sequence_from_pdb(os.path.abspath(pdb_file))

            pdb_file_path = os.path.abspath(pdb_file)
            # affinity = compute_affinity(pdb_file_path)

            # Add new columns to the DataFrame if they don't exist
            if 'sequence' not in df:
                df['sequence'] = None
            if 'mean plddt' not in df:
                df['mean plddt'] = None
            if 'max pae' not in df:
                df['max pae'] = None
            # if 'i_pae' not in df:
            #     df['i_pae'] = None
            # if 'affinity' not in df:
            #     df['affinity'] = None
            if 'absolute json path' not in df:
                df['absolute json path'] = None
            if 'absolute pdb path' not in df:
                df['absolute pdb path'] = None

            # Update the DataFrame with new values
            df.at[0, 'sequence'] = sequence
            df.at[0, 'mean plddt'] = avg_plddt
            df.at[0, 'max pae'] = max_pae
            # df.at[0, 'i_pae'] = i_pae
            # df.at[0, 'affinity'] = affinity
            df.at[0, 'absolute json path'] = os.path.abspath(json_file)
            df.at[0, 'absolute pdb path'] = pdb_file_path

    return df

def read_second_line_of_fasta(file_path):
    with open(file_path, 'r') as file:
        lines = file.readlines()
        if len(lines) >= 2:
            return lines[1].strip()
    return None

def reinsert_deletions(modified_seq, action_mask):
    # Remove '-' from action_mask and check if lengths match
    print('reinsert deletions - modified seq', modified_seq)
    print('reinsert deletions - actions_mask', action_mask)
    if len(modified_seq) != len(action_mask.replace('-', '')):
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

# def read_last_line_of_fasta(file_path):
#     with open(file_path, 'r') as file:
#         lines = file.readlines()
#         # Skip to the last non-empty line that starts with '>'
#         for last_line in reversed(lines):
#             if last_line.strip() and not last_line.startswith('>'):
#                 return last_line.strip()
#     return None