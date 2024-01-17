import pandas as pd
import numpy as np
import random
import sys
import os
import sequence_transformer
from generator_module import StateGenerator
from scorer_module import StateScorer
from utils import squeeze_seq
from utils import write_af2_update
from utils import concatenate_to_df

def compute_log_likelihood(sequence, LLmatrix):

    sequence = squeeze_seq(sequence)

    # Ensure that the length of the mutated sequence matches the number of columns in LLmatrix
    if len(sequence) != LLmatrix.shape[1]:
        raise ValueError("Length of sequence must match the number of columns in LLmatrix.")
    
    # Define the one-letter amino acid code
    # amino_acid_code = ''.join(runner.amino_acids) # ESM is using 'LAGVSERTIDPKQNFYMHWC' ordering
    amino_acid_code = ''.join('LAGVSERTIDPKQNFYMHWC')

    # Initialize total log likelihood
    total_log_likelihood = 0

    # Compute the total log likelihood of sequence
    for i, aa in enumerate(sequence):
        # Find the row index for this amino acid
        row_index = amino_acid_code.index(aa)
        
        # Add the log likelihood from the corresponding cell in LLmatrix
        total_log_likelihood += LLmatrix[row_index, i]

    return total_log_likelihood

def acceptance_probability(ref_score, mod_score, T):

    k = 1.
    p_mod = np.exp((ref_score - mod_score)/(k*T)) # TD: think carefully about the sign in the exponent
    print('acceptance probability', np.minimum(1.,p_mod))

    return np.minimum(1.,p_mod)

def action_bouncer(ref_score, mod_score, T):

    p_mod = acceptance_probability(ref_score, mod_score, T)

    return random.random() < p_mod

def sample_permissible_vector(seed, permissibility_seed, max_levenshtein_step_size, n_masks, alphabet):
    # Identify the indices where permissibility_seed has 'X' or '+'
    permissible_indices = [i for i, char in enumerate(permissibility_seed) if char in ['X', '+']]

    # Randomly select a sample size between 1 and the minimum of levenshtein_step_size and the length of permissible_indices
    levenshtein_step_size = random.randint(1, min(max_levenshtein_step_size, len(permissible_indices)))

    # Randomly select indices based on the levenshtein_step_size
    selected_indices = random.sample(permissible_indices, levenshtein_step_size)
        
    # Create a new mask based on the selected indices
    action_mask = list(seed)
    for index in selected_indices:
            action_mask[index] = permissibility_seed[index]
        
    action_mask = ''.join(action_mask)
    
    return action_mask, levenshtein_step_size

def sample_actions_for_mask(permissible_mask, permissibility_vector, alphabet):
    action_vector = []
    action_mask = []
    permissibility_vector = list(permissibility_vector)
    for i, char in enumerate(permissible_mask):
        if char in alphabet:
            action_vector.append('none')
            action_mask.append(char)
        elif char == 'X':
            action_vector.append('mutate')
            action_mask.append('X')
        elif char == '+':
            random_action = random.choice(['mutate', 'delete'])
            action_vector.append(random_action)
            if random_action == 'mutate':
                action_mask.append('X')
            elif random_action == 'delete': # important case
                action_mask.append('-')        
                permissibility_vector[i] = '-'
        # elif char == '-':
        #     action_mask.append('-')  
    
    action_mask = ''.join(action_mask)
    permissibility_vector = ''.join(permissibility_vector)

    return permissibility_vector, action_mask

def generate_proposed_state(t, seed, action_mask, cfg, outputs_directory, df):

    generator = StateGenerator(t, ['RFdiffusion+ProteinMPNN'], seed, action_mask, cfg, outputs_directory, df) # TD:
    modified_seq = generator.run()

    return modified_seq

def sample_action_mask(t, seed, permissibility_seed, action_residue_list, cfg, max_levenshtein_step_size):

    n_masks = 1
    alphabet = 'LAGVSERTIDPKQNFYMHWC'

    permissible_mask, levenshtein_step_size = sample_permissible_vector(seed, permissibility_seed, max_levenshtein_step_size, n_masks, alphabet)
    permissibility_vector, action_mask = sample_actions_for_mask(permissible_mask, permissibility_seed, alphabet)

    return permissibility_vector, action_mask, levenshtein_step_size

def score_sequence(t, seed, mod_seq, levenshtein_distance, LLmatrix_seed, cfg, outputs_directory): # TD: receive df as arugment and write the additional scores to frame; generalise to allow for plug-in of other scoring functions
    if squeeze_seq(mod_seq) !=[]:
        if levenshtein_distance==0:
            LL_mod = compute_log_likelihood(runner, mod_seq, LLmatrix_seed)
        elif levenshtein_distance>0:

            scorer = StateScorer(t, ['ESM2'], mod_seq, cfg, outputs_directory) # Note: currently only doing AF2 scoring for the selected design.
            df, LLmatrix_mod = scorer.run()
            # LLmatrix_mod = df.at[0, 'LLmatrix_sequence']
            LL_mod = compute_log_likelihood(mod_seq, LLmatrix_mod) # TD: normalization by sequence length?

    return LL_mod


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

            scorer = StateScorer(self.t, ['ESM2', 'AF2', 'Prodigy'], self.seed, self.cfg, self.outputs_directory)
            df, LLmatrix_seed = scorer.run()
            # LLmatrix_seed = df.at[0, 'LLmatrix_sequence']
            LL_seed = compute_log_likelihood(self.seed, LLmatrix_seed)
            if 'likelihood' not in df.columns:
                df['likelihood'] = None  # Initialize the column with None
            df.at[0, 'likelihood'] = LL_seed

            # supplement data frame by scores
            self.df = concatenate_to_df(df, self.df)

            action_residue_list = []
            sample_number = 1
            accept_flag = False
            while accept_flag is False:
                print('sample number', sample_number)

                permissibility_vector, action_mask, levenshtein_distance = sample_action_mask(self.t, self.seed, self.permissibility_seed, action_residue_list, self.cfg, self.max_levenshtein_step_size)
                print('levenshtein, permissible vector, mask:', levenshtein_distance, permissibility_vector, action_mask)

                mod_seq = generate_proposed_state(self.t, self.seed, action_mask, self.cfg, self.outputs_directory, df)

                LL_mod = score_sequence(self.t, self.seed, mod_seq, levenshtein_distance, LLmatrix_seed, self.cfg, self.outputs_directory) # TD: pass df to function

                accept_flag = action_bouncer(LL_seed, LL_mod, self.temperature) # rejection-sampling
                print('action accepted', accept_flag)

                squeezed_action_mask = squeeze_seq(action_mask)
                new_row = {
                    't': int(self.t),
                    'sample_number': int(sample_number),
                    'seed': squeeze_seq(self.seed),
                    'permissibility_seed': ''.join(self.permissibility_seed),
                    '(levenshtein-distance, mask)': (levenshtein_distance, squeezed_action_mask.replace('X', 'x')),
                    'modified_seq': mod_seq,
                    'permissibility_modified_seq': ''.join(permissibility_vector),
                    'acceptance_flag': accept_flag 
                }

                # Append the new row to the DataFrame
                df = pd.concat([self.df, pd.DataFrame([new_row])], ignore_index=True)

                # supplement data frame with sample_number and accept_flag, and return df

                sample_number += 1
        
            return mod_seq, permissibility_vector, (levenshtein_distance, squeeze_seq(action_mask)), levenshtein_distance, action_mask, df
