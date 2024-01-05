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
    for i, aa in enumerate(mutated_sequence):
        # Find the row index for this amino acid
        row_index = amino_acid_code.index(aa)
        
        # Add the log likelihood from the corresponding cell in LLmatrix
        total_log_likelihood += LLmatrix[row_index, i]

    return total_log_likelihood

def squeeze_seq(new_sequence):
    return ''.join(filter(lambda x: x != '-', new_sequence))

def boltzmann(ref_score, mod_score):

    # p_mod = exp(...)
    p_mod = mod_score

    return p_mod

def score_seq(seq):

    # Initialize the ESM2Runner with the default model
    runner = sequence_transformer.ESM2Runner()
    LLmatrix = runner.token_masked_marginal_log_likelihood_matrix(seq)

    LL = compute_log_likelihood(runner, seq, LLmatrix)

    return LL

# Randomly select an element from permissibility_seed where the selected action can be applied
def select_random_permissible_residue(permissibility_seed, selected_action):
    # Define the criteria for selecting an element based on the action
    action_criteria = {
        'mutate': lambda x: x == '+' or x == 'X',
        'delete': lambda x: x == '+'
    }
    
    # Filter the permissibility_seed for elements where the action can be applied
    permissible_elements = [i for i, x in enumerate(permissibility_seed) if action_criteria[selected_action](x)]
    
    # Randomly choose from the permissible elements with equal probabilities
    if not permissible_elements:
        print("Warning: No permissible elements found for the selected action.")
        return None
    
    return random.choice(permissible_elements)

# Randomly select one action from the list of permissible actions
def select_random_permissible_action(permissibility_seed, action_probabilities=None):
    
    # Define what makes each action permissible
    permissibility_criteria = {
        'mutate': lambda permissibility_seed: '+' in permissibility_seed or 'X' in permissibility_seed,
        'delete': lambda permissibility_seed: '+' in permissibility_seed
    }

    permissible_actions = []
    for action, is_permissible in permissibility_criteria.items():
        if is_permissible(permissibility_seed):
            permissible_actions.append(action)

    # warning if no permissible actions are found
    if not permissible_actions:
        print("Warning: No permissible actions found.")
    if not permissible_actions:
        return None  # No action to select if the list is empty
    
    if action_probabilities is None:
        # If no probabilities are provided, select with uniform probability
        action_probabilities = [1. / len(permissible_actions)] * len(permissible_actions)
    
    return random.choices(permissible_actions, weights=action_probabilities, k=1)[0]

def MLL_mutation(runner, LLmatrix):
    # Define the one-letter amino acid code
    amino_acid_code = ''.join(runner.amino_acids) # ESM is using 'LAGVSERTIDPKQNFYMHWC' ordering

    # Check if the LLmatrix has 20 rows corresponding to the amino acids
    if LLmatrix.shape[0] != len(amino_acid_code):
        raise ValueError("The LLmatrix should have 20 rows, one for each amino acid.")

    # Find the index of the maximum value in each column
    max_indices = np.argmax(LLmatrix, axis=0)

    # Map these indices to their corresponding amino acids
    MLL_mutations = [amino_acid_code[index] for index in max_indices]

    return MLL_mutations

def apply_permissible_action(seed, permissibility_seed, selected_action, selected_residue, MLL_mutations, cfg):

    modified_seq = list(seed)
    modified_permissibility_seq = list(permissibility_seed)
    if selected_action=='mutate':
        modified_seq[selected_residue] = MLL_mutations[selected_residue]
        modified_permissibility_seq = permissibility_seed
    elif selected_action=='delete':
        modified_seq[selected_residue] = '-'
        modified_permissibility_seq[selected_residue] = '-' 

    return modified_seq, modified_permissibility_seq

def select_and_apply_random_permissible_action(t, seed, permissibility_seed, action_residue_list, MLL_mutations, cfg):

    selected_action = select_random_permissible_action(permissible_actions, action_probabilities=None)
    selected_residue = select_random_permissible_residue(permissibility_seed, selected_action)

    if action_residue_tuple not in action_residue_list: # append and modify, if not done before
        action_residue_list.append(action_residue_tuple)
        print('action-residue pair:', action_residue_tuple)
        modified_seq, modified_permissibility_seq = apply_permissible_action(seed, permissibility_seed, selected_action, selected_residue, MLL_mutations, cfg)
        
    return modified_seq, modified_permissibility_seq, action_residue_list

