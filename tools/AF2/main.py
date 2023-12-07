# import glob
import os
import time
import pandas as pd
from AF2_module import AF2Runner
import csv

import hydra
from hydra import compose, initialize
from hydra.core.hydra_config import HydraConfig
from omegaconf import DictConfig, OmegaConf

def parse_log_and_fasta_to_csv(fasta_path, log_path):
    # Read and store sequences from the FASTA file
    sequences = {}
    with open(fasta_path, 'r') as fasta_file:
        identifier = ''
        for line in fasta_file:
            if line.startswith('>'):
                identifier = line[1:].strip()
            else:
                sequences[identifier] = line.strip()

    # Initialize the columns for the CSV
    columns = ['Sequence', 'Query', 'structure_pdbs', 'recycle', 'pLDDT', 'pTM', 'tol', 'rank']
    
    # Parse the log file and write to CSV
    csv_path = os.path.join(os.path.dirname(log_path), 'summary.csv')
    with open(log_path, 'r') as log_file, open(csv_path, 'w', newline='') as csv_file:
        csv_writer = csv.DictWriter(csv_file, fieldnames=columns)
        csv_writer.writeheader()

        current_query = ''
        for line in log_file:
            if 'Query' in line:
                # Extracting the query value (e.g., "1/2" or "2/2")
                current_query = line.split('Query')[1].split(':')[0].strip()

            if 'pLDDT' in line and 'pTM' in line:
                parts = line.split()
                data = parts[2:]  # Skip the first two substrings (date and time)
                row_data = {col: '' for col in columns}  # Initialize row with empty values

                # Extract the sequence number (XXX in "XXX/YYY") and get the corresponding sequence
                sequence_number = current_query.split('/')[0]
                row_data['Sequence'] = sequences.get(sequence_number, '')
                row_data['Query'] = current_query

                # Handle the first substring without an equal sign
                row_data['structure_pdbs'] = data[0] if data else ''

                # Process remaining substrings with equal signs
                for item in data[1:]:
                    if '=' in item:
                        key, value = item.split('=')
                        if key in columns:
                            row_data[key] = value

                csv_writer.writerow(row_data)

    print(f"CSV file written to {csv_path}")

def write_dataframe_to_fastas(dataframe, cfg):
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
        file_path = os.path.join(input_dir, f"seq_{row['sequence_number']}.fasta")
        with open(file_path, 'w') as file:
            file.write(f">{row['sequence_number']}\n{row['seq']}\n")
    return os.path.abspath(input_dir)


def find_fasta_file(directory_path):
    for root, dirs, files in os.walk(directory_path):
        for file in files:
            if file.endswith(".fasta"):
                return os.path.abspath(os.path.join(root, file))
    return None  # Return None if no .fasta file is found in the directory


def load_fasta_to_dataframe(fasta_file):
    sequences = []
    with open(fasta_file, 'r') as file:
        seq_num = 1
        for line in file:
            if line.startswith('>'):
                sequences.append({'sequence_number': seq_num, 'seq': ''})
                seq_num += 1
            else:
                sequences[-1]['seq'] += line.strip()

    return pd.DataFrame(sequences)

def seq2struc(df, outputs_directory, cfg):

    seq_input_dir = write_dataframe_to_fastas(df, cfg)

    af2_runner = AF2Runner(seq_input_dir, outputs_directory)
    af2_runner.run()
        
    return None

@hydra.main(version_base=None, config_path="conf", config_name="config")
def my_app(cfg: DictConfig) -> None:

    print(OmegaConf.to_yaml(cfg))
    print(f"Working directory : {os.getcwd()}")

    # defining output directory
    if cfg.outputs.directory is None:
        outputs_directory = hydra.core.hydra_config.HydraConfig.get().runtime.output_dir
    else:
        outputs_directory = cfg.outputs.directory
    print(f"Output directory : {outputs_directory}")

    fasta_file = find_fasta_file(cfg.inputs.directory) # load fasta with inital sequences and convert to data frame
    df = load_fasta_to_dataframe(fasta_file)

    start_time = time.time()

    seq2struc(df, outputs_directory, cfg)

    # create and write a csv file with sequence and metric information for each output structure
    parse_log_and_fasta_to_csv(os.path.abspath(fasta_file), os.path.join(outputs_directory, 'log.txt'))

    print("Sequence to structure complete...")
    end_time = time.time()
    duration = end_time - start_time
    print(f"executed in {duration:.2f} seconds.")

if __name__ == "__main__":
    my_app()