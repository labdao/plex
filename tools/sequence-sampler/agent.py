import pandas as pd
import numpy as np
import random
import sys
import os
import sequence_transformer

def squeeze_seq(new_sequence):
    return ''.join(filter(lambda x: x != '-', new_sequence))

def permissible_exhaustive_deletion(t, df):
    # Iterate over rows where 't' column value is t-1
    for index, row in df[df['t'] == t-1].iterrows():
        original_seq = row['original_seq']
        variant_seqs = row['permissible_variant_seq']
        seed_flags = row['seed_flag']
        if (t-1)==0:
            permissibility_vector = row['permissibility_vectors'][0] # TD: clean this up. should write the input differently already.
        else:
            permissibility_vector = row['permissibility_vectors']

        # Check if variant_seqs and seed_flags are lists and have the same length
        if isinstance(variant_seqs, list) and isinstance(seed_flags, list) and len(variant_seqs) == len(seed_flags):
            # Iterate over each sequence and its corresponding seed value
            for variant_seq, seed in zip(variant_seqs, seed_flags):
                # Process only if seed is True
                if seed:
                    variant_seq_list = list(variant_seq)

                    if '+' in variant_seq_list: # check whether there is anything to delete
                        # Iterate over the length of the variant_sequence
                        for n in range(len(variant_seq_list)): # TD: clean up to variant_seq?
                            new_sequence = list(variant_seq)
                            new_permissibility_vector = list(permissibility_vector)
                            if permissibility_vector[n]=='+':
                                # Create a new sequence replace the character at position n
                                new_sequence[n] = '-'
                                new_permissibility_vector[n] = '-'

                                # If variant_seq is a list, convert it to a string by joining its elements # TD: check why variant_seq sometimes ends up as a list
                                if isinstance(variant_seq, list):
                                    variant_seq = ''.join(variant_seq)

                                # Append a new row to the data frame
                                new_row = pd.DataFrame({
                                    't': [t],
                                    'seed_seq': [variant_seq],
                                    'original_seq': [original_seq],
                                    'shortened_seq': [squeeze_seq(new_sequence)],
                                    'permissible_shortened_seq': [new_sequence],
                                    'permissibility_vectors': [new_permissibility_vector]
                                })
                                df = pd.concat([df, new_row], ignore_index=True)

    return df

def compute_log_likelihood(runner, mutated_sequence, LL_matrix):

    # Ensure that the length of the mutated sequence matches the number of columns in LL_matrix
    if len(mutated_sequence) != LL_matrix.shape[1]:
        raise ValueError("Length of mutated_sequence must match the number of columns in LL_matrix.")
    
    # Define the one-letter amino acid code
    amino_acid_code = ''.join(runner.amino_acids) # ESM is using 'LAGVSERTIDPKQNFYMHWC' ordering

    # Initialize total log likelihood
    total_log_likelihood = 0

    # Compute the total log likelihood of the mutated sequence
    for i, aa in enumerate(mutated_sequence):
        # Find the row index for this amino acid
        row_index = amino_acid_code.index(aa)
        
        # Add the log likelihood from the corresponding cell in LL_matrix
        total_log_likelihood += LL_matrix[row_index, i]

    return total_log_likelihood

def greedy_choice_residue(runner, LL_matrix):
    # Define the one-letter amino acid code
    amino_acid_code = ''.join(runner.amino_acids) # ESM is using 'LAGVSERTIDPKQNFYMHWC' ordering

    # Check if the LL_matrix has 20 rows corresponding to the amino acids
    if LL_matrix.shape[0] != len(amino_acid_code):
        raise ValueError("The LL_matrix should have 20 rows, one for each amino acid.")

    # Find the index of the maximum value in each column
    max_indices = np.argmax(LL_matrix, axis=0)

    # Map these indices to their corresponding amino acids
    greedy_seq = [amino_acid_code[index] for index in max_indices]

    return greedy_seq

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

