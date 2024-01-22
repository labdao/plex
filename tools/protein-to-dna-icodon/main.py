import os
import json
import subprocess
import sys
from dnachisel import reverse_translate

def get_plex_job_inputs():
    json_str = os.getenv("PLEX_JOB_INPUTS")

    if json_str is None:
        raise ValueError("PLEX_JOB_INPUTS environment variable is missing.")

    try:
        data = json.loads(json_str)
        return data
    except json.JSONDecodeError:
        raise ValueError("PLEX_JOB_INPUTS is not a valid JSON string.")

def do_reverse_translate(protein_sequence):
    dna_sequence = reverse_translate(protein_sequence, randomize_codons=False, table='Standard')
    with open('/outputs/dna_sequence_after_reverse_translation.txt', 'w') as file:
        file.write(dna_sequence)

def three_to_one(residue):
    # Dictionary to convert three-letter codes to one-letter codes
    conversion = {
        'ALA': 'A', 'ARG': 'R', 'ASN': 'N', 'ASP': 'D', 'CYS': 'C',
        'GLN': 'Q', 'GLU': 'E', 'GLY': 'G', 'HIS': 'H', 'ILE': 'I',
        'LEU': 'L', 'LYS': 'K', 'MET': 'M', 'PHE': 'F', 'PRO': 'P',
        'SER': 'S', 'THR': 'T', 'TRP': 'W', 'TYR': 'Y', 'VAL': 'V'
    }
    return conversion.get(residue, '?')

def extract_sequence_from_pdb(pdb_file, chain_id):
    sequence = ''
    residues = {}
    with open(pdb_file, 'r') as file:
        for line in file:
            if line.startswith("ATOM") and line[21] == chain_id:
                residue_name = line[17:20].strip()
                residue_number = int(line[22:26].strip())
                # Ensure each residue is counted only once per chain
                if (residue_number, chain_id) not in residues:
                    residues[(residue_number, chain_id)] = residue_name
                    sequence += three_to_one(residue_name)
    return sequence

def extract_sequence(pdb_file):
    chains = ['A', 'B']  # List of chains to extract
    sequences = {}

    for chain_id in chains:
        sequence = extract_sequence_from_pdb(pdb_file, chain_id)
        sequences[chain_id] = sequence

    return sequences['B']

def main():
    # Get the job inputs from the environment variable
    try:
        job_inputs = get_plex_job_inputs()
        print("Job Inputs:", job_inputs)
    except ValueError as e:
        print(e)
        sys.exit(1)

    # Create /outputs directory if it doesn't exist
    os.makedirs("/outputs", exist_ok=True)

    input_file = job_inputs["input_file"]
    if(input_file.endswith('pdb')):
        protein_sequence = extract_sequence(input_file)
    elif(input_file.endswith('txt')):
        with open(input_file, 'r') as file:
            protein_sequence = file.read()

    do_reverse_translate(protein_sequence)
    
    species = job_inputs["species"]
    iterations = job_inputs["iterations"]
    make_more_optimal = job_inputs["make_more_optimal"]

    r_command = ["/usr/bin/Rscript", 
                 "iCodonScript.R",
                 species,
                 str(iterations),
                 str(make_more_optimal)
                ]
    subprocess.run(r_command)


    print("done")


if __name__ == "__main__":
    main()
