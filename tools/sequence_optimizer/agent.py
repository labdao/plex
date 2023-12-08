import pandas as pd
import random

class Agent:

    def __init__(self, t, df, reward, policy_flag):
        self.t = t
        self.df = df
        self.reward = reward
        self.policy_flag = policy_flag

    def apply_policy(self):
        if self.policy_flag == 'random_mutation':
            # Filter df_action based on the value of 't'
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
            # Filter df_action based on the value of 't'
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
        



# Example usage:
# agent = Agent(df, reward, policy_flag, t)
# new_df = agent.policy()

# import pandas as pd
# import random

# class Agent:

#     def __init__(self, df, reward, policy_flag, t):
#         self.df = df
#         self.reward = reward
#         self.policy_flag = policy_flag
#         self.t = t

#     def policy(self):
#         if self.policy_flag == 'random_mutation':
#             df_action = self.df[['sequence_number', 'seq']].copy()
#             amino_acids = 'ACDEFGHIKLMNPQRSTVWY'  # List of amino acids in one-letter code

#             for index, row in df_action.iterrows():
#                 seq = list(row['seq'])  # Convert string to list for mutation

#                 if len(seq) > 0:
#                     mutation_pos = random.randint(0, len(seq) - 1)  # Select random position
#                     original_residue = seq[mutation_pos]

#                     # Select a new amino acid different from the original
#                     new_residue = random.choice([aa for aa in amino_acids if aa != original_residue])

#                     # Perform the mutation
#                     seq[mutation_pos] = new_residue

#                     # Create strings for printing with brackets around the mutated position
#                     original_seq_for_printing = ''.join(seq[:mutation_pos]) + '[' + original_residue + ']' + ''.join(seq[mutation_pos + 1:])
#                     mutated_seq_for_printing = ''.join(seq[:mutation_pos]) + '[' + new_residue + ']' + ''.join(seq[mutation_pos + 1:])

#                     # Update the sequence in the DataFrame without brackets
#                     df_action.at[index, 'seq'] = ''.join(seq)

#                     # Print the original and mutated sequences with highlighted mutations
#                     print(f"Original: {original_seq_for_printing} -> Mutated: {mutated_seq_for_printing}")

#             return df_action
#         # elif  self.policy_flag == 'random_delete':
#         #     ...

