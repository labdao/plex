import pandas as pd
import numpy as np
import random
import sys
import os
import sequence_transformer

def infer_permissible_variant(mutated, pattern):
    # Initialize an index for the mutated string
    mutated_index = 0
    # Initialize the result list
    result = []

    # Loop through the pattern
    for char in pattern:
        if char != '-':
            # If the character in the pattern is different from the one in the mutated string
            # replace it with the mutated character
            if char != mutated[mutated_index]:
                result.append(mutated[mutated_index])
            else:
                result.append(char)
            mutated_index += 1
        else:
            # If the character is '-', keep it
            result.append(char)

    return result

def levenshtein_distance(s1, s2):
    if len(s1) < len(s2):
        return levenshtein_distance(s2, s1)

    # len(s1) >= len(s2)
    if len(s2) == 0:
        return len(s1)

    previous_row = range(len(s2) + 1)
    for i, c1 in enumerate(s1):
        current_row = [i + 1]
        for j, c2 in enumerate(s2):
            insertions = previous_row[j + 1] + 1
            deletions = current_row[j] + 1
            substitutions = previous_row[j] + (c1 != c2)
            current_row.append(min(insertions, deletions, substitutions))
        previous_row = current_row

    return previous_row[-1]

def action_constraint(t, df):
    # Ensure the 'action constraint' column exists
    if 'action_constraint' not in df.columns:
        df['action_constraint'] = None

    # Iterate over rows where 't' column value is t
    for index, row in df[(df['t'] == t)].iterrows():
        shortened_seq = row['shortened_seq']
        variant_seqs = row['variant_seq']
        levenshtein_distances = []

        # Compute Levenshtein distance for each sequence in the variant_seq list
        for variant_seq in variant_seqs:
            distance = levenshtein_distance(shortened_seq, variant_seq)
            levenshtein_distances.append(distance)

        # Update the 'action_constraint' column with the list of distances
        df.at[index, 'action_constraint'] = levenshtein_distances

    return df

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

def sample_permissible_vector(seed, permissibility_seed, levenshtein_step_size, n_masks, alphabet):
    # Identify the indices where permissibility_seed has 'X' or '+'
    permissible_indices = [i for i, char in enumerate(permissibility_seed) if char in ['X', '+']]

    # Randomly select indices based on the levenshtein_step_size
    selected_indices = random.sample(permissible_indices, min(levenshtein_step_size, len(permissible_indices)))
        
    # Create a new mask based on the selected indices
    action_mask = list(seed)
    for index in selected_indices:
            action_mask[index] = permissibility_seed[index]
        
    action_mask = ''.join(action_mask)
    
    return action_mask #, permissibility_seed

def sample_actions_for_mask(permissible_mask, permissibility_vector, alphabet):
    print()
    print('permissible_mask', permissible_mask)
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
            elif random_action == 'delete':
                action_mask.append('-')
                permissibility_vector[i] = '-' # important line
    
    action_mask = ''.join(action_mask)
    permissibility_vector = ''.join(permissibility_vector) # convert to string

    return permissibility_vector, action_mask

def apply_action(seed, selected_residue, selected_action, cfg):

    alphabet = 'LAGVSERTIDPKQNFYMHWC'

    modified_seq = list(seed)
    # modified_permissibility_seq = list(permissibility_seed)
    if selected_action=='X':
        print('applying mutation')
        letter_options = [letter for letter in alphabet if letter != modified_seq[selected_residue]]
        new_letter = random.choice(letter_options)
        modified_seq[selected_residue] = new_letter
        # modified_permissibility_seq = permissibility_seed
    elif selected_action=='-':
        print('applying deletion')
        letter_options = [letter for letter in alphabet if letter != modified_seq[selected_residue]]
        new_letter = random.choice(letter_options)
        modified_seq[selected_residue] = '-'
        # modified_permissibility_seq = permissibility_seed

    return modified_seq #, modified_permissibility_seq

