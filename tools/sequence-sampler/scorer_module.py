import subprocess
import os
import pandas as pd
import sequence_transformer
from AF2_module import AF2Runner
from utils import squeeze_seq
from utils import write_af2_update
from utils import compute_affinity
from utils import concatenate_to_df
import logging

class StateScorer:
    def __init__(self, evo_cycle, sequence, cfg, outputs_directory, df):
        self.evo_cycle = evo_cycle #- 1
        self.scorer_list = cfg.params.basic_settings.scorers.split(',')
        self.sequence = squeeze_seq(sequence)
        self.outputs_directory = outputs_directory
        self.cfg = cfg
        self.df = df

    def run_scoring(self):

        print('\n')
        print("Running scoring job...")
        df_score = pd.DataFrame() # initialize data frame
        for scorer in self.scorer_list:

            scorer_directory = os.path.join(self.outputs_directory, scorer)
            if not os.path.exists(scorer_directory):
                os.makedirs(scorer_directory, exist_ok=True)

            logging.info(f"Running {scorer}")
            if scorer=='ESM2':
                runner = sequence_transformer.ESM2Runner() # initialize ESM2Runner with the default model
                LLmatrix_sequence = runner.token_masked_marginal_log_likelihood_matrix(self.sequence)

                LL_mod = compute_log_likelihood(sequence, LLmatrix_mod) # TD: normalization by sequence length?

                if 'pseudolikelihood' not in df_score.columns:
                    df_score['pseudolikelihood'] = None  # Initialize the column with None
                if t==1:
                    df_score.at[0, 'pseudolikelihood'] = LL_mod
                else:
                    df_score.iloc[-1, df_score.columns.get_loc('pseudolikelihood')] = LL_mod
            
            elif scorer=='Colabfold':
                target_binder_sequence = f"{self.cfg.params.basic_settings.target_seq}:{self.sequence}" # TD: fix this; maybe load the target-sequence into the cfg from the pdb or fasta
                
                # include a function that combines binder and target sequence
                input_dir = os.path.join(self.cfg.inputs.directory, 'current_sequences')
                if os.path.exists(input_dir):
                    # If the folder already exists, empty the folder of all files
                    for file_name in os.listdir(input_dir):
                        file_path = os.path.join(input_dir, file_name)
                        if os.path.isfile(file_path):
                            os.remove(file_path)
                else:
                    os.makedirs(input_dir, exist_ok=True)

                file_path = os.path.join(input_dir, f"evo_cycle_{self.evo_cycle}.fasta")
                with open(file_path, 'w') as file:
                    file.write(f">evo_cycle_{self.evo_cycle}\n{target_binder_sequence}\n")

                seq_input_dir = os.path.abspath(input_dir)

                af2_runner = AF2Runner(seq_input_dir, scorer_directory)
                af2_runner.run()

                # append output as new columns of data frame
                df_score = write_af2_update(df_score, scorer_directory, json_pattern=f"evo_cycle_{self.evo_cycle}")
            
            elif scorer=='Prodigy': # not implemented yet

                pdb_file_path = df_score['absolute pdb path'].iloc[0]
                affinity = compute_affinity(pdb_file_path)
                if 'affinity' not in df_score.columns:
                    df_score['affinity'] = None  # Initialize the column with None

                # Assuming you have a single row, set the value of 'affinity' for the first row
                df_score.at[0, 'affinity'] = affinity

                if affinity is not None:
                    print(f"The affinity for the complex {pdb_file_path} is {affinity}")

            elif scorer=='Hamming':

                # Function to compute Hamming distance
                def compute_hamming_distance(seq1, seq2):
                    return sum(c1 != c2 for c1, c2 in zip(seq1, seq2))

                # Filter rows where 't' column is 0
                filtered_df = self.df[self.df['t'] == 0]

                # Compute Hamming distances for filtered rows
                hamming_distances = filtered_df['seed'].apply(lambda x: compute_hamming_distance(self.sequence, x))

                # Calculate the mean of the Hamming distances
                mean_hamming_distance = hamming_distances.mean()

                # Add the mean Hamming distance to df_score
                if 'mean_hamming_distance_to_init_seqs' not in df_score.columns:
                    df_score['mean_hamming_distance_to_init_seqs'] = None  # Initialize the column with None

                # Set the value of 'hamming_distance' for the first row
                df_score.at[0, 'mean_hamming_distance_to_init_seqs'] = mean_hamming_distance

                if mean_hamming_distance is not None:
                    logging.info(f"The mean Hamming distance for the selected sequences is {mean_hamming_distance}")

            #     df_score.to_csv(f"{scorer_directory}/output.csv", index=False)
            
            df_score.to_csv(f"{scorer_directory}/output.csv", index=False)
        
        logging.info(f"Scoring job complete. Results are in {self.outputs_directory}")

        # supplement data frame by scores
        df = concatenate_to_df(self.evo_cycle, df_score, self.df)

        return df
        # return df_score, LLmatrix_sequence

    def run(self):
        return self.run_scoring()