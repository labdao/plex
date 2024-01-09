# import glob
import os
import time
import pandas as pd
from AF2_module import AF2Runner
import csv

import hydra
from hydra import compose, initialize
from hydra.core.hydra_config import HydraConfig
from omegaconf import DictConfig, OmegaConf
import glob
import json

def prodigy_run(csv_path):
    # print('csv path', csv_path)
    df = pd.read_csv(csv_path)
    # print('csv file data frame', df)
    for i, r in df.iterrows():

        file_path = r['pdb']
        if pd.notna(file_path):
            print('file path', file_path)
            try:
                subprocess.run(
                    ["prodigy", "-q", file_path], stdout=open("temp.txt", "w"), check=True
                )
                with open("temp.txt", "r") as f:
                    lines = f.readlines()
                    if lines:  # Check if lines is not empty
                        affinity = float(lines[0].split(" ")[-1].split("/")[0])
                        df.loc[i, "affinity"] = affinity
                    else:
                        # print(f"No output from prodigy for {r['path']}")
                        print(f"No output from prodigy for {file_path}")
                        # Handle the case where prodigy did not produce output
            except subprocess.CalledProcessError:
                # print(f"Prodigy command failed for {r['path']}")
                print(f"Prodigy command failed for {file_path}")

    # export results
    df.to_csv(f"{csv_path}", index=None)

# Please note the following:
# - You need to replace 'path_to_your_pdb_file.pdb' and 'path_to_your_pae_matrix.npy' with the actual paths to your PDB file and PAE matrix file.
# - The selection strings for the receptor and ligand (receptor_selection and ligand_selection) need to be adjusted according to your specific system. The example uses chain IDs, but you might need to use different selection criteria.
# - The scheme='closest-heavy' argument in md.compute_contacts specifies that only heavy atoms (non-hydrogen) are considered for contact calculations. Adjust this if necessary.
# - The PAE matrix is assumed to be indexed in a way that corresponds to the residue indices in the PDB file. If there's an offset or different indexing, you'll need to adjust the indexing in the line that computes interface_pae_values.
# 
# Before running the code, ensure that you have mdtraj and numpy installed in your Python environment. You can install them using pip if necessary:
# pip install mdtraj numpy
def compute_ipae(pdb_file):

    import mdtraj as md
    import numpy as np

    # Load the PDB file using mdtraj
    # pdb_file = 'path_to_your_pdb_file.pdb'  # Replace with your PDB file path
    traj = md.load(pdb_file)

    # Define the selection strings for the receptor and ligand
    receptor_selection = 'A'  # Example selection string for the receptor
    ligand_selection = 'B'  # Example selection string for the ligand

    # Select the receptor and ligand atoms
    receptor_atoms = traj.topology.select(receptor_selection)
    ligand_atoms = traj.topology.select(ligand_selection)

    # Compute contacts between receptor and ligand
    # The distance_cutoff is in nanometers; adjust as needed
    distance_cutoff = 0.35
    contacts = md.compute_contacts(traj, contacts=(receptor_atoms, ligand_atoms), scheme='closest-heavy', cutoff=distance_cutoff)

    # Get the contact pairs and distances
    contact_pairs = contacts[0]
    contact_distances = contacts[1].reshape(-1)  # Flatten the distances array

    # Filter contact pairs based on the distance cutoff
    contact_indices = contact_pairs[contact_distances < distance_cutoff]

    # Assuming you have the PAE matrix loaded as a NumPy array
    pae_matrix = np.load('path_to_your_pae_matrix.npy')  # Replace with your PAE matrix file path

    # Compute the PAE at the interface
    # Map the contact indices to the PAE matrix
    # Note: You may need to adjust the indices based on how your PAE matrix is indexed
    interface_pae_values = pae_matrix[contact_indices[:, 0], contact_indices[:, 1]]

    # Output the results
    print("Contact residues (receptor, ligand):", contact_indices)
    print("PAE values at the interface:", interface_pae_values)

    # Compute the median PAE for the interface contacts
    median_pae_interface = np.median(interface_pae_values)

    # Output the median PAE value
    print("Median PAE at the interface:", median_pae_interface)

    return median_pae_interface

# def update_complex_summary(t, row, sequence, sequence_pseudoLL, df, directory, json_pattern):
#     summary_file = os.path.join(directory, 'folding_with_target_summary.csv')

#     # loop over JSON files that match the given pattern for the current iteration
#     for json_file in glob.glob(os.path.join(directory, f"{json_pattern}_scores*.json")):
#         with open(json_file, 'r') as file:
#             data = json.load(file)

