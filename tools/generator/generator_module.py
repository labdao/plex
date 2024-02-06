import subprocess
import os
import pandas as pd
import numpy as np
from utils import generate_contig
from utils import read_second_line_of_fasta
from utils import reinsert_deletions
import logging

class StateGenerator:
    def __init__(self, evo_cycle, sequence, action_mask, cfg, outputs_directory, df, permissibility_vector):
        self.evo_cycle = evo_cycle
        self.sequence = sequence
        self.action_mask = action_mask
        self.cfg = cfg
        self.target = cfg.params.basic_settings.target_seq
        self.outputs_directory = outputs_directory
        self.df = df
        self.generators = cfg.params.basic_settings.generators.split(',')
        self.alphabet = cfg.params.basic_settings.alphabet
        self.max_levenshtein_step_size = cfg.params.basic_settings.max_levenshtein_step_size
        self.permissibility_vector = permissibility_vector

    def run_generation(self):

        print('\n')
        logging.info(f"Running generating job...")
        df_generate = pd.DataFrame() # initialize data frame

        for generator in self.generators:

            generator_directory = os.path.join(self.outputs_directory, generator)
            if not os.path.exists(generator_directory):
                os.makedirs(generator_directory, exist_ok=True)

            if generator=='RFdiffusion+ProteinMPNN':
                logging.info(f"Running {generator}")

                if 'X' in self.action_mask: # check if there is any diffusion to be done

                    logging.info(f"diffusing...")

                    logging.info(f"action mask, {self.action_mask}")
                    contig = generate_contig(self.action_mask, self.target, starting_target_residue=None, end_target_residue=None)
                    logging.info(f"diffusion contig, {contig}")

                    # Set up the environment for the subprocess - required so that RFdiffussion can find its proper packages
                    env = os.environ.copy()
                    env['PYTHONPATH'] = "/app/RFdiffusion:" + env.get('PYTHONPATH', '')

                    command = [
                        'python', 'RFdiffusion/scripts/run_inference.py',
                        f'inference.output_prefix={os.path.join(generator_directory, f"evocycle_{self.evo_cycle}_motifscaffolding")}',
                        'inference.model_directory_path=RFdiffusion/models',
                        f'inference.input_pdb={self.df["absolute pdb path"].iloc[0]}',
                        'inference.num_designs=3',
                        f'contigmap.contigs={[contig]}'
                    ]

                    result = subprocess.run(command, capture_output=True, text=True, env=env)

                    # Check if the command was successful
                    if result.returncode == 0:
                        logging.info(f"#Inference script ran successfully")
                        logging.info(result.stdout)
                    else:
                        logging.info(f"#Error running inference script")
                        logging.info(result.stderr)
                    
                    logging.info(f"Running MPNN")

                    # Activate the conda environment 'mlfold'
                    subprocess.run(['conda', 'activate', 'mlfold'], shell=True) # TD: I think this can be removed - check this.

                    # Define the paths and parameters
                    path_to_PDB = os.path.join(generator_directory, f"evocycle_{self.evo_cycle}_motifscaffolding_0.pdb")
                    output_dir = generator_directory
                    chains_to_design = 'A'

                    # Create the output directory if it doesn't exist
                    os.makedirs(output_dir, exist_ok=True)

                    logging.info(f"pdb path, {path_to_PDB}")

                    # Define the command and arguments
                    command = [
                        'python', 'ProteinMPNN/protein_mpnn_run.py',
                        '--pdb_path', path_to_PDB,
                        '--pdb_path_chains', chains_to_design,
                        '--out_folder', output_dir,
                        '--num_seq_per_target', '8',
                        '--sampling_temp', '0.1',
                        '--seed', '37',
                        '--batch_size', '1'
                    ]

                    # Run the command
                    subprocess.run(command, capture_output=True, text=True)

                    # Usage
                    fasta_file_path = os.path.join(generator_directory, f"seqs/evocycle_{self.evo_cycle}_motifscaffolding_0.fa")
                    modified_seq = read_second_line_of_fasta(fasta_file_path)
                    logging.info(f"modified sequence after ProteinMPNN, {modified_seq}")

                    # insert the deletions back into the sequence:
                    modified_seq = reinsert_deletions(modified_seq, self.action_mask)

                else:
                    modified_seq = self.sequence

                modified_seq = list(modified_seq) # insert the new deletions
                for i, char in enumerate(self.action_mask):
                    if char=='-':
                        if modified_seq[i]!='-':
                            logging.info(f"deleting residue")
                            modified_seq[i] = '-'

                return ''.join(modified_seq)
        
        logging.info(f"Generating job complete. Results are in {self.outputs_directory}")

    def run(self):
        return self.run_generation()
