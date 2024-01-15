import subprocess
import os
import pandas as pd
import sequence_transformer
from AF2_module import AF2Runner
from utils import squeeze_seq
from utils import write_af2_update

class StateScorer:
    def __init__(self, evo_cycle, scorer_list, sequence, cfg, outputs_directory):
        self.evo_cycle = evo_cycle
        self.scorer_list = scorer_list
        self.sequence = squeeze_seq(sequence)
        self.outputs_directory = outputs_directory
        self.cfg = cfg

    def run_scoring(self):

        print('\n')
        print("Running scoring job...")
        df_score = pd.DataFrame() # initialize data frame
        for scorer in self.scorer_list:

            scorer_directory = os.path.join(self.outputs_directory, scorer)
            if not os.path.exists(scorer_directory):
                os.makedirs(scorer_directory, exist_ok=True)

            if scorer=='ESM2':
                print(f"Running {scorer}")
                runner = sequence_transformer.ESM2Runner() # initialize ESM2Runner with the default model
                LLmatrix_sequence = runner.token_masked_marginal_log_likelihood_matrix(self.sequence)

                scores_to_add = {
                    'LLmatrix_sequence': LLmatrix_sequence #,
                    # 'LL_sequence': LL_sequence
                }
                for column_name, column_data in scores_to_add.items():
                    df_score[column_name] = pd.Series([column_data])
            
            elif scorer=='AF2':
                print(f"Running {scorer}")
                target_binder_sequence = f"{self.cfg.params.basic_settings.target_seq}:{self.sequence}"
                
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
                df_score = write_af2_update(df_score, scorer_directory, json_pattern=f"evo_cycle_{self.evo_cycle}\n{target_binder_sequence}\n")
                df_score.to_csv(f"{scorer_directory}/output.csv", index=False)
            
            # elif scorer=='Prodigy': # not implemented yet

                # Usage example
                # pdb_file_path = os.path.abspath(pdb_file)
                # affinity = compute_affinity(pdb_file_path)
                # if affinity is not None:
                #     print(f"The affinity for the file {pdb_file_path} is {affinity}")

                # def compute_affinity(file_path):
                # if pd.notna(file_path):
                #     try:
                #         # Run Prodigy and capture the output in temp.txt
                #         subprocess.run(
                #             ["prodigy", "-q", file_path], stdout=open("temp.txt", "w"), check=True
                #         )
                #         # Read the output from temp.txt
                #         with open("temp.txt", "r") as f:
                #             lines = f.readlines()
                #             if lines:  # Check if lines is not empty
                #                 # Extract the affinity value from the output
                #                 affinity = float(lines[0].split(" ")[-1].split("/")[0])
                #                 return affinity
                #             else:
                #                 print(f"No output from prodigy for {file_path}")
                #                 return None  # No output from Prodigy
                #     except subprocess.CalledProcessError:
                #         print(f"Prodigy command failed for {file_path}")
                #         return None  # Prodigy command failed
                # else:
                #     print("Invalid file path")
                #     return None  # Invalid file path provided
        
        print(f"Scoring job complete. Results are in {self.outputs_directory}")

        return df_score

    def run(self):
        return self.run_scoring()