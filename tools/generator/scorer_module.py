import os
import glob
import pandas as pd
from Bio.PDB import PDBParser, Superimposer, PDBIO
import sequence_transformer
from AF2_module import AF2Runner
from omegafold_module import Omegafold
from utils import squeeze_seq
from utils import write_af2_update
from utils import compute_affinity
from utils import concatenate_to_df
from utils import compute_log_likelihood
import logging
from utils import check_gpu_availability

class Scorer:

    def __init__(self, cfg, outputs_directory):

        self.cfg = cfg
        self.outputs_directory = outputs_directory

    def run(self, t, sequence, df):

        scorer_list = self.cfg.params.basic_settings.scorers.split(',')

        logging.info(f"Running scoring job...")
        df_score = pd.DataFrame() # initialize data frame
        for scorer in scorer_list:

            scorer_directory = os.path.join(self.outputs_directory, scorer)
            if not os.path.exists(scorer_directory):
                os.makedirs(scorer_directory, exist_ok=True)

            logging.info(f"Running {scorer}")
            if scorer=='ESM2' or scorer=='esm2':
                runner = sequence_transformer.ESM2Runner() # initialize ESM2Runner with the default model
                LLmatrix_sequence = runner.token_masked_marginal_log_likelihood_matrix(squeeze_seq(sequence))

                LL_mod = compute_log_likelihood(sequence, LLmatrix_sequence) # TD: normalization by sequence length?

                if 'pseudolikelihood' not in df_score.columns:
                    df_score['pseudolikelihood'] = None  # Initialize the column with None

                # Set the value of 'pseudolikelihood' for the first row
                df_score.at[0, 'pseudolikelihood'] = LL_mod

            elif scorer=='Colabfold' or scorer=='colabfold':
                target_binder_sequence = f"{self.cfg.params.basic_settings.target_seq}:{squeeze_seq(sequence)}"
                
                # include a function that combines binder and target sequence
                input_dir = os.path.join(self.cfg.inputs.directory, 'current_sequences')
                if os.path.exists(input_dir):
                    # If the folder already exists, empty the folder of all files
                    for file_name in os.listdir(input_dir):
                        file_path = os.path.join(input_dir, file_name)
                        if os.path.isfile(file_path):
                            os.remove(file_path)
                else:
                    os.makedirs(input_dir, exist_ok=True)

                file_path = os.path.join(input_dir, f"design_cycle_{t}.fasta")
                with open(file_path, 'w') as file:
                    file.write(f">design_cycle_{t}\n{target_binder_sequence}\n")

                seq_input_dir = os.path.abspath(input_dir)

                check_gpu_availability()

                af2_runner = AF2Runner(seq_input_dir, scorer_directory)
                af2_runner.run()

                # append output as new columns of data frame
                df_score = write_af2_update(df_score, scorer_directory, json_pattern=f"design_cycle_{t}")

            elif scorer=='omegafold_with_alignment':

                def parse_b_factors(pdb_file_path):
                    b_factors = []
                    with open(pdb_file_path, 'r') as file:
                        for line in file:
                            if line.startswith("ATOM"):
                                b_factor = float(line[60:66].strip())
                                b_factors.append(b_factor)
                    return sum(b_factors) / len(b_factors) if b_factors else 0
                
                def fold_initial_complex():
                
                    target_binder_sequence = f"{self.cfg.params.basic_settings.target_seq}:{squeeze_seq(sequence)}"
                    
                    # include a function that combines binder and target sequence
                    input_dir = os.path.join(self.cfg.inputs.directory, 'current_sequences')
                    if os.path.exists(input_dir):
                        # If the folder already exists, empty the folder of all files
                        for file_name in os.listdir(input_dir):
                            file_path = os.path.join(input_dir, file_name)
                            if os.path.isfile(file_path):
                                os.remove(file_path)
                    else:
                        os.makedirs(input_dir, exist_ok=True)

                    file_path = os.path.join(input_dir, f"design_cycle_{0}.fasta")
                    with open(file_path, 'w') as file:
                        file.write(f">design_cycle_{t}\n{target_binder_sequence}\n")

                    seq_input_dir = os.path.abspath(input_dir)

                    check_gpu_availability()

                    af2_runner = AF2Runner(seq_input_dir, scorer_directory)
                    af2_runner.run()

                    # append output as new columns of data frame
                    df_score = pd.DataFrame() # initialize data frame
                    df_score = write_af2_update(df_score, scorer_directory, json_pattern=f"design_cycle_{0}")

                    return df_score

                if t==0: # run colabfold
                    logging.info("running colabfold")
                    df_score = fold_initial_complex()

                elif t>=1: # run omegafold and create complex by alignment
                    logging.info("running omegafold and alignment")

                    binder_sequence = f"{squeeze_seq(sequence)}"
                    
                    # include a function that combines binder and target sequence
                    input_dir = os.path.join(self.cfg.inputs.directory, 'current_sequences')
                    if os.path.exists(input_dir):
                        # If the folder already exists, empty the folder of all files
                        for file_name in os.listdir(input_dir):
                            file_path = os.path.join(input_dir, file_name)
                            if os.path.isfile(file_path):
                                os.remove(file_path)
                    else:
                        os.makedirs(input_dir, exist_ok=True)

                    file_path = os.path.join(input_dir, f"design_cycle_{t}.fasta")
                    with open(file_path, 'w') as file:
                        file.write(f">design_cycle_{t}\n{binder_sequence}\n")

                    check_gpu_availability()

                    file_path_abs = os.path.abspath(file_path)

                    omegafold_runner = Omegafold(file_path_abs, scorer_directory)
                    omegafold_runner.run()

                    # extract b-scores and compute average
                    pdb_file_path = os.path.join(scorer_directory, f"design_cycle_{t}.pdb")
                    average_b_score = parse_b_factors(pdb_file_path)
                    logging.info(f"The average B-score is: {average_b_score}")

                    if 'average_b-score' not in df_score.columns:
                        df_score['average_b-score'] = None  # Initialize the column with None

                    df_score.at[0, 'average_b-score'] = average_b_score

                    search_pattern = os.path.join(scorer_directory, 'design_cycle_0_unrelaxed_rank_001*.pdb')
                    pdb_files = glob.glob(search_pattern)
                    original_complex_file_path = os.path.abspath(pdb_files[0]) if pdb_files else None
                
                    pdb_parser = PDBParser()
                    structure_complex = pdb_parser.get_structure('original_complex', original_complex_file_path)
                    structure_modified = pdb_parser.get_structure('modified_B', pdb_file_path)
                    
                    original_chain_B = structure_complex[0]['B']
                    modified_chain_A = structure_modified[0]['A']

                    fixed_atoms = [atom for atom in original_chain_B.get_atoms() if atom.get_id() == 'CA'] # Prepare for alignment by selecting backbone atoms (e.g., alpha carbons)
                    moving_atoms = [atom for atom in modified_chain_A.get_atoms() if atom.get_id() == 'CA']

                    if len(fixed_atoms) == 0 or len(moving_atoms) == 0 or len(fixed_atoms) != len(moving_atoms): # check that the lists are not empty and have the same length; otherwise, alignment is not possible
                        raise ValueError("Alignment cannot be performed due to unequal numbers of backbone atoms or no backbone atoms found.")

                    super_imposer = Superimposer() # Perform the alignment
                    super_imposer.set_atoms(fixed_atoms, moving_atoms)
                    super_imposer.apply(modified_chain_A.get_atoms())

                    new_structure = structure_complex[0].copy()
                    new_structure.detach_child('B')
                    modified_chain_A.id = 'B'

                    chain_to_renumber = modified_chain_A # renumber residues to start counting from 1 instead of 0
                    temp_residue_number = 10000
                    for residue in chain_to_renumber.get_residues():
                        residue.id = (' ', temp_residue_number, ' ')
                        temp_residue_number += 1

                    new_residue_number = 1
                    for residue in chain_to_renumber.get_residues():
                        residue.id = (' ', new_residue_number, ' ')
                        new_residue_number += 1

                    new_structure.add(modified_chain_A)

                    io = PDBIO()
                    io.set_structure(new_structure)
                    path_to_complex = os.path.join(scorer_directory, f"design_cycle_{t}.pdb")
                    io.save(path_to_complex)

                    if 'absolute pdb path' not in df_score.columns:
                        df_score['absolute pdb path'] = None  # Initialize the column with None

                    # Set the value of 'pseudolikelihood' for the first row
                    df_score.at[0, 'absolute pdb path'] = path_to_complex

                    print(f"New PDB with chain A and aligned modified chain B created as {path_to_complex}")

            elif scorer=='omegafold_initial_fold':

                def parse_b_factors(pdb_file_path):
                    b_factors = []
                    with open(pdb_file_path, 'r') as file:
                        for line in file:
                            if line.startswith("ATOM"):
                                b_factor = float(line[60:66].strip())
                                b_factors.append(b_factor)
                    return sum(b_factors) / len(b_factors) if b_factors else 0
    
                def renumber_chain_a_residues(pdb_file_path):
                    # Read the existing PDB file
                    with open(pdb_file_path, 'r') as file:
                        lines = file.readlines()

                    # Write back the lines with updated residue numbers for chain A
                    with open(pdb_file_path, 'w') as file:
                        for line in lines:
                            if line.startswith("ATOM") or line.startswith("TER"):
                                chain_id = line[21]
                                if chain_id == 'A':
                                    residue_number = int(line[22:26].strip()) + 1  # Increment the residue number by 1
                                    new_line = line[:22] + f"{residue_number:4d}" + line[26:]
                                    file.write(new_line)
                                else:
                                    file.write(line)
                            else:
                                file.write(line)

                def AF2_fold_complex():
                
                    target_binder_sequence = f"{self.cfg.params.basic_settings.target_seq}:{squeeze_seq(sequence)}"
                    
                    # include a function that combines binder and target sequence
                    input_dir = os.path.join(self.cfg.inputs.directory, 'current_sequences')
                    if os.path.exists(input_dir):
                        for file_name in os.listdir(input_dir):
                            file_path = os.path.join(input_dir, file_name)
                            if os.path.isfile(file_path):
                                os.remove(file_path)
                    else:
                        os.makedirs(input_dir, exist_ok=True)

                    file_path = os.path.join(input_dir, f"design_cycle_{t}.fasta")
                    with open(file_path, 'w') as file:
                        file.write(f">design_cycle_{t}\n{target_binder_sequence}\n")

                    seq_input_dir = os.path.abspath(input_dir)

                    check_gpu_availability()

                    af2_runner = AF2Runner(seq_input_dir, scorer_directory)
                    af2_runner.run()

                    # append output as new columns of data frame
                    df_score = pd.DataFrame() # initialize data frame
                    df_score = write_af2_update(df_score, scorer_directory, json_pattern=f"design_cycle_{t}")


                    return df_score

                if t==0: # run colabfold
                    logging.info("running omegafold")

                    # target_binder_sequence = f"{self.cfg.params.basic_settings.target_seq}:{squeeze_seq(sequence)}"
                    binder_sequence = f"{self.cfg.params.basic_settings.target_seq}"

                    # include a function that combines binder and target sequence
                    input_dir = os.path.join(self.cfg.inputs.directory, 'current_sequences')
                    if os.path.exists(input_dir):
                        # If the folder already exists, empty the folder of all files
                        for file_name in os.listdir(input_dir):
                            file_path = os.path.join(input_dir, file_name)
                            if os.path.isfile(file_path):
                                os.remove(file_path)
                    else:
                        os.makedirs(input_dir, exist_ok=True)

                    file_path = os.path.join(input_dir, f"design_cycle_{t}.fasta")
                    with open(file_path, 'w') as file:
                        file.write(f">design_cycle_{t}\n{binder_sequence}\n")

                    check_gpu_availability()

                    file_path_abs = os.path.abspath(file_path)

                    omegafold_runner = Omegafold(file_path_abs, scorer_directory)
                    omegafold_runner.run()

                    # extract b-scores and compute average
                    pdb_file_path = os.path.join(scorer_directory, f"design_cycle_{t}.pdb")
                    average_b_score = parse_b_factors(pdb_file_path)
                    logging.info(f"The average B-score is: {average_b_score}")

                    if 'average_b-score' not in df_score.columns:
                        df_score['average_b-score'] = None  # Initialize the column with None

                    df_score.at[0, 'average_b-score'] = average_b_score

                    renumber_chain_a_residues(pdb_file_path)
                    path_to_complex = pdb_file_path

                    if 'absolute pdb path' not in df_score.columns:
                        df_score['absolute pdb path'] = None  # Initialize the column with None

                    df_score.at[0, 'absolute pdb path'] = path_to_complex

                    print(f"New PDB with chain A and fully masked chain B created as {path_to_complex}")
                                
                elif t==1: # run colabfold
                    logging.info("running colabfold")
                    df_score = AF2_fold_complex()

                elif t>=2: # run omegafold and create complex by alignment
                    logging.info("running omegafold and alignment")

                    binder_sequence = f"{squeeze_seq(sequence)}"
                    
                    input_dir = os.path.join(self.cfg.inputs.directory, 'current_sequences')
                    if os.path.exists(input_dir):
                        for file_name in os.listdir(input_dir):
                            file_path = os.path.join(input_dir, file_name)
                            if os.path.isfile(file_path):
                                os.remove(file_path)
                    else:
                        os.makedirs(input_dir, exist_ok=True)

                    file_path = os.path.join(input_dir, f"design_cycle_{t}.fasta")
                    with open(file_path, 'w') as file:
                        file.write(f">design_cycle_{t}\n{binder_sequence}\n")

                    check_gpu_availability()

                    file_path_abs = os.path.abspath(file_path)

                    omegafold_runner = Omegafold(file_path_abs, scorer_directory)
                    omegafold_runner.run()

                    # extract b-scores and compute average
                    pdb_file_path = os.path.join(scorer_directory, f"design_cycle_{t}.pdb")
                    average_b_score = parse_b_factors(pdb_file_path)
                    logging.info(f"The average B-score is: {average_b_score}")

                    if 'average_b-score' not in df_score.columns:
                        df_score['average_b-score'] = None  # Initialize the column with None

                    df_score.at[0, 'average_b-score'] = average_b_score

                    search_pattern = os.path.join(scorer_directory, 'design_cycle_1_unrelaxed_rank_001*.pdb') # important to take cycle_1 here and not cycle_0
                    pdb_files = glob.glob(search_pattern)
                    original_complex_file_path = os.path.abspath(pdb_files[0]) if pdb_files else None
                
                    pdb_parser = PDBParser()
                    structure_complex = pdb_parser.get_structure('original_complex', original_complex_file_path)
                    structure_modified = pdb_parser.get_structure('modified_B', pdb_file_path)
                    
                    original_chain_B = structure_complex[0]['B']
                    modified_chain_A = structure_modified[0]['A']

                    fixed_atoms = [atom for atom in original_chain_B.get_atoms() if atom.get_id() == 'CA'] # Prepare alignment by selecting backbone (alpha carbon) atoms
                    moving_atoms = [atom for atom in modified_chain_A.get_atoms() if atom.get_id() == 'CA']

                    if len(fixed_atoms) == 0 or len(moving_atoms) == 0 or len(fixed_atoms) != len(moving_atoms): # check that the lists are not empty and have the same length; otherwise, alignment is not possible
                        raise ValueError("Alignment cannot be performed due to unequal numbers of backbone atoms or no backbone atoms found.")

                    super_imposer = Superimposer() # Perform the alignment
                    super_imposer.set_atoms(fixed_atoms, moving_atoms)
                    super_imposer.apply(modified_chain_A.get_atoms())

                    new_structure = structure_complex[0].copy()
                    new_structure.detach_child('B')
                    modified_chain_A.id = 'B'

                    chain_to_renumber = modified_chain_A # renumber residues to start counting from 1 instead of 0
                    temp_residue_number = 10000
                    for residue in chain_to_renumber.get_residues():
                        residue.id = (' ', temp_residue_number, ' ')
                        temp_residue_number += 1

                    new_residue_number = 1
                    for residue in chain_to_renumber.get_residues():
                        residue.id = (' ', new_residue_number, ' ')
                        new_residue_number += 1

                    new_structure.add(modified_chain_A)

                    io = PDBIO()
                    io.set_structure(new_structure)
                    path_to_complex = os.path.join(scorer_directory, f"design_cycle_{t}.pdb")
                    io.save(path_to_complex)

                    if 'absolute pdb path' not in df_score.columns:
                        df_score['absolute pdb path'] = None  # Initialize the column with None

                    # Set the value of 'pseudolikelihood' for the first row
                    df_score.at[0, 'absolute pdb path'] = path_to_complex

                    print(f"New PDB with chain A and aligned modified chain B created as {path_to_complex}")
            
            elif scorer=='Prodigy' or scorer=='prodigy': # not implemented yet

                pdb_file_path = df_score['absolute pdb path'].iloc[0]
                affinity = compute_affinity(pdb_file_path)
                if 'affinity' not in df_score.columns:
                    df_score['affinity'] = None  # Initialize the column with None

                # Assuming you have a single row, set the value of 'affinity' for the first row
                df_score.at[0, 'affinity'] = affinity

                if affinity is not None:
                    print(f"Affinity for complex {pdb_file_path} is {affinity}")

            elif scorer=='Hamming' or 'hamming':

                # Function to compute Hamming distance
                def compute_hamming_distance(seq1, seq2):
                    return sum(c1 != c2 for c1, c2 in zip(seq1, seq2))

                # Filter rows where 't' column is 0
                filtered_df = df[df['t'] == 0]

                # Compute Hamming distances for filtered rows
                hamming_distances = filtered_df['modified_seq'].apply(lambda x: compute_hamming_distance(squeeze_seq(sequence), squeeze_seq(x)))

                # Calculate the mean of the Hamming distances
                mean_hamming_distance = hamming_distances.mean()

                # Add the mean Hamming distance to df_score
                if 'mean_hamming_distance_to_init_seqs' not in df_score.columns:
                    df_score['mean_hamming_distance_to_init_seqs'] = None  # Initialize the column with None

                # Set the value of 'hamming_distance' for the first row
                df_score.at[0, 'mean_hamming_distance_to_init_seqs'] = mean_hamming_distance

                if mean_hamming_distance is not None:
                    logging.info(f"Mean Hamming distance for selected sequences is {mean_hamming_distance}")
            
            df_score.to_csv(f"{scorer_directory}/output.csv", index=False) # TD: treat the case when no scorer is given. currently, even when there is no scorer, something seems to be written
        
        logging.info(f"Scoring job complete. Results are in {self.outputs_directory}")

        df = concatenate_to_df(t, df_score, df)

        return df
