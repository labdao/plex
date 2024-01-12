import subprocess
import os
import pandas as pd
from utils import squeeze_seq

class StateGenerator:
    def __init__(self, evo_cycle, generator_list, sequence, cfg, outputs_directory):
        self.evo_cycle = evo_cycle
        self.generator_list = generator_list
        self.sequence = squeeze_seq(sequence)
        self.outputs_directory = outputs_directory
        self.cfg = cfg
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
            #     print(f"Running {generator}")
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

            
            elif generator=='delete+substitute':
                print(f"Running {generator}")

                alphabet = 'LAGVSERTIDPKQNFYMHWC'
        
                modified_seq = list(self.seed)
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