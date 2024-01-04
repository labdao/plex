import os
from AF2_module import AF2Runner
import json
import random
import pandas as pd
import numpy as np
import glob
import sequence_transformer
import shutil
import subprocess

def print_rows_with_t(t, df):
    """
    This function prints all rows of the DataFrame 'df' 
    where the 't'-column has the value t.
    """
    # Filter the DataFrame based on the condition
    filtered_df = df[df['t'] == t]
    
    # Print the filtered DataFrame
    print(filtered_df)

def create_complex_df(t, MML_df, outputs_directory, cfg): # create a data frame for the protein complex
    # Filter rows from MML_df where 't' column value is t
    filtered_df = MML_df[MML_df['t'] == t]

    # Create a new DataFrame with specified columns
    complex_df = pd.DataFrame(columns=['t', 'target_seq', 'binder_seq', 'complex', 'sequence_pseudo_LL', 'mean plddt', 'max pae', 'json', 'pdb'])

    # Extract target_seq from cfg
    target_seq = cfg.params.basic_settings.target_seq

    ### FOR THE ORIGINAL SEQUENCE ###
    # Assuming binder_seq is the sequence you want from MML_df for t=0
    binder_seq = filtered_df.iloc[0]['original_seq']  # Change 'variant_seq' to the correct column name if different

    # Create the 'complex' value as concatenation of target_seq and binder_seq
    complex_value = f"{target_seq}:{binder_seq}"

    # Construct the row with t=0 as a DataFrame
    row_df = pd.DataFrame([{
        't': 0,
        'target_seq': target_seq,
        'binder_seq': binder_seq,
        'complex': complex_value,
        # Leaving other columns blank
    }])

    # Concatenate new row DataFrame to complex_df
    complex_df = pd.concat([complex_df, row_df], ignore_index=True)

    # Write complex_df to a CSV file
    complex_df.to_csv('folding_with_target_summary.csv', index=False)

    # Extract corresponding string sequence
    sequence = complex_value

    folder = 'sequence_to_fold'
    if os.path.exists(folder):
        shutil.rmtree(folder)
    os.makedirs(folder, exist_ok=True)

    # Write into a fasta file in sequence_to_fold directory
    fasta_file_path = os.path.join(outputs_directory, folder)
    print('fasta_file_path', fasta_file_path)
    fasta_filename = os.path.join(folder, f"complex_t0.fasta")
    with open(fasta_filename, 'w') as fasta_file:
        fasta_file.write(f">1\n{sequence}\n") 

    print("starting AF2 MML sequence")

    af2_runner = AF2Runner(folder, outputs_directory)
    af2_runner.run()

    print("done folding")

    json_pattern = f"complex_t0"
    complex_df = update_complex_summary(0, complex_df, sequence, None, complex_df, outputs_directory, json_pattern)

    ### FOR THE FINAL SHORTENED SEQUENCE ###
    # Create a new DataFrame with specified columns
    complex_df = pd.DataFrame(columns=['t', 'target_seq', 'binder_seq', 'complex', 'sequence_pseudo_LL', 'mean plddt', 'max pae', 'json', 'pdb'])

    # Assuming binder_seq is the sequence you want from MML_df for t=0
    binder_seq = filtered_df.iloc[0]['variant_seq']  # Change 'variant_seq' to the correct column name if different

    # Create the 'complex' value as concatenation of target_seq and binder_seq
    complex_value = f"{target_seq}:{binder_seq}"

    # Construct the row with t=0 as a DataFrame
    row_df = pd.DataFrame([{
        't': cfg.params.basic_settings.number_of_evo_cycles,
        'target_seq': target_seq,
        'binder_seq': binder_seq,
        'complex': complex_value,
        # Leaving other columns blank
    }])

    # Concatenate new row DataFrame to complex_df
    complex_df = pd.concat([complex_df, row_df], ignore_index=True)

    # Write complex_df to a CSV file
    complex_df.to_csv('folding_with_target_summary.csv', index=False)

    # Extract corresponding string sequence
    sequence = complex_value

    folder = 'sequence_to_fold'
    if os.path.exists(folder):
        shutil.rmtree(folder)
    os.makedirs(folder, exist_ok=True)

    # Write into a fasta file in sequence_to_fold directory
    fasta_file_path = os.path.join(outputs_directory, folder)
    print('fasta_file_path', fasta_file_path)
    fasta_filename = os.path.join(folder, f"complex_tFINAL.fasta")
    with open(fasta_filename, 'w') as fasta_file:
        fasta_file.write(f">1\n{sequence}\n") 

    print("starting AF2 MML sequence")

    af2_runner = AF2Runner(folder, outputs_directory)
    af2_runner.run()

    print("done folding")

    json_pattern = f"complex_tFINAL"
    complex_df = update_complex_summary(cfg.params.basic_settings.number_of_evo_cycles, complex_df, sequence, None, complex_df, outputs_directory, json_pattern)

    return complex_df

