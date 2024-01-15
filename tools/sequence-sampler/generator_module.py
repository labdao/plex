import subprocess
import os
import pandas as pd
import random
from utils import squeeze_seq
from utils import generate_contig
from utils import read_last_line_of_fasta

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
                command = [
                    'python', 'RFdiffusion/scripts/run_inference.py',
                    '--output_prefix', os.path.join(generator_directory, '/motifscaffolding'),
                    '--model_directory_path', '/inputs/models',
                    '--input_pdb', self.df['absolute pdb path'],
                    '--num_designs', '1',
                    '--contigmap_contigs', contig
                ]

                result = subprocess.run(command, capture_output=True, text=True)

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
                path_to_PDB = os.path.join(generator_directory, '/motifscaffolding_1.pdb') # '/inputs/diffusion_test.pdb'
                output_dir = generator_directory
                chains_to_design = 'B'

                # Create the output directory if it doesn't exist
                os.makedirs(output_dir, exist_ok=True)

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
                fasta_file_path = os.path.join(generator_directory, '/motifscaffolding_1.fasta')
                modified_seq = read_last_line_of_fasta(fasta_file_path)
                print('modified sequence after ProteinMPNN', last_line)

                return ''.join(modified_seq) # TD: extract the sequence properly here from the seqs output directory!

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

# docker run command for rfdiffusion:

# docker run -it --rm --gpus all \
#   -v $HOME/models:$HOME/models \
#   -v $HOME/inputs:$HOME/inputs \
#   -v $HOME/outputs:$HOME/outputs \
#   rfdiffusion \
#   inference.output_prefix=$HOME/outputs/motifscaffolding \
#   inference.model_directory_path=$HOME/models \
#   inference.input_pdb=$HOME/inputs/5TPN.pdb \
#   inference.num_designs=3 \
#   'contigmap.contigs=[10-40/A163-181/10-40]'

                # docker run -it --rm --gpus all \
                #   -v $HOME/models:$HOME/models \
                #   -v $HOME/inputs:$HOME/inputs \
                #   -v $HOME/outputs:$HOME/outputs \
                #   rfdiffusion \
                #   inference.output_prefix=$HOME/outputs/motifscaffolding \
                #   inference.model_directory_path=$HOME/models \
                #   inference.input_pdb=$HOME/inputs/5TPN.pdb \
                #   inference.num_designs=3 \
                #   'contigmap.contigs=[10-40/A163-181/10-40]'

                # # last lines of some docker file
                # WORKDIR /app/RFdiffusion

                # ENV DGLBACKEND="pytorch"

                # ENTRYPOINT ["python3.9", "scripts/run_inference.py"]

                # inference.output_prefix=$HOME/outputs/motifscaffolding \
                # inference.model_directory_path=$HOME/models \
                # inference.input_pdb=$HOME/inputs/5TPN.pdb \
                # inference.num_designs=3 \
                # 'contigmap.contigs=[10-40/A163-181/10-40]'