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

    # Parse the log file and write to CSV
    csv_path = os.path.join(os.path.dirname(log_path), 'summary.csv')
    with open(log_path, 'r') as log_file, open(csv_path, 'w', newline='') as csv_file:
        csv_writer = csv.writer(csv_file)
        csv_writer.writerow(['Sequence', 'Query', 'Data'])

        current_seq = ''
        for line in log_file:
            if 'Query' in line:
                current_seq = line.split(':')[1].strip()
            if 'pLDDT' in line:
                data = line.split(',')[1].strip()
                csv_writer.writerow([sequences.get(current_seq, ''), current_seq, data])

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

    # Example usage
    parse_log_and_fasta_to_csv(os.path.abspath(fasta_file), os.path.join(outputs_directory, 'log.txt'))

    print("Sequence to structure complete...")
    end_time = time.time()
    duration = end_time - start_time
    print(f"executed in {duration:.2f} seconds.")

if __name__ == "__main__":
    my_app()