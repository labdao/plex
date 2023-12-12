import pandas as pd
import random
import sys
import os

def exhaustive_deletion(t, df):
    # Iterate over rows where 't' column value is t-1
    for index, row in df[df['t'] == t-1].iterrows():
        original_seq = row['original_seq']
        variant_seqs = row['variant_seq']
        seed_flags = row['seed_flag']

        # Check if variant_seqs and seed_flags are lists and have the same length
        if isinstance(variant_seqs, list) and isinstance(seed_flags, list) and len(variant_seqs) == len(seed_flags):
            # Iterate over each sequence and its corresponding seed value
            for variant_seq, seed in zip(variant_seqs, seed_flags):
                # Process only if seed is True
                if seed:
                    variant_seq_list = list(variant_seq)

                    # Iterate over the length of the variant_sequence
                    for n in range(len(variant_seq_list)):
                        # Create a new sequence excluding the character at position n
                        new_sequence = ''.join(variant_seq_list[:n] + variant_seq_list[n+1:])
                        # Append a new row to the data frame
                        new_row = pd.DataFrame({'t': [t], 'original_seq': original_seq, 'shortened_seq': [new_sequence]})
                        df = pd.concat([df, new_row], ignore_index=True)

    return df

def mutate_single_residue(t, df):
    # Define the one-letter amino acid alphabet
    alphabet = 'ACDEFGHIKLMNPQRSTVWY'

    # Iterate over rows where 't' column value is t
    for index, row in df[df['t'] == t].iterrows():
        original_seq = row['shortened_seq']
        mutated_sequences = []

        # Iterate over the length of the shortened_sequence
        for i in range(len(original_seq)):

            # Mutate only one residue at a time

            mutated_seq = list(original_seq)
            permissible_aas = [aa for aa in alphabet if aa != mutated_seq[i]]
            mutated_seq[i] = random.choice(permissible_aas)
            mutated_sequences.append(''.join(mutated_seq))

        df.at[index, 'variant_seq'] = mutated_sequences

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

import numpy as np

def action_ranking(t, df):
    # Ensure the 'action_score' column exists
    if 'action_score' not in df.columns:
        df['action_score'] = None

    # Iterate over rows where 't' column value is t
    for index, row in df[df['t'] == t].iterrows():
        variant_seqs = row['variant_seq']

        # Generate a list of action scores
        action_scores = [-np.log(np.random.uniform(0, 1)) for _ in range(len(variant_seqs))]

        # Update the 'action_score' column with the list of scores
        df.at[index, 'action_score'] = action_scores

    return df


