import os
from AF2_module import AF2Runner
import os
import json
import random
import pandas as pd
import numpy as np
import glob
import sequence_transformer

def write_dataframe_to_fastas(t, dataframe, cfg):
    input_dir = os.path.join(cfg.inputs.directory, 'current_sequences')
    if os.path.exists(input_dir):
        # If the folder already exists, empty the folder of all files
        for file_name in os.listdir(input_dir):
            file_path = os.path.join(input_dir, file_name)
            if os.path.isfile(file_path):
                os.remove(file_path)
    else:
        os.makedirs(input_dir, exist_ok=True)

    for index, row in dataframe.iterrows():
        file_path = os.path.join(input_dir, f"seq_{row['sequence_number']}_t_{t}.fasta")
        with open(file_path, 'w') as file:
            file.write(f">{row['sequence_number']}\n{row['seq']}\n")
    return os.path.abspath(input_dir)

def supplement_dataframe(t, df, directory_path):
    # Ensure DataFrame has a proper index
    df.reset_index(drop=True, inplace=True)

    # Adding new columns only if they don't already exist
    for col in ['top rank json', 'top rank pdb', 'plddt', 'max_pae', 'ptm']:
        if col not in df.columns:
            df[col] = np.nan

    for index, row in df.iterrows():
        if row['t'] == t:
            sequence_number = row['sequence_number']

            # Search for the matching json and pdb files
            json_pattern = os.path.join(directory_path, f'seq_{sequence_number}_t_{t}_scores_rank_001*.json')
            pdb_pattern = os.path.join(directory_path, f'seq_{sequence_number}_t_{t}_unrelaxed_rank_001*.pdb')

            json_files = glob.glob(json_pattern)
            pdb_files = glob.glob(pdb_pattern)

            # Assuming the first match is the desired file
            if json_files:
                df.at[index, 'top rank json'] = json_files[0]
                with open(json_files[0], 'r') as file:
                    data = json.load(file)
                    df.at[index, 'plddt'] = np.mean(data['plddt'])
                    df.at[index, 'max_pae'] = data['max_pae']
                    df.at[index, 'ptm'] = data['ptm']

            if pdb_files:
                df.at[index, 'top rank pdb'] = pdb_files[0]

    return df

