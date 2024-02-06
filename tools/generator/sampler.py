import pandas as pd
import random
from generator_module import StateGenerator
from scorer_module import StateScorer
from utils import squeeze_seq
import logging

def sequence_bouncer(t, df, cfg, accept_flag):

    if cfg.params.basic_settings.bouncer_flag=='open-door':

        acceptance_probability = 1.0
        logging.info(f"acceptance probability: {acceptance_probability}")
        accept_flag = True
        logging.info(f"accept_flag: {accept_flag}")

        return True

def sample_permissible_vector(seed, permissibility_seed, max_levenshtein_step_size, alphabet):

    logging.info(f"permissibility_seed {permissibility_seed}")
    permissible_indices = [i for i, char in enumerate(permissibility_seed) if char in ['X']]

    step_size_sampling = False # TD: possibly set this flag as input through cfg
    if step_size_sampling == True:
        # randomly select indices based on random levenshtein_step_size
        levenshtein_step_size = random.randint(1, min(max_levenshtein_step_size, len(permissible_indices)))
        selected_indices = random.sample(permissible_indices, levenshtein_step_size)
    else:
        selected_indices = permissible_indices
        levenshtein_step_size = max_levenshtein_step_size
            
    # Create a new mask based on the selected indices
    permissibile_mask = list(seed)
    for i, char in enumerate(permissibility_seed):
        if char=='-':
            permissibile_mask[i] = '-'
        elif i in selected_indices:
            permissibile_mask[i] = permissibility_seed[i]
        
    permissibile_mask = ''.join(permissibile_mask)
    
    return permissibile_mask, levenshtein_step_size

def sample_actions_for_mask(permissible_mask, permissibility_vector, alphabet):
    action_mask = []
    permissibility_vector = list(permissibility_vector)
    for i, char in enumerate(permissible_mask):
        if char in alphabet:
            action_mask.append(char)
        elif char == 'X':
            action_mask.append('X')
        # elif char == '+':
        #     action_mask.append('+')
        # elif char == '+': # originally used to create a mask that already encodes the deletion; without it, the action mask is equal to the the permissibility vector
        #     random_action = random.choice(['mutate', 'delete'])
        #     if random_action == 'mutate':
        #         action_mask.append('X')
        #     elif random_action == 'delete': # important case; it can be debated whether this should maybe be part of the generator
        #         action_mask.append('-')        
        #         permissibility_vector[i] = '-' 
        elif char == '-':
            action_mask.append('-')      

    action_mask = ''.join(action_mask)
    permissibility_vector = ''.join(permissibility_vector)

    return permissibility_vector, action_mask

def generate_proposed_state(t, seed, action_mask, cfg, outputs_directory, df, permissibility_vector):

    generator = StateGenerator(t, seed, action_mask, cfg, outputs_directory, df, permissibility_vector)
    modified_seq, modified_permissibility_vector = generator.run()

    return modified_seq, modified_permissibility_vector

def sample_action_mask(t, seed, permissibility_seed, action_residue_list, cfg, max_levenshtein_step_size):

    alphabet = 'LAGVSERTIDPKQNFYMHWC'

    permissible_mask, levenshtein_step_size = sample_permissible_vector(seed, permissibility_seed, max_levenshtein_step_size, alphabet)
    permissibility_vector, action_mask = sample_actions_for_mask(permissible_mask, permissibility_seed, alphabet)

    return permissibility_vector, action_mask, levenshtein_step_size

def score_sequence_fullmetrics(t, sequence, cfg, outputs_directory, df):
    if squeeze_seq(sequence) !=[]:

        scorer = StateScorer(t, sequence, cfg, outputs_directory, df)
        df = scorer.run()

    return df


class Sampler:

    def __init__(self, t, seed, permissibility_seed, cfg, outputs_directory, df):
        self.t = t
        self.seed = seed
        self.permissibility_seed = permissibility_seed
        self.cfg = cfg
        self.policy_flag = cfg.params.basic_settings.policy_flag
        self.temperature = cfg.params.basic_settings.temperature
        self.max_levenshtein_step_size = cfg.params.basic_settings.max_levenshtein_step_size
        self.outputs_directory = outputs_directory
        self.df = df

    def apply_policy(self):

        if self.policy_flag == 'policy_sampling':

            if self.t==1: # compute scores for initial sequence
                self.df = score_sequence_fullmetrics(self.t-1, self.seed, self.cfg, self.outputs_directory, self.df)

            action_residue_list = []
            sample_number = 1
            accept_flag = False
            while accept_flag is False:
                logging.info(f"sample number {sample_number}")

                permissibility_vector, action_mask, levenshtein_distance = sample_action_mask(self.t, self.seed, self.permissibility_seed, action_residue_list, self.cfg, self.max_levenshtein_step_size)

                mod_seq, permissibility_mod_seq = generate_proposed_state(self.t, self.seed, action_mask, self.cfg, self.outputs_directory, self.df, permissibility_vector)
                if permissibility_mod_seq is None:
                    permissibility_mod_seq = permissibility_vector

                # concat new row to data frame
                squeezed_action_mask = squeeze_seq(action_mask)
                new_row = {
                    't': int(self.t),
                    'sample_number': int(sample_number),
                    'seed': squeeze_seq(self.seed),
                    'permissibility_seed': ''.join(self.permissibility_seed),
                    '(levenshtein-distance, mask)': (levenshtein_distance, squeezed_action_mask.replace('X', 'x')),
                    'modified_seq': mod_seq,
                    'permissibility_modified_seq': ''.join(permissibility_mod_seq),
                    'acceptance_flag': False
                }
                self.df = pd.concat([self.df, pd.DataFrame([new_row])], ignore_index=True)

                self.df = score_sequence_fullmetrics(self.t, mod_seq, self.cfg, self.outputs_directory, self.df)

                accept_flag = sequence_bouncer(self.t, self.df, self.cfg, accept_flag)

                self.df.iloc[-1, self.df.columns.get_loc('acceptance_flag')] = accept_flag

                sample_number += 1
        

            return mod_seq, permissibility_mod_seq, (levenshtein_distance, squeeze_seq(action_mask)), levenshtein_distance, action_mask, self.df