# def prodigy_run(csv_path, pdb_path):
def prodigy_run(csv_path):
    # print('csv path', csv_path)
    df = pd.read_csv(csv_path)
    # print('csv file data frame', df)
    for i, r in df.iterrows():

        file_path = r['pdb']
        if pd.notna(file_path):
            print('file path', file_path)
            try:
                subprocess.run(
                    ["prodigy", "-q", file_path], stdout=open("temp.txt", "w"), check=True
                )
                with open("temp.txt", "r") as f:
                    lines = f.readlines()
                    if lines:  # Check if lines is not empty
                        affinity = float(lines[0].split(" ")[-1].split("/")[0])
                        df.loc[i, "affinity"] = affinity
                    else:
                        # print(f"No output from prodigy for {r['path']}")
                        print(f"No output from prodigy for {file_path}")
                        # Handle the case where prodigy did not produce output
            except subprocess.CalledProcessError:
                # print(f"Prodigy command failed for {r['path']}")
                print(f"Prodigy command failed for {file_path}")

    # export results
    df.to_csv(f"{csv_path}", index=None)

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

def update_summary(t, row, sequence, sequence_pseudoLL, df, directory, json_pattern, index, i):
    summary_file = os.path.join(directory, 'MML_summary.csv')

    # Loop over JSON files that match the given pattern for the current iteration
    for json_file in glob.glob(os.path.join(directory, json_pattern)):
        with open(json_file, 'r') as file:
            data = json.load(file)

        print('data', data)

        # Compute average plddt
        avg_plddt = sum(data['plddt']) / len(data['plddt'])

        # Get max_pae value
        max_pae = data['max_pae']

        # Find corresponding PDB file
        json_pattern = f"sequence_Time{t}_TablRow{index}_VariantIdx{i}_unrelaxed"
        pdb_file = None
        rank_str = json_file.split('rank')[1].split('.')[0]
        for pdb in glob.glob(os.path.join(directory, f'{json_pattern}*rank{rank_str}.pdb')): # maybe add json pattern here, like is already done for the protein complex routine
            pdb_file = pdb
            break

        # Prepare new row
        new_row = {
            't': t,
            'original_seq': row['original_seq'],
            'variant_seq': sequence,
            'sequence_pseudo_LL': sequence_pseudoLL,
            'mean plddt': avg_plddt,
            'max pae': max_pae,
            'json': os.path.abspath(json_file),
            'pdb': os.path.abspath(pdb_file) if pdb_file else None
        }

        # Concatenate new row to DataFrame
        df = pd.concat([df, pd.DataFrame([new_row])], ignore_index=True)

    if not df.empty:
        if os.path.exists(summary_file):
            df.to_csv(summary_file, mode='a', header=False, index=False)
        else:
            df.to_csv(summary_file, index=False)

    return df

