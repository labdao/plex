import pandas as pd
import numpy as np
import random
import sys
import os
import sequence_transformer
from state_generator import StateGenerator
# from state_scorer import StateScorer

def compute_log_likelihood(runner, mutated_sequence, LLmatrix):

    # Ensure that the length of the mutated sequence matches the number of columns in LLmatrix
    if len(mutated_sequence) != LLmatrix.shape[1]:
        raise ValueError("Length of mutated_sequence must match the number of columns in LLmatrix.")
    
    # Define the one-letter amino acid code
    amino_acid_code = ''.join(runner.amino_acids) # ESM is using 'LAGVSERTIDPKQNFYMHWC' ordering

    # Initialize total log likelihood
    total_log_likelihood = 0

    # Compute the total log likelihood of the mutated sequence
    squeezed_mutated_sequence = squeeze_seq(mutated_sequence)
    for i, aa in enumerate(squeezed_mutated_sequence):
        # Find the row index for this amino acid
        row_index = amino_acid_code.index(aa)
        
        # Add the log likelihood from the corresponding cell in LLmatrix
        total_log_likelihood += LLmatrix[row_index, i]

    return total_log_likelihood

def squeeze_seq(new_sequence):
    return ''.join(filter(lambda x: x != '-', new_sequence))

def acceptance_probability(ref_score, mod_score, T):

    k = 1.
    p_mod = np.exp((ref_score - mod_score)/(k*T)) # TD: think carefully about the sign in the exponent
    print('acceptance probability', np.minimum(1.,p_mod))

    return np.minimum(1.,p_mod)

def action_bouncer(ref_score, mod_score, T):

    p_mod = acceptance_probability(ref_score, mod_score, T)

    return random.random() < p_mod

def score_seq(seq): # TD: normalisation of LL by sequence length!?

    # Initialize the ESM2Runner with the default model
    runner = sequence_transformer.ESM2Runner()
    LLmatrix = runner.token_masked_marginal_log_likelihood_matrix(seq)

    LL = compute_log_likelihood(runner, seq, LLmatrix)

    return LL

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
        elif char == '-':
            action_mask.append('-')  
    
    action_mask = ''.join(action_mask)
    permissibility_vector = ''.join(permissibility_vector)

    return permissibility_vector, action_mask

def generate_proposed_state(seed, action_mask, cfg):

    generator = StateGenerator('simple_generator', seed, action_mask, cfg)
    modified_seq = generator.run()

    return modified_seq

def sample_action_mask(t, seed, permissibility_seed, action_residue_list, cfg, max_levenshtein_step_size):

    n_masks = 1
    alphabet = 'LAGVSERTIDPKQNFYMHWC'

    permissible_mask, levenshtein_step_size = sample_permissible_vector(seed, permissibility_seed, max_levenshtein_step_size, n_masks, alphabet)
    permissibility_vector, action_mask = sample_actions_for_mask(permissible_mask, permissibility_seed, alphabet)

    return permissibility_vector, action_mask, levenshtein_step_size

def score_sequence(seed, mod_seq, levenshtein_distance, LLmatrix_seed, runner): # TD: generalise to allow for plug-in of other scoring functions
    if mod_seq !=[]:
        if levenshtein_distance==0:
            LL_mod = compute_log_likelihood(runner, mod_seq, LLmatrix_seed)
        elif levenshtein_distance>0:
            LLmatrix_mod = runner.token_masked_marginal_log_likelihood_matrix(mod_seq)
            LL_mod = compute_log_likelihood(runner, mod_seq, LLmatrix_mod)

    return LL_mod


class Sampler:

    def __init__(self, t, seed, permissibility_seed, cfg):
        self.t = t
        self.seed = seed
        self.permissibility_seed = permissibility_seed
        self.cfg = cfg
        self.policy_flag = cfg.params.basic_settings.policy_flag
        self.temperature = cfg.params.basic_settings.temperature
        self.max_levenshtein_step_size = cfg.params.basic_settings.max_levenshtein_step_size

    def apply_policy(self):

        if self.policy_flag == 'policy_sampling':

            levenshtein_step_size = 2

            # Initialize the ESM2Runner with the default model
            runner = sequence_transformer.ESM2Runner()
            LLmatrix_seed = runner.token_masked_marginal_log_likelihood_matrix(self.seed)
            LL_seed = compute_log_likelihood(runner, self.seed, LLmatrix_seed)

            action_residue_list = []
            sample_number = 1
            accept_flag = False
            while accept_flag is False:
                print('sample number', sample_number)

                permissibility_vector, action_mask, levenshtein_distance = sample_action_mask(self.t, self.seed, self.permissibility_seed, action_residue_list, self.cfg, self.max_levenshtein_step_size)
                print('levenshtein, permissible vector, mask:', levenshtein_distance, permissibility_vector, action_mask)

                mod_seq = generate_proposed_state(self.seed, action_mask, self.cfg)

                LL_mod = score_sequence(self.seed, squeeze_seq(mod_seq), levenshtein_distance, LLmatrix_seed, runner)

                accept_flag = action_bouncer(LL_seed, LL_mod, self.temperature)
                print('action accepted', accept_flag)

                sample_number += 1
        
            return mod_seq, permissibility_vector, (levenshtein_distance, squeeze_seq(action_mask)), levenshtein_distance, action_mask