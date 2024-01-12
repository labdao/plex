import random

class StateGenerator:
    def __init__(self, generator_type, seed, action_mask, cfg):
        self.generator_type = generator_type
        self.seed = seed
        self.action_mask = action_mask
        self.cfg = cfg

    def generate_state(self):

        if self.generator_type=='simple_generator':

            alphabet = 'LAGVSERTIDPKQNFYMHWC'
    
            modified_seq = list(self.seed)
            for i, char in enumerate(self.action_mask):
                if char not in alphabet:
                    if char=='X':
                        print('applying mutation')
                        letter_options = [letter for letter in alphabet if letter != modified_seq[i]]
                        new_letter = random.choice(letter_options)
                        modified_seq[i] = new_letter
                    elif char=='-': # could be picky here only perform the deletion if it has not been applied to that position before
                        if modified_seq[i]!='-':
                            print('applying deletion')
                            modified_seq[i] = '-'

            return ''.join(modified_seq)

        elif self.generator_type=='diffusion_generator':

            # run alpha fold
            # run RFdiffusion

            return ''.join(modified_seq)

    def run(self):
        return self.generate_state()

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

# TD: from action mask, generate contig