class Agent:

    def __init__(self, t, df, reward, cfg):
        self.t = t
        self.df = df
        self.reward = reward
        self.policy_flag = cfg.params.basic_settings.policy_flag
        self.cfg = cfg

    def apply_policy(self):   


        if self.policy_flag == 'delete_and_mutate_ESM':

            # read seq from data frame
            if self.t == 1:  # Adjust formatting of input data frame in the first iteration

                # renaming sequence_number column to original_sequence and write the values
                self.df.rename(columns={'sequence_number': 'original_seq'}, inplace=True)
                self.df['original_seq'] = self.df['seq']

                # Inserting an empty column named 'shortened_seq'
                self.df.insert(2, 'shortened_seq', '')

                # Renaming the column 'seq' to 'variant_seq' and converting values to lists
                self.df['seq'] = self.df['seq'].apply(lambda x: [x])
                self.df.rename(columns={'seq': 'variant_seq'}, inplace=True)

                # Adding a column 'seed_flag' after 'variant_seq' and setting its value to a list containing True
                variant_seq_index = self.df.columns.get_loc('variant_seq')
                self.df.insert(variant_seq_index + 1, 'seed_flag', [[True]] * len(self.df))

            # perform exhaustive deletion and return a list of shortened_sequences
            df = exhaustive_deletion(self.t, self.df)

            df = mutate_single_residue(self.t, df) # to be replaced by greedy_sampling

            df = action_constraint(self.t, df)

            ## action ranking
            df = action_ranking(self.t, df)
        
            return df, pd.DataFrame.empty
        
        elif self.policy_flag == 'random_mutation':

            # select subset of sequences (those of the current step) on which to which to apply the policy
            if self.t == 0:
                df_action = self.df[self.df['t'] == 0][['t', 'sequence_number', 'seq']].copy()
            else:
                df_action = self.df[self.df['t'] == self.t - 1][['t', 'sequence_number', 'seq']].copy()    

            amino_acids = 'ACDEFGHIKLMNPQRSTVWY'  # List of amino acids in one-letter code   

            # Update the 't' column in df_action to the current value of self.t
            df_action['t'] = self.t

            for index, row in df_action.iterrows():
                seq = list(row['seq'])  # Convert string to list for mutation

                if len(seq) > 0:
                    mutation_pos = random.randint(0, len(seq) - 1)  # Select random position
                    original_residue = seq[mutation_pos]

                    # Select a new amino acid different from the original
                    new_residue = random.choice([aa for aa in amino_acids if aa != original_residue])

                    # Perform the mutation
                    seq[mutation_pos] = new_residue

                    # Update the sequence in the DataFrame without brackets
                    df_action.at[index, 'seq'] = ''.join(seq)

                    print('df_action', df_action)

                    # Optional: Print the original and mutated sequences with highlighted mutations
                    original_seq_for_printing = ''.join(seq[:mutation_pos]) + '[' + original_residue + ']' + ''.join(seq[mutation_pos + 1:])
                    mutated_seq_for_printing = ''.join(seq[:mutation_pos]) + '[' + new_residue + ']' + ''.join(seq[mutation_pos + 1:])
                    print(f"Original: {original_seq_for_printing} -> Mutated: {mutated_seq_for_printing}")

            # Concatenate self.df with df_action
            result_df = pd.concat([self.df, df_action])

            return result_df, df_action
        
        elif self.policy_flag == 'random_deletion':

            # select subset of sequences (those of the current step) on which to which to apply the policy
            if self.t == 0:
                df_action = self.df[self.df['t'] == 0][['t', 'sequence_number', 'seq']].copy()
            else:
                df_action = self.df[self.df['t'] == self.t - 1][['t', 'sequence_number', 'seq']].copy()   

            for index, row in df_action.iterrows():
                seq = list(row['seq'])  # Convert string to list for deletion

                if len(seq) > 1:  # Ensure there is at least one residue to delete
                    deletion_pos = random.randint(0, len(seq) - 1)  # Select random position

                    original_residue = seq[deletion_pos]
                    original_seq_for_printing = ''.join(seq[:deletion_pos]) + '[' + original_residue + ']' + ''.join(seq[deletion_pos + 1:])
                    seq_after_delete_for_printing = ''.join(seq[:deletion_pos]) + ''.join(seq[deletion_pos + 1:])

                    # Perform the deletion
                    del seq[deletion_pos]

                    # Update the sequence in the DataFrame after deletion
                    df_action.at[index, 'seq'] = ''.join(seq)

                    # Optional: Print the sequence before and after deletion
                    print(f"Sequence before deletion: {original_seq_for_printing} -> Sequence after deletion: {seq_after_delete_for_printing}")

            # Update the 't' column in df_action to the current value of self.t
            df_action['t'] = self.t

            # Concatenate self.df with df_action
            result_df = pd.concat([self.df, df_action])

            return result_df, df_action
        

