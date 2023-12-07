#!/usr/bin/env python3

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

def three_to_one(residue):
    # Dictionary to convert three-letter codes to one-letter codes
    conversion = {
        'ALA': 'A', 'ARG': 'R', 'ASN': 'N', 'ASP': 'D', 'CYS': 'C',
        'GLN': 'Q', 'GLU': 'E', 'GLY': 'G', 'HIS': 'H', 'ILE': 'I',
        'LEU': 'L', 'LYS': 'K', 'MET': 'M', 'PHE': 'F', 'PRO': 'P',
        'SER': 'S', 'THR': 'T', 'TRP': 'W', 'TYR': 'Y', 'VAL': 'V'
    }
    return conversion.get(residue, '?')

def main():
    pdb_file = "VTNC_HUMAN_2-41.pdb"
    chains = ['A', 'B']  # List of chains to extract
    sequences = {}

    for chain_id in chains:
        sequence = extract_sequence_from_pdb(pdb_file, chain_id)
        sequences[chain_id] = sequence

    # Write the sequences to a text file
    with open("amino_acid_sequences.txt", 'w') as output_file:
        for chain_id in chains:
            output_file.write(f"Chain {chain_id}: {sequences[chain_id]}\n")
        
        # Write the concatenated sequence of Chain A and Chain B
        if 'A' in sequences and 'B' in sequences:
            concatenated_sequence = sequences['A'] + ":" + sequences['B']
            output_file.write(f"Concatenated Sequence (Chain A:Chain B): {concatenated_sequence}\n")

if __name__ == "__main__":
    main()