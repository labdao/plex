from utils import squeeze_seq
import logging


class Sampler:

    def __init__(self,cfg, outputs_directory, generator, selector, scorer, evolve, n_samples):

        self.cfg = cfg
        self.outputs_directory = outputs_directory
        self.scorer = scorer
        self.generator = generator
        self.selector = selector
        self.evolve = evolve
        self.n_samples = n_samples

    def _create_permissible_vector(self, seed, permissibility_seed):

        logging.info(f"permissibility_seed {permissibility_seed}")

        permissible_indices = [i for i, char in enumerate(permissibility_seed) if char in ['X', '*']]
        selected_indices = permissible_indices
                
        permissibile_mask = list(seed)
        for i, char in enumerate(permissibility_seed):
            if char=='-':
                permissibile_mask[i] = '-'
            elif i in selected_indices:
                permissibile_mask[i] = permissibility_seed[i]
            
        permissibile_mask = ''.join(permissibile_mask)
        
        return permissibile_mask

    def run(self, t, seed, permissibility_seed, df):

        if t==1:
            df = self.scorer.run(t-1, seed, df)

        sample_number = 1
        accept_flag = False
        while accept_flag is False:
            logging.info(f"sample number {sample_number}")
                
            permissibility_vector = self._create_permissible_vector(seed, permissibility_seed)

            mod_seq, permissibility_mod_seq, df = self.generator.run(t, sample_number, seed, permissibility_seed, permissibility_vector, df)

            if squeeze_seq(mod_seq) !=[]:
                df = self.scorer.run(t, mod_seq, df)

            if self.evolve == False and sample_number >= self.n_samples:
                df.iloc[-1, df.columns.get_loc('acceptance_flag')] = False
                break
            elif self.evolve == True:
                accept_flag = self.selector.run(t, df, accept_flag)
                df.iloc[-1, df.columns.get_loc('acceptance_flag')] = accept_flag

            sample_number += 1
        
        
        return mod_seq, permissibility_mod_seq, df