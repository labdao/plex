import json
import os
import subprocess

from enum import Enum
from typing import Dict, List


class ScatteringMethod(Enum):
    DOT_PRODUCT = 'dotproduct'
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


def plex_init(toolpath: str, scatteringMethod="dotProduct", plex_path="plex", **kwargs):
    cwd = os.getcwd()
    plex_work_dir = os.environ.get("PLEX_WORK_DIR", os.path.dirname(os.path.dirname(cwd)))

    # Convert kwargs dictionary to a JSON string
    inputs = json.dumps(kwargs)

    cmd = [plex_path, "init", "-t", toolpath, "-i", inputs, f"--scatteringMethod={scatteringMethod}"]

    io_json_cid = ""
    with subprocess.Popen(cmd, stdout=subprocess.PIPE, bufsize=1, universal_newlines=True, cwd=plex_work_dir) as p:
        for line in p.stdout:
            if "Pinned IO JSON CID:" in line:
                parts = line.split()
                io_json_cid = parts[-1]
            print(line, end='')

    return io_json_cid


def plex_vectorize(io_path: str, tool_cid: str, outputDir="", plex_path="plex"):
    cwd = os.getcwd()
    plex_work_dir = os.environ.get("PLEX_WORK_DIR", os.path.dirname(os.path.dirname(cwd)))

    cmd = [plex_path, "vectorize", "-i", io_path, "-t", tool_cid, "-o", outputDir]
    with subprocess.Popen(cmd, stdout=subprocess.PIPE, bufsize=1, universal_newlines=True, cwd=plex_work_dir) as p:
        outvects = ""
        for line in p.stdout:
            if "Output Vectors were saved at:" in line:
                parts = line.split()
                io_vector_outpath = parts[-1]
                with open(io_vector_outpath, 'r') as f:
                    outvects = json.load(f)
                os.remove(io_vector_outpath)
            print(line, end='')
    return outvects


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
