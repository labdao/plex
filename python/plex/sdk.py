import json
import os
import subprocess
import logging
import itertools

from enum import Enum
from tempfile import TemporaryDirectory
from typing import Dict, List, Union



class ScatteringMethod(Enum):
    DOT_PRODUCT = 'dot_product'
    FLAT_CROSSPRODUCT = 'flat_crossproduct'
    NESTED_CROSSPRODUCT = 'nested_crossproduct'

# TODO: Move most of this logic to the go client
def generate_io_graph_from_tool(tool_filepath, scattering_method=ScatteringMethod.DOT_PRODUCT, **kwargs):
    # Open the file and load its content
    with open(tool_filepath, 'r') as f:
        tool = json.load(f)
    
    # Check if all kwargs are in the tool's inputs
    for arg in kwargs:
        if arg not in tool['inputs']:
            logging.error(f'The argument {arg} is not in the tool inputs.')
            logging.info(f'Available keys: {list(tool["inputs"].keys())}')
            raise ValueError(f'The argument {arg} is not in the tool inputs.')
    
    # Scattering methods
    if scattering_method == ScatteringMethod.DOT_PRODUCT:
        if len(set(len(x) for x in kwargs.values())) != 1:
            logging.error('All input arguments must have the same length for dot_product scattering method.')
            raise ValueError('All input arguments must have the same length for dot_product scattering method.')
        inputs_list = list(zip(*kwargs.values()))
    elif scattering_method in [ScatteringMethod.FLAT_CROSSPRODUCT, ScatteringMethod.NESTED_CROSSPRODUCT]:
        inputs_list = list(itertools.product(*kwargs.values()))
    else:
        logging.error(f'Invalid scattering method: {scattering_method}')
        raise ValueError(f'Invalid scattering method: {scattering_method}')
    
    # Build the io_json_graph
    io_json_graph = []
    for inputs in inputs_list:
        io_json_graph.append({
            'tool': tool_filepath,
            'inputs': {arg: {'class': tool['inputs'][arg]['type'], 'filepath': filepath} for arg, filepath in zip(kwargs.keys(), inputs)},
            'outputs': {arg: {'class': tool['outputs'][arg]['type'], 'filepath': ''} for arg in tool['outputs']},
            'state': 'created',
            'errMsg': '',
        })
    
    return io_json_graph

def run_plex(io: Union[Dict, List[Dict]], concurrency=1, local=False, verbose=False, retry=False, showAnimation=False, plex_path="./plex"):
    if not (isinstance(io, dict) or (isinstance(io, list) and all(isinstance(i, dict) for i in io))):
        raise ValueError('io must be a dict or a list of dicts')

    io_json_path = ""
    # Use a context manager for the temporary directory
    with TemporaryDirectory() as temp_dir:

        # Generate the JSON file name in the temporary directory
        json_file_path = os.path.join(temp_dir, 'io_data.json')

        # Save the io data to the JSON file
        with open(json_file_path, 'w') as json_file:
            json.dump(io, json_file, indent=4)

        cwd = os.getcwd()
        plex_work_dir = os.environ.get("PLEX_WORK_DIR",os.path.dirname(os.path.dirname(cwd)))
        cmd = [plex_path, "-input-io", json_file_path, "-concurrency", str(concurrency)]

        if local:
            cmd.append("-local=true")

        if verbose:
            cmd.append("-verbose=true")

        if retry:
            cmd.append("-retry=true")

        if not showAnimation: # default is true in the CLI
            cmd.append("-show-animation=false")

        with subprocess.Popen(cmd, stdout=subprocess.PIPE, bufsize=1, universal_newlines=True, cwd=plex_work_dir) as p:
            for line in p.stdout:
                if "Initialized IO file at:" in line:
                    parts = line.split()
                    io_json_path = parts[-1]
                print(line, end='')
    return io_json_path

def print_io_graph_status(io_graph):
    state_count = {}

    # Iterate through the io_list and count the occurrences of each state
    for io in io_graph:
        state = io['state']
        if state in state_count:
            state_count[state] += 1
        else:
            state_count[state] = 1

    # Print the total number of IOs
    print(f"Total IOs: {len(io_graph)}")

    # Print the number of IOs in each state
    for state, count in state_count.items():
        print(f"IOs in {state} state: {count}")

def mint_nft(tool_filepath, io_json_path, web3=True, plex_path="./plex"):
    # check if io_json_path is a valid file path
    if not os.path.isfile(io_json_path):
        raise ValueError('io_json_path must be a valid file path')
    
    # check if tool_filepath is a valid file path
    if not os.path.isfile(tool_filepath):
        raise ValueError('tool_filepath must be a valid file path')
    
    cwd = os.getcwd()
    plex_work_dir = os.environ.get("PLEX_WORK_DIR", os.path.dirname(os.path.dirname(cwd)))
    cmd = [plex_path, "-tool", tool_filepath, "-input-io", io_json_path]

    if web3:
        cmd.append("-web3=true")

    with subprocess.Popen(cmd, stdout=subprocess.PIPE, bufsize=1, universal_newlines=True, cwd=plex_work_dir) as p:
        for line in p.stdout:
            print(line, end='')