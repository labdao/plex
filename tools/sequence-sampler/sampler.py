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

def compute_log_likelihood(sequence, LLmatrix): # TD: move into the scorer module, or even utils or sequence-transformer

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

def sequence_bouncer(t, df, cfg):

    T = cfg.params.basic_settings.temperature
    # weights = {'pseudolikelihood': .7, 'mean plddt': .2, 'affinity': .1}
    scoring_weights = cfg.params.basic_settings.scoring_weights
    scoring_weights = scoring_weights.split(',')
    weights = {'pseudolikelihood': float(scoring_weights[0]), 'mean plddt': float(scoring_weights[1]), 'affinity': float(scoring_weights[2])}
    DeltaE = 0.
    scoring_metrics = cfg.params.basic_settings.scoring_metrics
    scoring_weights = cfg.params.basic_settings.scoring_metrics
    for metric in scoring_metrics.split(','):

        ref_metric = df.iloc[t-1][metric]
        mod_metric = df.iloc[t][metric]
        # if metric=='pseudolikelihood': # useful when transforming metrics
        #     ref_metric = ...
        #     mod_metric = ...
        # if metric=='mean plddt':
        #     ref_metric = ...
        #     mod_metric = ...
        # if metric=='affinity':
        #     ref_metric = ...
        #     mod_metric = ...
        
        DeltaE += weights[metric] * (ref_metric - mod_metric)

    p_mod = np.exp(DeltaE / T)
    print('acceptance probability', np.minimum(1.,p_mod))

    return random.random() < p_mod

def sample_permissible_vector(seed, permissibility_seed, max_levenshtein_step_size, alphabet):
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
    
    action_mask = ''.join(action_mask)
    permissibility_vector = ''.join(permissibility_vector)

    return permissibility_vector, action_mask

def generate_proposed_state(t, seed, action_mask, cfg, outputs_directory, df):

    generator = StateGenerator(t, ['RFdiffusion+ProteinMPNN'], seed, action_mask, cfg, outputs_directory, df) # TD:
    modified_seq = generator.run()

    return modified_seq

def sample_action_mask(t, seed, permissibility_seed, action_residue_list, cfg, max_levenshtein_step_size):

    alphabet = 'LAGVSERTIDPKQNFYMHWC'

    permissible_mask, levenshtein_step_size = sample_permissible_vector(seed, permissibility_seed, max_levenshtein_step_size, alphabet)
    permissibility_vector, action_mask = sample_actions_for_mask(permissible_mask, permissibility_seed, alphabet)

    return permissibility_vector, action_mask, levenshtein_step_size

def score_sequence_fullmetrics(t, sequence, cfg, outputs_directory, df): # TD: receive df as arugment and write the additional scores to frame; generalise to allow for plug-in of other scoring functions
    if squeeze_seq(sequence) !=[]:
        scorer = StateScorer(t, ['ESM2', 'Colabfold', 'Prodigy'], sequence, cfg, outputs_directory) # Note: currently only doing AF2 scoring for the selected design.
        df_scorer, LLmatrix_mod = scorer.run()
        LL_mod = compute_log_likelihood(sequence, LLmatrix_mod) # TD: normalization by sequence length?

        if 'pseudolikelihood' not in df_scorer.columns:
            df_scorer['pseudolikelihood'] = None  # Initialize the column with None
        if t==1:
            df_scorer.at[0, 'pseudolikelihood'] = LL_mod
        else:
            df_scorer.iloc[-1, df_scorer.columns.get_loc('pseudolikelihood')] = LL_mod

        # supplement data frame by scores
        df = concatenate_to_df(df_scorer, df)

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
                self.df = score_sequence_fullmetrics(self.t, self.seed, self.cfg, self.outputs_directory, self.df)

            action_residue_list = []
            sample_number = 1
            accept_flag = False
            while accept_flag is False:
                print('sample number', sample_number)

                permissibility_vector, action_mask, levenshtein_distance = sample_action_mask(self.t, self.seed, self.permissibility_seed, action_residue_list, self.cfg, self.max_levenshtein_step_size)
                print('levenshtein, permissible vector, mask:', levenshtein_distance, permissibility_vector.replace('X', 'x'), action_mask.replace('X', 'x'))

                mod_seq = generate_proposed_state(self.t, self.seed, action_mask, self.cfg, self.outputs_directory, self.df)

                squeezed_action_mask = squeeze_seq(action_mask)
                new_row = {
                    't': int(self.t),
                    'sample_number': int(sample_number),
                    'seed': squeeze_seq(self.seed),
                    'permissibility_seed': ''.join(self.permissibility_seed),
                    '(levenshtein-distance, mask)': (levenshtein_distance, squeezed_action_mask.replace('X', 'x')),
                    'modified_seq': mod_seq,
                    'permissibility_modified_seq': ''.join(permissibility_vector),
                    'acceptance_flag': False
                }
                # concat the new row to the DataFrame
                self.df = pd.concat([self.df, pd.DataFrame([new_row])], ignore_index=True)

                self.df = score_sequence_fullmetrics(self.t, mod_seq, self.cfg, self.outputs_directory, self.df)

                accept_flag = sequence_bouncer(self.t, self.df, self.cfg)
                self.df.iloc[-1, self.df.columns.get_loc('acceptance_flag')] = accept_flag
                print('action accepted', accept_flag)

                sample_number += 1
        
            return mod_seq, permissibility_vector, (levenshtein_distance, squeeze_seq(action_mask)), levenshtein_distance, action_mask, self.df