### some old code snippets
            # ## some pseudo code:
            # # replacing landscape, alphabet, starting_sequence = get_landscape(args):
            # alphabet = 'ACDEFGHIKLMNPQRSTVWY'
            # starting_sequence = row['seq']
            # # and landscape is now given by the plddt from 

            # print('landscape', landscape)
            # print('alphabet', alphabet)
            # print('starting_seq', starting_sequence)
            # model = get_model(args, alphabet=alphabet, starting_sequence=starting_sequence)
            # explorer = get_algorithm(args, model=model, alphabet=alphabet, starting_sequence=starting_sequence)

            # runner = Runner(args)
            # runner.run(landscape, starting_sequence, model, explorer)   


                # if self.policy_flag == 'proximal_exploration':
        #     # keep the oracle/environment and put random smapling of sequences into it; always keep the total number of sequences that are fed into the agent fixed (think seeds)
        #     # Implement the ESM in the argent 

        #     alphabet = 'ACDEFGHIKLMNPQRSTVWY'  # alphabet of amino acids in one letter code   

        #     # load the wild-type sequence and its score f(wt) from the first row (zeroth iteration entry) of df.
        #     dt_wt = 

        #     # select subset of sequences (those of the current step) on which to which to apply the policy
        #     if self.t == 0:
        #         df_action = self.df[self.df['t'] == 0][['t', 'sequence_number', 'seq']].copy()
        #     else:
        #         df_action = self.df[self.df['t'] == self.t - 1][['t', 'sequence_number', 'seq']].copy()    

        #     # Update the 't' column in df_action to the current value of self.t
        #     df_action['t'] = self.t

        #     for index, row in df_action.iterrows():
        #         seq = list(row['seq'])  # Convert string to list for mutation

        #         if len(seq) > 0:
        #             mutation_pos = random.randint(0, len(seq) - 1)  # Select random position
        #             original_residue = seq[mutation_pos]

        #             # Select a new amino acid different from the original
        #             new_residue = random.choice([aa for aa in amino_acids if aa != original_residue])

        #             # Perform the mutation
        #             seq[mutation_pos] = new_residue

        #             # Update the sequence in the DataFrame without brackets
        #             df_action.at[index, 'seq'] = ''.join(seq)

        #             print('df_action', df_action)

        #             # Optional: Print the original and mutated sequences with highlighted mutations
        #             original_seq_for_printing = ''.join(seq[:mutation_pos]) + '[' + original_residue + ']' + ''.join(seq[mutation_pos + 1:])
        #             mutated_seq_for_printing = ''.join(seq[:mutation_pos]) + '[' + new_residue + ']' + ''.join(seq[mutation_pos + 1:])
        #             print(f"Original: {original_seq_for_printing} -> Mutated: {mutated_seq_for_printing}")

        #     # Concatenate self.df with df_action
        #     result_df = pd.concat([self.df, df_action])

        #     return result_df, df_action           

        #     for index, row in df_action.iterrows():

        #         args = get_args()

        #         # # Some pseudo code / this is def run() function from the evals_utils script adapted and put in the context of our code.
        #         # D^{t-1} is given by the entire (!) df table

        #         # train the model f^hat on D^{t-1}

        #         # run pex with f^hat and D^{t-1}, this gives a new set of M sequences, s_i^(t). Probably store them in df_action.

        #     # # Concatenate self.df with df_action
        #     # result_df = pd.concat([self.df, df_action])
        
        #     # next step is to return results_df and the oracle is computing the truth f(s_i^(t)) for us and supplement the full df table with it.
        #     return result_df, df_action     

# def exhaustive_deletion(t, df):
#     # Iterate over rows where 't' column value is t-1 and 'seed' is True
#     for index, row in df[(df['t'] == t-1) & (df['seed'] == True)].iterrows():
#         variant_seq = row['variant_seq']
#         variant_seq_list = list(variant_seq)

#         # Iterate over the length of the variant_sequence
#         for n in range(len(variant_seq_list)):
#             # Create a new sequence excluding the character at position n
#             new_sequence = ''.join(variant_seq_list[:n] + variant_seq_list[n+1:])
#             # Append a new row to the data frame
#             new_row = pd.DataFrame({'t': [t], 'shortened_seq': [new_sequence]})
#             df = pd.concat([df, new_row], ignore_index=True)

#     return df

            # if self.t == 1: # adjust formatting of input data frame in the first iteration
            #     # Inserting an empty column named 'shortened_seq'
            #     self.df.insert(2, 'shortened_seq', '')

            #     # Renaming the column 'seq' to 'variant_seq'
            #     self.df.rename(columns={'seq': 'variant_seq'}, inplace=True)

            #     # Adding a column 'seed' after 'variant_seq' and setting its value to True
            #     variant_seq_index = self.df.columns.get_loc('variant_seq')
            #     # self.df.insert(variant_seq_index + 1, 'seed', True)
            #     self.df.insert(variant_seq_index + 1, 'seed', [[True]] * len(self.df))