class Sampler:

    def __init__(self, t, seed, cfg):
        self.t = t
        self.seed = seed
        self.cfg = cfg
        self.policy_flag = cfg.params.basic_settings.policy_flag

    def apply_policy(self):

        if self.policy_flag == 'policy_sampling':

            # Initialize the ESM2Runner with the default model
            runner = sequence_transformer.ESM2Runner()
            LLmatrix_seed = runner.token_masked_marginal_log_likelihood_matrix(self.seed)
            LL_seed = score_seq(self.seed)
            p_ref = LL_seed # need to replace this later with proper implementation of boltzmann

            action_residue_list = []
            MLL_mutations = MLL_mutation(runner, LLmatrix_seed)
            sample_number = 1
            while p_mod < p_ref:
                print('sample number', sample_number)
                mod_seq = select_and_apply_random_permissible_action(t, seed, permissibility_seed, action_residue_list, MLL_mutations, cfg)
                if mod_seq !=[]:
                    LL_mod = score_seq(mod_seq) 
                    p_mod = boltzmann(LL_seed, LL_mod)
                sample_number += 1
        
            return mod_seq



### OLD CODE ###

# seed = mod_seq
# # add all info to data frame

# compute the action constraint (Levenshtein distance)
# df = action_constraint(self.t, df)

# if self.t == 1: # read seq from data in first step
                # # TD: no need for the shorted_seq anymore 

                # # renaming sequence_number column to original_sequence and write the values
                # self.df.rename(columns={'sequence_number': 'original_seq'}, inplace=True)
                # self.df['original_seq'] = self.df['seq']

                # # Inserting an empty column named 'shortened_seq'
                # self.df.insert(2, 'seed_seq', '')
                # self.df.insert(3, 'shortened_seq', '')

                # # Renaming the column 'seq' to 'variant_seq' and converting values to lists
                # self.df['seq'] = self.df['seq'].apply(lambda x: [x])
                # self.df.rename(columns={'seq': 'variant_seq'}, inplace=True)

                # self.df['permissibility_vectors'] = self.df['permissibility_vectors'].apply(lambda x: [x])

                # # Adding a column 'seed_flag' after 'variant_seq' and setting its value to a list containing True
                # variant_seq_index = self.df.columns.get_loc('variant_seq')
                # self.df.insert(variant_seq_index + 1, 'seed_flag', [[True]] * len(self.df))

                # # Add code to modify the 'variant_seq' strings and create a new column 'permissible_variant_seq'
                # self.df['permissible_variant_seq'] = self.df.apply(
                #     lambda row: [
                #         ''.join(
                #             char if pv_char != '-' else '-' 
                #             for char, pv_char in zip(var_seq, pv)
                #         ) 
                #         for var_seq, pv in zip(row['variant_seq'], row['permissibility_vectors'])
                #     ], axis=1
                # )

    # squeezed_permisibility_list = list(squeeze_seq(permissibility_vectors))
    # if squeezed_permisibility_list[i]=='X' or squeezed_permisibility_list[i]=='+':
    # # if squeezed_permisibility_list[i]=='X':

    #     # Mutate only one residue at a time
    #             mutated_seq = list(shortened_seq)
    #             mutated_seq[i] = greedy_mutations[i]
    #             mutated_sequences.append(''.join(mutated_seq))

    #             permissible_variant = infer_permissible_variant(mutated_seq, pattern=permissible_shortened_seq)
    #             permissible_mutated_sequences.append(permissible_variant) # TD: at the moment we only append if the sequence has been mutated after the deletion. should take care of the shortened sequences which cannot be varied; these would continue as purely deleting along the tree

        # df.at[index, 'variant_seq'] = mutated_sequences
        # df.at[index, 'permissible_variant_seq'] = permissible_mutated_sequences
        # # Ensure the 'action_score' column exists
        # if 'action_score' not in df.columns:
        #     df['action_score'] = None
        # df.at[index, 'action_score'] = LL_mutated_sequence