def update_complex_summary(t, row, sequence, sequence_pseudoLL, df, directory, json_pattern):
    summary_file = os.path.join(directory, 'folding_with_target_summary.csv')

    # Loop over JSON files that match the given pattern for the current iteration
    for json_file in glob.glob(os.path.join(directory, f"{json_pattern}_scores*.json")):
        with open(json_file, 'r') as file:
            data = json.load(file)

        # Compute average plddt
        avg_plddt = sum(data['plddt']) / len(data['plddt'])

        # Get max_pae value
        max_pae = data['max_pae']

        # Find corresponding PDB file
        pdb_file = None
        rank_str = json_file.split('rank')[1].split('.')[0]
        for pdb in glob.glob(os.path.join(directory, f"{json_pattern}_unrelaxed_rank{rank_str}*.pdb")):
            pdb_file = pdb
            break

        # Prepare new row
        new_row = {
            't': t,
            'target_seq': row['target_seq'].iloc[0],
            'binder_seq': row['binder_seq'].iloc[0],
            'complex': sequence,
            'sequence_pseudo_LL': sequence_pseudoLL,
            'mean plddt': avg_plddt,
            'max pae': max_pae,
            'json': os.path.abspath(json_file),
            'pdb': os.path.abspath(pdb_file) if pdb_file else None
        }

        # Concatenate new row to DataFrame
        df = pd.concat([df, pd.DataFrame([new_row])], ignore_index=True)

    if not df.empty:
        if os.path.exists(summary_file):
            df.to_csv(summary_file, mode='a', header=False, index=False)
        else:
            df.to_csv(summary_file, index=False)

    return df

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

def permissibility_sampling_pseudoLLSelection(t, df, cfg):
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

def action_selection(t, df, cfg):

    if t>0:
        # df = modified_sampling_set_pseudoLLSelection(t, df, cfg)
        df = permissibility_sampling_pseudoLLSelection(t, df, cfg)

    return df

class Oracle:
    def __init__(self, t, df, df_action, outputs_directory, cfg):

        self.t = t
        self.df = df
        self.df_action = df_action
        self.outputs_directory = outputs_directory
        self.cfg = cfg

    def run(self):

        ### run likelihood-based evolution ### 
        df = action_selection(self.t, self.df, self.cfg)
        print_rows_with_t(self.t, df)

        ### run AF2 ### 
        # # prepare input sequences as fastas and run AF2
        # seq_input_dir = write_dataframe_to_fastas(self.t, self.df_action, self.cfg)

        MML_df = pd.DataFrame(columns=['t', 'original_seq', 'variant_seq', 'sequence_pseudo_LL',  'mean plddt', 'max pae', 'json', 'pdb'])

        # print('mml df', MML_df)

        # if self.t==self.cfg.params.basic_settings.number_of_evo_cycles:
        if self.t>=1:

            # Before the loop, clear the sequence_to_fold directory
            folder = 'sequence_to_fold'
            if os.path.exists(folder):
                shutil.rmtree(folder)
            os.makedirs(folder, exist_ok=True)

            for index, row in df[df['t'] == self.t].iterrows():
                variant_list = row['variant_seq']
                seed_flags = row['seed_flag']
                pseudo_LL = row['variant_pseudoLL']

                # Check if there are any True entries in seed_flags
                for i, flag in enumerate(seed_flags):
                    if flag:
                        # Extract corresponding string sequence
                        sequence = variant_list[i]
                        sequence_pseudoLL = pseudo_LL[i]

                        # Write into a fasta file in sequence_to_fold directory
                        fasta_file_path = os.path.join(self.outputs_directory, folder)
                        print('fasta_file_path', fasta_file_path)
                        fasta_filename = os.path.join(folder, f"sequence_Time{self.t}_TablRow{index}_VariantIdx{i}.fasta")
                        with open(fasta_filename, 'w') as fasta_file:
                            fasta_file.write(f">1\n{sequence}\n") 

                        print("starting AF2 MML sequence")

                        af2_runner = AF2Runner(folder, self.outputs_directory)
                        af2_runner.run()

                        print("done folding")

                        json_pattern = f"sequence_Time{self.t}_TablRow{index}_VariantIdx{i}_scores*.json"
                        MML_df = update_summary(self.t, row, sequence, sequence_pseudoLL, MML_df, self.outputs_directory, json_pattern, index, i)
        
        if self.t==self.cfg.params.basic_settings.number_of_evo_cycles:

            # print('mml df', MML_df)

            # fold original and target seqs together and save pdb
            # fold the final evolution step MML and target seqs toegther and save pdb
            complex_df = create_complex_df(self.t, MML_df, self.outputs_directory, self.cfg)
            print(complex_df)

            print("running Prodigy")
            prodigy_run(f"{self.outputs_directory}/folding_with_target_summary.csv")
            
        return df