import subprocess
import os
import pandas as pd
import random
from utils import squeeze_seq
from utils import generate_contig
from utils import read_second_line_of_fasta

class StateGenerator:
    def __init__(self, evo_cycle, generator_list, sequence, action_mask, cfg, outputs_directory, df):
        self.evo_cycle = evo_cycle
        self.generator_list = generator_list
        self.sequence = sequence # to squeeze or not to squeeze?
        self.action_mask = action_mask
        self.cfg = cfg
        self.target = cfg.params.basic_settings.target_seq
        self.outputs_directory = outputs_directory
        self.df = df
        # take data frame as input and retrieve the pdb as absolute path

    def run_generation(self):

        print('\n')
        print("Running generating job...")
        df_generate = pd.DataFrame() # initialize data frame
        for generator in self.generator_list:

            generator_directory = os.path.join(self.outputs_directory, generator)
            if not os.path.exists(generator_directory):
                os.makedirs(generator_directory, exist_ok=True)

            if generator=='RFdiffusion+ProteinMPNN':
                print(f"Running {generator}")

                print(os.getcwd())

                contig = generate_contig(self.action_mask, self.target, starting_target_residue=None, end_target_residue=None)
                print('contig for diffusion', contig)
                # define arguments and run RFDiffusion

                # Set up the environment for the subprocess - required so that RFdiffussion can find its proper packages
                env = os.environ.copy()
                env['PYTHONPATH'] = "/app/RFdiffusion:" + env.get('PYTHONPATH', '')

                command = [
                    'python', 'RFdiffusion/scripts/run_inference.py',
                    f'inference.output_prefix={os.path.join(generator_directory, f"evocycle_{self.evo_cycle}_motifscaffolding")}',
                    'inference.model_directory_path=RFdiffusion/models',
                    f'inference.input_pdb={self.df["absolute pdb path"].iloc[0]}',
                    'inference.num_designs=1',
                    f'contigmap.contigs={[contig]}'
                ]

                result = subprocess.run(command, capture_output=True, text=True, env=env)

                # Check if the command was successful
                if result.returncode == 0:
                    print("Inference script ran successfully")
                    print(result.stdout)
                else:
                    print("Error running inference script")
                    print(result.stderr)
                
                print('Running MPNN')

                # Activate the conda environment 'mlfold'
                subprocess.run(['conda', 'activate', 'mlfold'], shell=True)

                # Define the paths and parameters
                path_to_PDB = os.path.join(generator_directory, f"evocycle_{self.evo_cycle}_motifscaffolding_0.pdb")
                output_dir = generator_directory
                chains_to_design = 'A'

                # Create the output directory if it doesn't exist
                os.makedirs(output_dir, exist_ok=True)

                print("pdb path", path_to_PDB)

                # Define the command and arguments
                command = [
                    'python', 'ProteinMPNN/protein_mpnn_run.py',
                    '--pdb_path', path_to_PDB,
                    '--pdb_path_chains', chains_to_design,
                    '--out_folder', output_dir,
                    '--num_seq_per_target', '1',
                    '--sampling_temp', '0.1',
                    '--seed', '37',
                    '--batch_size', '1'
                ]

                # Run the command
                subprocess.run(command, capture_output=True, text=True)

                # Usage
                fasta_file_path = os.path.join(generator_directory, f"seqs/evocycle_{self.evo_cycle}_motifscaffolding_0.fa")
                modified_seq = read_second_line_of_fasta(fasta_file_path)
                print('modified sequence after ProteinMPNN', modified_seq)

                return ''.join(modified_seq)

            elif generator=='delete+substitute':
                print(f"Running {generator}")

                alphabet = 'LAGVSERTIDPKQNFYMHWC'
        
                modified_seq = list(self.sequence)
                for i, char in enumerate(self.action_mask):
                    if char not in alphabet:
                        if char=='X':
                            print('applying mutation')
                            letter_options = [letter for letter in alphabet if letter != modified_seq[i]]
                            new_letter = random.choice(letter_options)
                            modified_seq[i] = new_letter
                        elif char=='-':
                            if modified_seq[i]!='-':
                                print('applying deletion')
                                modified_seq[i] = '-'

                return ''.join(modified_seq)
        
        print(f"Generating job complete. Results are in {self.outputs_directory}")

    def run(self):
        return self.run_generation()