#         # Compute average plddt
#         avg_plddt = sum(data['plddt']) / len(data['plddt'])

#         # Get max_pae value
#         max_pae = data['max_pae']

#         # Find corresponding PDB file
#         pdb_file = None
#         rank_str = json_file.split('rank')[1].split('.')[0]
#         for pdb in glob.glob(os.path.join(directory, f"{json_pattern}_unrelaxed_rank{rank_str}*.pdb")):
#             pdb_file = pdb
#             break

#         # Prepare new row
#         new_row = {
#             't': t, # delete
#             'target_seq': row['target_seq'].iloc[0],
#             'binder_seq': row['binder_seq'].iloc[0],
#             'complex': sequence,
#             'sequence_pseudo_LL': sequence_pseudoLL, # delete
#             'mean plddt': avg_plddt,
#             'max pae': max_pae, # add i_pae # add rmsd
#             'json': os.path.abspath(json_file),
#             'pdb': os.path.abspath(pdb_file) if pdb_file else None
#         }

#         # Concatenate new row to DataFrame
#         df = pd.concat([df, pd.DataFrame([new_row])], ignore_index=True)

#     if not df.empty:
#         if os.path.exists(summary_file):
#             df.to_csv(summary_file, mode='a', header=False, index=False)
#         else:
#             df.to_csv(summary_file, index=False)

#     return df

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

        # Prepare new row
        new_row = {
            # 'sequence': 'extract from pdb',
            'mean plddt': avg_plddt,
            'max pae': max_pae,
            # 'i_pae': ,
            # 'rmsd': ,
            'absolute json path': os.path.abspath(json_file),
            'absolute pdb path': os.path.abspath(pdb_file) if pdb_file else None
        }

        # use absolute pdb path here and compute i_pae
        # compute_ipae(pdb_file)

        # Concatenate new row to DataFrame
        df = pd.concat([df, pd.DataFrame([new_row])], ignore_index=True)

    if not df.empty:
        if os.path.exists(summary_file):
            df.to_csv(summary_file, mode='a', header=False, index=False)
        else:
            df.to_csv(summary_file, index=False)

    return df

def parse_log_and_fasta_to_csv(fasta_path, log_path): # remove this
    # Read and store sequences from the FASTA file
    sequences = {}
    with open(fasta_path, 'r') as fasta_file:
        identifier = ''
        for line in fasta_file:
            if line.startswith('>'):
                identifier = line[1:].strip()
            else:
                sequences[identifier] = line.strip()

    # Initialize the columns for the CSV
    columns = ['Sequence', 'Query', 'structure_pdbs', 'recycle', 'pLDDT', 'pTM', 'tol', 'rank']
    
    # Parse the log file and write to CSV
    csv_path = os.path.join(os.path.dirname(log_path), 'summary.csv')
    with open(log_path, 'r') as log_file, open(csv_path, 'w', newline='') as csv_file:
        csv_writer = csv.DictWriter(csv_file, fieldnames=columns)
        csv_writer.writeheader()

        current_query = ''
        for line in log_file:
            if 'Query' in line:
                # Extracting the query value (e.g., "1/2" or "2/2")
                current_query = line.split('Query')[1].split(':')[0].strip()

            if 'pLDDT' in line and 'pTM' in line:
                parts = line.split()
                data = parts[2:]  # Skip the first two substrings (date and time)
                row_data = {col: '' for col in columns}  # Initialize row with empty values

                # Extract the sequence number (XXX in "XXX/YYY") and get the corresponding sequence
                sequence_number = current_query.split('/')[0]
                row_data['Sequence'] = sequences.get(sequence_number, '')
                row_data['Query'] = current_query

                # Handle the first substring without an equal sign
                row_data['structure_pdbs'] = data[0] if data else ''

                # Process remaining substrings with equal signs
                for item in data[1:]:
                    if '=' in item:
                        key, value = item.split('=')
                        if key in columns:
                            row_data[key] = value

                csv_writer.writerow(row_data)


    
    # print("running Prodigy")
    # prodigy_run(f"{self.outputs_directory}/folding_with_target_summary.csv")

    # need to write the prodigy scores to the csv; for this initialize the columns above

    print(f"CSV file written to {csv_path}")

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
    # parse_log_and_fasta_to_csv(os.path.abspath(fasta_file), os.path.join(outputs_directory, 'log.txt'))
    create_summary(outputs_directory, json_pattern='seq_1')

    print("Sequence to structure complete...")
    end_time = time.time()
    duration = end_time - start_time
    print(f"executed in {duration:.2f} seconds.")

if __name__ == "__main__":
    my_app()