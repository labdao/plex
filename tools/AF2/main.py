import glob
import os
import time
import signal
import sys
import random
import string
import subprocess
import re
import json
# import numpy as np

import subprocess
from folding import AF2Runner
import hydra
from hydra import compose, initialize
from hydra.core.hydra_config import HydraConfig
from omegaconf import DictConfig, OmegaConf

# import pandas as pd
import yaml

def seq_to_struc(fasta_file, cfg):

    # # check whether input_seqs is a fasta or a sequence. If fasta, then extract the sequences

    # folded_seqs = [] # save as a fasta and overwrite in the 
    # run loop over the sequences
    af2_runner = AF2Runner(fasta_file, cfg.outputs.directory)
    af2_runner.run()

        # produce n AF2 model in a for loop
        # write add the seq, pdb file identifier (absolute path), and metrics to a csv file
        # extract sequences from pdb and append to fasta 
        # append pdbs file identifier to list of pdbs
    
    # call sequence_to_structure recursively seq_to_struc([], )


@hydra.main(version_base=None, config_path="conf", config_name="config_seq2struc")
def my_app(cfg: DictConfig) -> None:

    print(OmegaConf.to_yaml(cfg))
    print(f"Working directory : {os.getcwd()}")

    # defining output directory
    if cfg.outputs.directory is None:
        outputs_directory = hydra.core.hydra_config.HydraConfig.get().runtime.output_dir
    else:
        outputs_directory = cfg.outputs.directory
    print(f"Output directory : {outputs_directory}")

    fasta_file = cfg.inputs.directory
    
    # load fasta file with list of sequences
    start_time = time.time()

    # sequence to structure function

    seq_to_struc(fasta_file, cfg)

    print("sequence to structure complete...")
    end_time = time.time()
    duration = end_time - start_time
    print(f"executed in {duration:.2f} seconds.")

if __name__ == "__main__":
    my_app()

## goal definition
# df with sequences to compare within the benchmark
# for every sequence in the task run folding - k times, default for k = 1 (this will require you to do a bit of fasta glue code)
# add the paths of the folded proteins to the
# for every structure path run the minimize protein function
# add the columns with the new scores to the data frame 
# the final output of running this script and an input csv is an output csv that holds additional columns which point to structures, minimized structures, af-metrics as well as the minimization metrics