def permissible_likelihood_based_mutation(t, df):

    # Initialize the ESM2Runner with the default model
    runner = sequence_transformer.ESM2Runner()

    # Iterate over rows where 't' column value is t
    for index, row in df[df['t'] == t].iterrows():
        shortened_seq = row['shortened_seq']
        permissible_shortened_seq = row['permissible_shortened_seq']
        permissibility_vectors = row['permissibility_vectors']
        mutated_sequences = []
        LL_mutated_sequence = []
        permissible_mutated_sequences = []

        LL_matrix = runner.token_masked_marginal_log_likelihood_matrix(shortened_seq) # TD: could save some compute time here, by only consider those residues which are permissible to mutation
        # print('check sum', np.sum(np.exp(LL_matrix), axis=0)) # convert to probabilities and compute sum for each column

        print('LLLLLLLLLL_matrix', LL_matrix)
        greedy_mutations = greedy_choice_residue(runner, LL_matrix)

        # Iterate over the length of the shortened_sequence
        for i in range(len(shortened_seq)):

            squeezed_permisibility_list = list(squeeze_seq(permissibility_vectors))
            if squeezed_permisibility_list[i]=='X' or squeezed_permisibility_list[i]=='+':
            # if squeezed_permisibility_list[i]=='X':

                # Mutate only one residue at a time
                mutated_seq = list(shortened_seq)
                mutated_seq[i] = greedy_mutations[i]
                mutated_sequences.append(''.join(mutated_seq))

                permissible_variant = infer_permissible_variant(mutated_seq, pattern=permissible_shortened_seq)
                permissible_mutated_sequences.append(permissible_variant) # TD: at the moment we only append if the sequence has been mutated after the deletion. should take care of the shortened sequences which cannot be varied; these would continue as purely deleting along the tree

                LL_greedy = compute_log_likelihood(runner, mutated_seq, LL_matrix)
                LL_mutated_sequence.append(LL_greedy)

        df.at[index, 'variant_seq'] = mutated_sequences
        df.at[index, 'permissible_variant_seq'] = permissible_mutated_sequences
        # Ensure the 'action_score' column exists
        if 'action_score' not in df.columns:
            df['action_score'] = None
        df.at[index, 'action_score'] = LL_mutated_sequence

    return df

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

# def permissible_actions(t, df):

#     # # pseudo code:
#     # - loop over all seed sequences and carry out all permissible actions on them, create a pool of candidate sequences
#     # - 

#     # Iterate over rows where 't' column value is t-1
#     for index, row in df[df['t'] == t-1].iterrows():
#         original_seq = row['original_seq']
#         variant_seqs = row['permissible_variant_seq']
#         seed_flags = row['seed_flag']
#         if (t-1)==0:
#             permissibility_vector = row['permissibility_vectors'][0] # TD: clean this up. should write the input differently already.
#         else:
#             permissibility_vector = row['permissibility_vectors']

#         # Check if variant_seqs and seed_flags are lists and have the same length
#         if isinstance(variant_seqs, list) and isinstance(seed_flags, list) and len(variant_seqs) == len(seed_flags):
#             # Iterate over each sequence and its corresponding seed value
#             for variant_seq, seed in zip(variant_seqs, seed_flags):
#                 # Process only if seed is True
#                 if seed:
#                     variant_seq_list = list(variant_seq)

#                     if '+' in variant_seq_list: # check whether there is anything to delete
#                         # Iterate over the length of the variant_sequence
#                         for n in range(len(variant_seq_list)): # TD: clean up to variant_seq?
#                             new_sequence = list(variant_seq)
#                             new_permissibility_vector = list(permissibility_vector)
#                             if permissibility_vector[n]=='+':
#                                 # Create a new sequence replace the character at position n
#                                 new_sequence[n] = '-'
#                                 new_permissibility_vector[n] = '-'

#                                 # If variant_seq is a list, convert it to a string by joining its elements # TD: check why variant_seq sometimes ends up as a list
#                                 if isinstance(variant_seq, list):
#                                     variant_seq = ''.join(variant_seq)

#                                 # Append a new row to the data frame
#                                 new_row = pd.DataFrame({
#                                     't': [t],
#                                     'seed_seq': [variant_seq],
#                                     'original_seq': [original_seq],
#                                     'shortened_seq': [squeeze_seq(new_sequence)],
#                                     'permissible_shortened_seq': [new_sequence],
#                                     'permissibility_vectors': [new_permissibility_vector]
#                                 })
#                                 df = pd.concat([df, new_row], ignore_index=True)

#     return df