def apply_action_vector(seed, action_mask, cfg):

    alphabet = 'LAGVSERTIDPKQNFYMHWC'
    
    modified_seq = list(seed)
    for i, char in enumerate(action_mask):
        if char not in alphabet:
            action = char
            modified_seq = apply_action(seed, i, action, cfg)

    return modified_seq

def sample_action_mask(t, seed, permissibility_seed, action_residue_list, cfg, max_levenshtein_step_size): # 

    max_levenshtein_step_size = 1

    n_masks = 1
    alphabet = 'LAGVSERTIDPKQNFYMHWC'

    permissible_mask = sample_permissible_vector(seed, permissibility_seed, max_levenshtein_step_size, n_masks, alphabet)
    permissibility_vector, action_mask = sample_actions_for_mask(permissible_mask, permissibility_seed, alphabet) # TD: check whether this line is correct!
    print('permissible_vector after', permissibility_vector)

    return permissibility_vector, action_mask

def propose_state(t, seed, action_mask, cfg):

    modified_seq = apply_action_vector(seed, action_mask, cfg)

    return squeeze_seq(modified_seq)

def score_sequence(seed, mod_seq, LLmatrix_seed, runner): # TD: generalise to allow for plug-in of other scoring functions
    if mod_seq !=[]:
        if len(mod_seq)==len(seed): # TD: change to, if levenshtein distance == 0
            LL_mod = compute_log_likelihood(runner, mod_seq, LLmatrix_seed)
        elif len(mod_seq)<len(seed): # TD: change to, if levenshtein distance > 0
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
            # MLL_mutations = MLL_mutation(runner, LLmatrix_seed)
            sample_number = 1
            accept_flag = False
            while accept_flag is False:
                print('sample number', sample_number)

                permissibility_vector, action_mask = sample_action_mask(self.t, self.seed, self.permissibility_seed, action_residue_list, self.cfg, self.max_levenshtein_step_size)
                print('permissible vector, mask', permissibility_vector, action_mask)

                mod_seq = propose_state(self.t, self.seed, action_mask, self.cfg)

                LL_mod = score_sequence(self.seed, mod_seq, LLmatrix_seed, runner)

                accept_flag = action_bouncer(LL_seed, LL_mod, self.temperature)
                print('action accepted', accept_flag)

                sample_number += 1
        
            return mod_seq, permissibility_vector, (levenshtein_step_size, squeeze_seq(action_mask)), levenshtein_step_size, action_mask


###### OLD CODE ######
# mod_seq, mod_permissibility_seq, action_residue_list, selected_action, action_residue_pair = select_and_apply_random_permissible_action(self.t, self.seed, self.permissibility_seed, action_residue_list, MLL_mutations, self.cfg)


# # Randomly select an element from permissibility_seed where the selected action can be applied
# def select_random_permissible_residue(permissibility_seed, selected_action):
#     # Define the criteria for selecting an element based on the action
#     action_criteria = {
#         'mutate': lambda x: x == '+' or x == 'X',
#         'delete': lambda x: x == '+'
#     }
    
#     # Filter the permissibility_seed for elements where the action can be applied
#     permissible_elements = [i for i, x in enumerate(permissibility_seed) if action_criteria[selected_action](x)]
    
#     # Randomly choose from the permissible elements with equal probabilities
#     if not permissible_elements:
#         print("Warning: No permissible elements found for the selected action.")
#         return None
    
#     return random.choice(permissible_elements)

# # Randomly select one action from the list of permissible actions
# def select_random_permissible_action(permissibility_seed, action_probabilities=None):
    
#     # Define what makes each action permissible
#     permissibility_criteria = {
#         'mutate': lambda permissibility_seed: '+' in permissibility_seed or 'X' in permissibility_seed,
#         'delete': lambda permissibility_seed: '+' in permissibility_seed
#     }

#     permissible_actions = []
#     for action, is_permissible in permissibility_criteria.items():
#         if is_permissible(permissibility_seed):
#             permissible_actions.append(action)

