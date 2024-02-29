import pandas as pd
from utils import squeeze_seq

from generators.RFdiffProteinMPNN import RFdiffusionProteinMPNNGenerator
from generators.complete_sequence import complete_sequence_Generator
from generators.RME_generator import RMEGenerator

class GenerationArgs:
    def __init__(self, evo_cycle, sequence, permissibility_vector, df, cfg, outputs_directory, generator_name):
        self.evo_cycle = evo_cycle
        self.sequence = sequence
        self.permissibility_vector = permissibility_vector
        self.df = df
        self.cfg = cfg
        self.outputs_directory = outputs_directory
        self.generator_name = generator_name
        self.target = cfg.params.basic_settings.target_seq
        self.alphabet = cfg.params.basic_settings.alphabet
        self.max_levenshtein_step_size = cfg.params.basic_settings.max_levenshtein_step_size
        self.num_designs = cfg.params.RFdiffusion_settings.inference.num_designs
        self.num_seqs = cfg.params.pMPNN_settings.num_seqs
        self.hotspots = cfg.params.RFdiffusion_settings.hotspots

class Generator:
    def __init__(self, cfg, outputs_directory):

        self.cfg = cfg
        self.outputs_directory = outputs_directory

    def _get_generator(self, generator_name):
        
        if generator_name == 'RFdiff+ProteinMPNN':
            return RFdiffusionProteinMPNNGenerator()
        elif generator_name == 'RFdiff+ProteinMPNN+ESM2':
            return RMEGenerator()
        elif generator_name == 'complete_sequence':
            return complete_sequence_Generator()
        # ... add other generators to this list  ...
        else:
            raise ValueError(f"Unknown generator: {generator_name}")

    def run(self, t, sample_number, seed, permissibility_seed, permissibility_vector, df):

        generator_name = self.cfg.params.basic_settings.generator
        if t == 0:
            generator_name = 'complete_sequence'
        args = GenerationArgs(t, seed, permissibility_vector, df, self.cfg, self.outputs_directory, generator_name)
        generator = self._get_generator(generator_name)
        modified_seq, permissibility_vector = generator.generate(args)

        new_row = {
            't': int(t),
            'sample_number': int(sample_number),
            'seed': squeeze_seq(seed),
            'permissibility_seed': ''.join(permissibility_seed),
            '(levenshtein-distance, mask)': (None, squeeze_seq(permissibility_vector).replace('X', 'x')),
            'modified_seq': modified_seq,
            'permissibility_modified_seq': ''.join(permissibility_vector),
            'acceptance_flag': False
        }
        df = pd.concat([df, pd.DataFrame([new_row])], ignore_index=True)

        return modified_seq, permissibility_vector, df