def modified_sampling_set_pseudoLLSelection(t, df, cfg):
    N = cfg.params.basic_settings.max_number_of_offspring_kept # number of same-seed offspring kept
    k = cfg.params.basic_settings.k  # max number of samples

    runner = sequence_transformer.ESM2Runner()

    # Initialize 'variant_pseudoLL' column if it doesn't exist
    if 'variant_pseudoLL' not in df.columns:
        df['variant_pseudoLL'] = df.apply(lambda row: [np.nan] * len(row['variant_seq']), axis=1)

    # Iterate over rows where 't' column value is t
    for index, row in df[df['t'] == t].iterrows():
        action_scores = row['action_score']
        variant_list = row['variant_seq']
        seed_flags = row['seed_flag']

        # Determine the list of distinct strings in variant_list
        distinct_variants = list(set(variant_list))

        # Initialize seed_flags with all False
        seed_flags = [False] * len(variant_list)

        # Loop over variant_list and set one element in seed_flags to True for each distinct element
        for i, variant in enumerate(variant_list):
            if variant in distinct_variants:
                seed_flags[i] = True
                distinct_variants.remove(variant)

        # Filter and sort action scores based on seed_flags
        filtered_scores = [(score, i) for i, (score, flag) in enumerate(zip(action_scores, seed_flags)) if flag]
        filtered_scores.sort(reverse=True, key=lambda x: x[0])

        # Select top N elements
        top_indices = {index for _, index in filtered_scores[:N]}

        # Update seed_flags: True if in top N, otherwise False
        updated_seed_flags = [i in top_indices for i, _ in enumerate(seed_flags)]
        df.at[index, 'seed_flag'] = updated_seed_flags

        variant_pseudoLL = [np.nan] * len(variant_list)  # Initialize with np.nan

        # Compute pseudoLL for elements with True in updated_seed_flags
        for i, (variant, flag) in enumerate(zip(variant_list, updated_seed_flags)):
            if flag:
                pseudoLL = runner.sequence_pseudo_log_likelihoods_scalar(variant)
                if isinstance(pseudoLL, (int, float)):  # Ensure pseudoLL is numeric
                    variant_pseudoLL[i] = pseudoLL

        # Update the DataFrame
        df.at[index, 'variant_pseudoLL'] = variant_pseudoLL

    # Flatten all pseudoLL values across the relevant rows, ignoring NaNs
    t_rows = df[df['t'] == t]
    all_pseudoLLs = pd.Series([item for sublist in t_rows['variant_pseudoLL'].tolist() for item in sublist])
    all_pseudoLLs = all_pseudoLLs.dropna()

    # Find the k largest variant_pseudoLL values
    top_k_values = all_pseudoLLs.nlargest(k)

    # Extract the variant sequences of with the k-largest pseudoLL values
    top_k_variants = []
    for index, row in t_rows.iterrows():
        for idx, val in enumerate(row['variant_pseudoLL']):
            if val in top_k_values.values:
                variant = df.at[index, 'variant_seq'][idx]
                top_k_variants.append(variant)

    # Compute the list of unique top variants
    unique_top_variants = list(set(top_k_variants))
    print('Unique top-k variants:', unique_top_variants)

    # Update seed_flags based on the unique top variants
    for index, row in t_rows.iterrows():
        seed_flags = row['seed_flag']
        for idx, variant in enumerate(row['variant_seq']):
            if variant in unique_top_variants:
                seed_flags[idx] = True
                unique_top_variants.remove(variant)  # Remove the variant as it's already used
            else:
                seed_flags[idx] = False  # Ensure other flags are set to False
        df.at[index, 'seed_flag'] = seed_flags

    # Print the corresponding variant_seqs along with their seed_flag and pseudoLL
    print("\nTop", k, "variants:")
    for index, row in t_rows.iterrows():
        for variant, flag, pseudoLL in zip(row['variant_seq'], row['seed_flag'], row['variant_pseudoLL']):
            if pseudoLL in top_k_values.values:
                print(f"Variant: {variant}, Seed Flag: {flag}, pseudoLL: {pseudoLL}")


    return df


def modified_sampling_set_pseudoLL(t, df, cfg):
    N = cfg.params.basic_settings.max_number_of_offspring_kept

    runner = sequence_transformer.ESM2Runner()

    # Check if the 'variant_pseudoLL' column exists in df, if not, create it
    if 'variant_pseudoLL' not in df.columns:
        df['variant_pseudoLL'] = None

    # Iterate over rows where 't' column value is t
    for index, row in df[df['t'] == t].iterrows():
        action_scores = row['action_score']
        variant_list = row['variant_seq']
        seed_flags = row['seed_flag']

        # Determine the list of distinct strings in variant_list
        distinct_variants = list(set(variant_list))

        # Initialize seed_flags with all False
        seed_flags = [False] * len(variant_list)

        # Loop over variant_list and set one element in seed_flags to True for each distinct element
        for i, variant in enumerate(variant_list):
            if variant in distinct_variants:
                seed_flags[i] = True
                distinct_variants.remove(variant)

        # Filter and sort action scores based on seed_flags
        filtered_scores = [(score, i) for i, (score, flag) in enumerate(zip(action_scores, seed_flags)) if flag]
        filtered_scores.sort(reverse=True, key=lambda x: x[0])

        # Select top N elements
        top_indices = {index for _, index in filtered_scores[:N]}

        # Update seed_flags: True if in top N, otherwise False
        updated_seed_flags = [i in top_indices for i, _ in enumerate(seed_flags)]
        df.at[index, 'seed_flag'] = updated_seed_flags

        # Initialize the pseudoLL list with None
        variant_pseudoLL = [None] * len(variant_list)

        # Compute pseudoLL for elements with True in updated_seed_flags
        for i, (variant, flag) in enumerate(zip(variant_list, updated_seed_flags)):
            if flag:
                pseudoLL = runner.sequence_pseudo_log_likelihoods_scalar(variant)
                variant_pseudoLL[i] = pseudoLL

        # Update the DataFrame
        df.at[index, 'variant_pseudoLL'] = variant_pseudoLL

    return df