#     # warning if no permissible actions are found
#     if not permissible_actions:
#         print("Warning: No permissible actions found.")
#     if not permissible_actions:
#         return None  # No action to select if the list is empty
    
#     if action_probabilities is None:
#         # If no probabilities are provided, select with uniform probability
#         action_probabilities = [1. / len(permissible_actions)] * len(permissible_actions)
    
#     return random.choices(permissible_actions, weights=action_probabilities, k=1)[0]

# def apply_permissible_action(seed, permissibility_seed, selected_action, selected_residue, MLL_mutations, cfg):

#     modified_seq = list(seed)
#     modified_permissibility_seq = list(permissibility_seed)
#     if selected_action=='mutate':
#         aa_alphabet = 'LAGVSERTIDPKQNFYMHWC'
#         aa_options = [aa for aa in aa_alphabet if aa != modified_seq[selected_residue]]
#         new_amino_acid = random.choice(aa_options)
#         modified_seq[selected_residue] = new_amino_acid
#         modified_permissibility_seq = permissibility_seed
#     elif selected_action=='delete':
#         modified_seq[selected_residue] = '-'
#         modified_permissibility_seq[selected_residue] = '-' 

#     return modified_seq, modified_permissibility_seq

# def select_and_apply_random_permissible_action(t, seed, permissibility_seed, action_residue_list, MLL_mutations, cfg):

#     selected_action = select_random_permissible_action(permissibility_seed, action_probabilities=None)
#     selected_residue = select_random_permissible_residue(permissibility_seed, selected_action)

#     action_residue_pair = (selected_residue, selected_action)
#     print('action-residue pair:', action_residue_pair)

#     if action_residue_pair not in action_residue_list: # append and modify, if action-residue pair has not been sampled previously
#         action_residue_list.append(action_residue_pair)
#         mod_seq, modified_permissibility_seq = apply_permissible_action(seed, permissibility_seed, selected_action, selected_residue, MLL_mutations, cfg)

#     return mod_seq, modified_permissibility_seq, action_residue_list, selected_action, action_residue_pair

# def MLL_mutation(runner, LLmatrix):
#     # Define the one-letter amino acid code
#     amino_acid_code = ''.join(runner.amino_acids) # ESM is using 'LAGVSERTIDPKQNFYMHWC' ordering

#     # Check if the LLmatrix has 20 rows corresponding to the amino acids
#     if LLmatrix.shape[0] != len(amino_acid_code):
#         raise ValueError("The LLmatrix should have 20 rows, one for each amino acid.")

#     # Find the index of the maximum value in each column
#     max_indices = np.argmax(LLmatrix, axis=0)

#     # Map these indices to their corresponding amino acids
#     MLL_mutations = [amino_acid_code[index] for index in max_indices]

#     return MLL_mutations

# def apply_action(seed, permissibility_seed, selected_action, selected_residue, cfg):

#     modified_seq = list(seed)
#     modified_permissibility_seq = list(permissibility_seed)
#     if selected_action=='mutate':
#         print('applying mutation')
#         aa_alphabet = 'LAGVSERTIDPKQNFYMHWC'
#         aa_options = [aa for aa in aa_alphabet if aa != modified_seq[selected_residue]]
#         new_amino_acid = random.choice(aa_options)
#         modified_seq[selected_residue] = new_amino_acid
#         modified_permissibility_seq = permissibility_seed
#     elif selected_action=='delete':
#         modified_seq[selected_residue] = '-'
#         modified_permissibility_seq[selected_residue] = '-' 

#     return modified_seq, modified_permissibility_seq

# def apply_action_vector(seed, permissibility_seed, action_vector, cfg):
#     modified_seq = list(seed)

#     for i, action in enumerate(action_vector):
#         if action != 'none':
#             modified_seq, modified_permissibility_seq = apply_action(seed, permissibility_seed, action, i, cfg)

#     return modified_seq, modified_permissibility_seq # need to treat the case where all actions are none / how can this happen in the first place?