class Agent:

    def __init__(self, t, df, reward, cfg):
        self.t = t
        self.df = df
        self.reward = reward
        self.policy_flag = cfg.params.basic_settings.policy_flag
        self.cfg = cfg

    def apply_policy(self):   

        if self.policy_flag == 'delete_and_mutate_ESM':

            # read seq from data frame (first step)
            if self.t == 1:  # Adjust formatting of input data frame in the first iteration

                # renaming sequence_number column to original_sequence and write the values
                self.df.rename(columns={'sequence_number': 'original_seq'}, inplace=True)
                self.df['original_seq'] = self.df['seq']

                # Inserting an empty column named 'shortened_seq'
                self.df.insert(2, 'seed_seq', '')
                self.df.insert(3, 'shortened_seq', '')

                # Renaming the column 'seq' to 'variant_seq' and converting values to lists
                self.df['seq'] = self.df['seq'].apply(lambda x: [x])
                self.df.rename(columns={'seq': 'variant_seq'}, inplace=True)

                self.df['permissibility_vectors'] = self.df['permissibility_vectors'].apply(lambda x: [x])

                # Adding a column 'seed_flag' after 'variant_seq' and setting its value to a list containing True
                variant_seq_index = self.df.columns.get_loc('variant_seq')
                self.df.insert(variant_seq_index + 1, 'seed_flag', [[True]] * len(self.df))

                # Add code to modify the 'variant_seq' strings and create a new column 'permissible_variant_seq'
                self.df['permissible_variant_seq'] = self.df.apply(
                    lambda row: [
                        ''.join(
                            char if pv_char != '-' else '-' 
                            for char, pv_char in zip(var_seq, pv)
                        ) 
                        for var_seq, pv in zip(row['variant_seq'], row['permissibility_vectors'])
                    ], axis=1
                )

            # perform exhaustive deletion and return a list of shortened_sequences
            # df = exhaustive_deletion(self.t, self.df)
            df = permissible_exhaustive_deletion(self.t, self.df)

            # select mutation based on greedy sampling
            # df = likelihood_based_mutation(self.t, df)
            df = permissible_likelihood_based_mutation(self.t, df)

            # compute the action constraint (Levenshtein distance)
            df = action_constraint(self.t, df)

            # ## action ranking
            # df = action_ranking(self.t, df)
        
            return df, pd.DataFrame.empty

        if self.policy_flag == 'free_evolution': # policy that does not force mutation and deletion; also keep past sequences (i.e. admit different length seqs in same evolutionary step)

            if self.t == 1: # read seq from data in first step

                # TD: no need for the shorted_seq anymore 

                # renaming sequence_number column to original_sequence and write the values
                self.df.rename(columns={'sequence_number': 'original_seq'}, inplace=True)
                self.df['original_seq'] = self.df['seq']

                # Inserting an empty column named 'shortened_seq'
                self.df.insert(2, 'seed_seq', '')
                self.df.insert(3, 'shortened_seq', '')

                # Renaming the column 'seq' to 'variant_seq' and converting values to lists
                self.df['seq'] = self.df['seq'].apply(lambda x: [x])
                self.df.rename(columns={'seq': 'variant_seq'}, inplace=True)

                self.df['permissibility_vectors'] = self.df['permissibility_vectors'].apply(lambda x: [x])

                # Adding a column 'seed_flag' after 'variant_seq' and setting its value to a list containing True
                variant_seq_index = self.df.columns.get_loc('variant_seq')
                self.df.insert(variant_seq_index + 1, 'seed_flag', [[True]] * len(self.df))

                # Add code to modify the 'variant_seq' strings and create a new column 'permissible_variant_seq'
                self.df['permissible_variant_seq'] = self.df.apply(
                    lambda row: [
                        ''.join(
                            char if pv_char != '-' else '-' 
                            for char, pv_char in zip(var_seq, pv)
                        ) 
                        for var_seq, pv in zip(row['variant_seq'], row['permissibility_vectors'])
                    ], axis=1
                )

            # # pseudo code:
            # - carry out all permissible actions on the seed sequences  
            #    df = permissible_actions(self.t, self.df)  
            # - from the 
            #    df = permissible_likelihood_based_mutation(self.t, df)        

            # # perform exhaustive deletion and return a list of shortened_sequences
            # # df = exhaustive_deletion(self.t, self.df)
            # df = permissible_exhaustive_deletion(self.t, self.df)

            # # select mutation based on greedy sampling
            # df = permissible_likelihood_based_mutation(self.t, df)

            # compute the action constraint (Levenshtein distance)
            df = action_constraint(self.t, df)
        
            return df, pd.DataFrame.empty