def modified_sampling_set(t, df, cfg):
    N = cfg.params.basic_settings.max_number_of_offspring_kept

    # Iterate over rows where 't' column value is t
    for index, row in df[df['t'] == t].iterrows():
        action_scores = row['action_score']
        variant_list = row['variant_seq']
        seed_flags = row['seed_flag']

        # Determine the list of distinct strings in variant_list
        distinct_variants = list(set(variant_list))

        # Initialize seed_flags with all False
        seed_flags = [False] * len(variant_list)

        # Loop over variant_list and set one element in seed_flags to True for each distinct element
        for i, variant in enumerate(variant_list):
            if variant in distinct_variants:
                seed_flags[i] = True
                distinct_variants.remove(variant)

        # Filter and sort action scores based on seed_flags
        filtered_scores = [(score, i) for i, (score, flag) in enumerate(zip(action_scores, seed_flags)) if flag]
        filtered_scores.sort(reverse=True, key=lambda x: x[0])

        # Select top N elements
        top_indices = {index for _, index in filtered_scores[:N]}

        # Update seed_flags: True if in top N, otherwise False
        updated_seed_flags = [i in top_indices for i, _ in enumerate(seed_flags)]
        df.at[index, 'seed_flag'] = updated_seed_flags

    return df

def action_selection(t, df, cfg):

    if t>0:
        # df = modified_sampling_set(t, df, cfg)
        # df = modified_sampling_set_pseudoLL(t, df, cfg)
        df = modified_sampling_set_pseudoLLSelection(t, df, cfg)

    return df

def print_rows_with_t(t, df):
    """
    This function prints all rows of the DataFrame 'df' 
    where the 't'-column has the value t.
    """
    # Filter the DataFrame based on the condition
    filtered_df = df[df['t'] == t]
    
    # Print the filtered DataFrame
    print(filtered_df)

class Oracle:
    def __init__(self, t, df, df_action, outputs_directory, cfg):

        self.t = t
        self.df = df
        self.df_action = df_action
        self.outputs_directory = outputs_directory
        self.cfg = cfg

    def run(self):

        ### likelihood based evolution ### 
        df = action_selection(self.t, self.df, self.cfg)
        print_rows_with_t(self.t, df)

        ### AF2 Runner ### 
        # # prepare input sequences as fastas and run AF2
        # seq_input_dir = write_dataframe_to_fastas(self.t, self.df_action, self.cfg)

        if self.t==...:

            # loop over the sequences which have seed_flag==True at time t, clear the seqs_to_fold directory in the outputs directory, and then write fastas for each most likely sequence, and run the AF2. Then supplement the data frame with metrics and path the pdbs.
            # write them to a separate csv in which contains only the most likely sequences use is an id to the identify also the AF2 results folder.  

            print("starting repeat number ", n)
            af2_runner = AF2Runner(seq_input_dir, self.outputs_directory)
            af2_runner.run()

            # # complete df data frame with info
            supplemented_dataframe = supplement_dataframe(self.t, self.df, self.outputs_directory)

        return df
    

### OLD CODE SNIPPETS ###
# def sampling_set(t, df, cfg):
#     k = cfg.params.basic_settings.k  # max number of samples

#     # Iterate over rows where 't' column value is t
#     for index, row in df[df['t'] == t].iterrows():
#         action_ranking = row['action_score']
#         variant_list = row['variant_seq']
#         length_of_ranking = len(action_ranking)

#         # Determine the list of distinct strings in variant_list
#         distinct_variants = list(set(variant_list))

#         # Initialize seed_flags with all False
#         seed_flags = [False] * len(variant_list)

#         # Loop over variant_list and set one element in seed_flags to True for each distinct element
#         for i, variant in enumerate(variant_list):
#             if variant in distinct_variants:
#                 seed_flags[i] = True
#                 distinct_variants.remove(variant)

#         df.at[index, 'seed_flag'] = seed_flags

#         # # Handling action_ranking as in your original function
#         # if k < length_of_ranking:
#         #     sampled_indices = random.sample(range(length_of_ranking), k)
#         #     action_seed_flags = [index in sampled_indices for index in range(length_of_ranking)]
#         #     # Combine the two seed flags
#         #     combined_seed_flags = [sf and asf for sf, asf in zip(seed_flags, action_seed_flags)]
#         #     df.at[index, 'seed_flag'] = combined_seed_flags
#         # else:
#         #     df.at[index, 'seed_flag'] = seed_flags

#     return df
