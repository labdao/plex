import json
import os
import subprocess
import logging
import itertools
import tempfile

from enum import Enum
from typing import Dict, List


class ScatteringMethod(Enum):
    DOT_PRODUCT = 'dot_product'
    FLAT_CROSSPRODUCT = 'flat_crossproduct'  # can depreciate
    NESTED_CROSSPRODUCT = 'nested_crossproduct'  # can depreciate
    CROSS_PRODUCT = 'cross_product'


class CoreTools(Enum):
    EQUIBIND = "QmZ2HarAgwZGjc3LBx9mWNwAQkPWiHMignqKup1ckp8NhB"
    DIFFDOCK = "QmSzetFkveiQYZ5FgpZdHHfsjMWYz5YzwMAvqUgUFhFPMM"
    COLABFOLD_MINI = "QmcRH74qfqDBJFku3mEDGxkAf6CSpaHTpdbe1pMkHnbcZD"
    COLABFOLD_STANDARD = "QmXnM1VpdGgX5huyU3zTjJovsu42KPfWhjxhZGkyvy9PVk"
    COLABFOLD_LARGE = "QmPYqMy19VFFuYztL6b5ruo4Kw4JWT583emStGrSYTH5Yi"
    BAM2FASTQ = "QmbPUirWiWCv9sgdHLekf5AnoCdw4QPU2SyfGGKs9JRRbq"
    ODDT = "QmUx7NdxkXXZvbK1JXZVUYUBqsevWkbVxgTzpWJ4Xp4inf"
    RFDIFFUSION = "QmXnCBCtoYuPyGsEJVpjn5regHfFSYa8kx44e22XxDX2t2"
    REPEATMODELER = "QmZdXxnUt1sFFR39CfkEUgiioUBf6qP5CUs8TCb7Wqn4MC"
    GNINA = "QmYfGaWzxwi8HiWLdiX4iQXuuLXVKYrr6YC3DknEvZeSne"
    BATCH_DLKCAT = "QmThdvypN8gDDwwyNnpSYsdwvyxCET8s1jym3HZCTaBzmD"
    OPENBABEL_PDB_TO_SDF = "QmbbDSDZJp8G7EFaNKsT7Qe7S9iaaemZmyvS6XgZpdR5e3"
    OPENBABEL_RMSD = "QmUxrKgAs5r42xVki4vtMskJa1Z7WA64wURkwywPMch7dA"


def isFilePath(file_path: str):
    if os.path.isfile(file_path):
        return True
    elif os.path.isdir(file_path):
        return True
    elif os.path.isfile(file_path):
        return True
    elif "/" in file_path:
        return True
    else:
        return False


def generate_io_graph_from_tool(toolCID, scattering_method=ScatteringMethod.DOT_PRODUCT, **kwargs):
    if isFilePath(toolCID):
        print("Use plex_upload to upload a tool to IPFS")
        raise ValueError(f'The tool argument must be a CID, not a filepath: {tool}')

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
    elif scattering_method in [ScatteringMethod.CROSS_PRODUCT, ScatteringMethod.FLAT_CROSSPRODUCT, ScatteringMethod.NESTED_CROSSPRODUCT]:
        inputs_list = list(itertools.product(*kwargs.values()))
    else:
        logging.error(f'Invalid scattering method: {scattering_method}')
        raise ValueError(f'Invalid scattering method: {scattering_method}')

    # Build the io_json_graph
    io_json_graph = []
    for inputs in inputs_list:
        io_json_graph.append({
            'tool': {},
            'inputs': {arg: {'class': tool['inputs'][arg]['type'], 'filepath': filepath} for arg, filepath in zip(kwargs.keys(), inputs)},
            'outputs': {arg: {'class': tool['outputs'][arg]['type'], 'filepath': ''} for arg in tool['outputs']},
            'state': 'created',
            'errMsg': '',
        })

    # Save io_json_graph to temp file
    with tempfile.NamedTemporaryFile(delete=False, suffix=".json") as temp:
        json.dump(io_json_graph, temp)
        tempIoJsonFilePath = temp.name

    io_json_cid = plex_upload(tempIoJsonFilePath)
    os.remove(tempIoJsonFilePath)

    return io_json_cid, io_json_graph


def plex_upload(filePath: str, plex_path="plex"):
    cmd = [plex_path, "upload", "-p", filePath]

    with subprocess.Popen(cmd, stdout=subprocess.PIPE, bufsize=1, universal_newlines=True) as p:
        for line in p.stdout:
            if "FILL IN WITH CORRECT LOG:" in line:
                parts = line.split()
                io_json_cid = parts[-1]
            print(line, end='')
    return io_json_cid


def plex_create(toolpath: str, inputDir: str, layers=2, outputDir="", verbose=False, showAnimation=False, concurrency="1", annotations=[], plex_path="plex"):
    cwd = os.getcwd()
    plex_work_dir = os.environ.get("PLEX_WORK_DIR", os.path.dirname(os.path.dirname(cwd)))
    cmd = [plex_path, "create", "-t", toolpath, "-i", inputDir, f"--layers={layers}"]

    if outputDir:
        cmd.append(f"-o={outputDir}")

    if verbose:
        cmd.append("-v=true")

    if concurrency:
        cmd.append(f"--concurrency={concurrency}")

    if annotations:
        cmd.append(f"--annotations={annotations.join(',')}")

    if not showAnimation: # default is true in the CLI
        cmd.append("--showAnimation=false")

    with subprocess.Popen(cmd, stdout=subprocess.PIPE, bufsize=1, universal_newlines=True, cwd=plex_work_dir) as p:
        for line in p.stdout:
            if "Initial IO JSON file CID:" in line:
                parts = line.split()
                io_json_cid = parts[-1]
            print(line, end='')
    return io_json_cid


def plex_run(io_json_cid: str, outputDir="", verbose=False, showAnimation=False, concurrency="1", annotations=[], plex_path="plex"):
    cwd = os.getcwd()
    plex_work_dir = os.environ.get("PLEX_WORK_DIR", os.path.dirname(os.path.dirname(cwd)))
    cmd = [plex_path, "run", "-i", io_json_cid]

    if outputDir:
        cmd.append(f"-o={outputDir}")

    if verbose:
        cmd.append("-v=true")

    if concurrency:
        cmd.append(f"--concurrency={concurrency}")

    if annotations:
        cmd.append(f"--annotations={annotations.join(',')}")

    if not showAnimation: # default is true in the CLI
        cmd.append("--showAnimation=false")

    with subprocess.Popen(cmd, stdout=subprocess.PIPE, bufsize=1, universal_newlines=True, cwd=plex_work_dir) as p:
        for line in p.stdout:
            if "Completed IO JSON CID:" in line:
                parts = line.split()
                io_json_cid = parts[-1]
            if "Initialized IO file at:" in line:
                parts = line.split()
                io_json_local_filepath = parts[-1]
            print(line, end='')
    return io_json_cid, io_json_local_filepath


def plex_mint(io_json_cid: str, imageCid="", plex_path="plex"):
    cwd = os.getcwd()
    plex_work_dir = os.environ.get("PLEX_WORK_DIR", os.path.dirname(os.path.dirname(cwd)))
    cmd = [plex_path, "mint", "-i", io_json_cid]

    if imageCid:
        cmd.append(f"-imageCid={imageCid}")

    with subprocess.Popen(cmd, stdout=subprocess.PIPE, bufsize=1, universal_newlines=True, cwd=plex_work_dir) as p:
        for line in p.stdout:
            print(line, end='')


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
