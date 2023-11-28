import sys
import csv
from Bio import PDB
from collections import defaultdict
import os
from os.path import splitext
import itertools
import random

def find_chain_residue_range(pdb_path, chain_id):
    """
    Finds the start and end residue sequence indices for a given chain in a PDB file.
    """
    parser = PDB.PDBParser(QUIET=True)
    structure = parser.get_structure('protein', pdb_path)
    for model in structure:
        for chain in model:
            if chain.id == chain_id:
                residues = list(chain)
                if residues:
                    start_residue = residues[0].id[1]
                    end_residue = residues[-1].id[1]
                    return start_residue, end_residue
    return None, None

def find_interacting_residues(pdb_path, cutoff):
    """
    Finds interacting residues within a protein structure given a PDB file.
    """
    parser = PDB.PDBParser(QUIET=True)
    structure = parser.get_structure('protein', pdb_path)
    interactions = defaultdict(list)

    for chain_1 in structure.get_chains():
        for chain_2 in structure.get_chains():
            if chain_1 == chain_2:
                continue
            for residue_1 in chain_1:
                for residue_2 in chain_2:
                    if any(atom_1 - atom_2 < cutoff for atom_1 in residue_1 for atom_2 in residue_2):
                        interactions[chain_1.id].append(residue_1.id[1])
                        interactions[chain_2.id].append(residue_2.id[1])

    for chain, residues in interactions.items():
        interactions[chain] = sorted(set(residues))

    return interactions

def create_contact_domains(contacts, domain_threshold):
    """
    Groups interacting residues into domains.
    """
    contact_domains = []
    while contacts:
        help_list = [contacts.pop(0)]
        while contacts and contacts[0] - help_list[-1] < domain_threshold:
            help_list.append(contacts.pop(0))

        # contact_domains.append(help_list)
        domain_range = list(range(help_list[0], help_list[-1] + 1))
        contact_domains.append(domain_range)

    return contact_domains

def create_sublists(contact_domains, chain_length):
    """
    Creates sublists from contact domains, including non-contact domains.
    """
    reference_list = list(range(1, chain_length + 1))
    output = []
    contact_flag = []
    last_contact_end = 0

    for domain in contact_domains:
        non_contact_start = last_contact_end + 1
        non_contact_end = domain[0] - 1
        if non_contact_start <= non_contact_end:
            output.append(list(range(non_contact_start, non_contact_end + 1)))
            contact_flag.append(False)

        output.append(domain)
        contact_flag.append(True)
        last_contact_end = domain[-1]

    if last_contact_end < chain_length:
        output.append(list(range(last_contact_end + 1, chain_length + 1)))
        contact_flag.append(False)

    return output, contact_flag

def identify_domains(pdb_path, interactions, binder_chain, domain_threshold):
    """
    Identifies domains based on interactions.
    """
    if binder_chain not in interactions:
        return [], []

    contacts = interactions[binder_chain]
    contact_domains = create_contact_domains(contacts, domain_threshold)
    _, chain_length = find_chain_residue_range(pdb_path, binder_chain)
    domains, contact_flag = create_sublists(contact_domains, chain_length)

    return domains, contact_flag

def prompt_tokens(domains, contact_flag, chain_name):
    """
    Generates prompt tokens based on domains and contact flags.
    """
    token_alphabet = []

    for i, (domain, contact) in enumerate(zip(domains, contact_flag)):
        if domain:
            token = {
                'label': i + 1,
                'domain_length': len(domain),
                'domain': f'{chain_name}{domain[0]}-{domain[-1]}',
                'contact': contact
            }
            token_alphabet.append(token)

    return token_alphabet

def generate_combinations(contact_flags, n_samples, p_masking_contact_domain, p_masking_noncontact_domain):
    """
    Generates combinations of True/False based on contact and non-contact domain probabilities.
    """
    probabilities = [p_masking_contact_domain if flag else p_masking_noncontact_domain for flag in contact_flags]
    bool_lists = []
    for _ in range(n_samples):
        bool_list = [random.random() > p for p in probabilities]
        bool_lists.append(bool_list)
    return bool_lists

def generate_full_prompts(pdb_path, binder_chain, target_chain, target_start_residue, target_end_residue, cutoff, n_samples, p_masking_contact_domain, p_masking_noncontact_domain, domain_threshold):
    """
    Computes the inter-protein contact domains. Generates an alphabet of domain tokens. Randomly masks contact and non-contact domains according to preset probabilities.
    """
    interactions = find_interacting_residues(pdb_path, cutoff)
    domains, contact_flag = identify_domains(pdb_path, interactions, binder_chain, domain_threshold)
    token_alphabet = prompt_tokens(domains, contact_flag, binder_chain)
    # print('alphabet of tokens', token_alphabet)

    contact_flags = [token['contact'] for token in token_alphabet]
    domain_lengths = [token['domain_length'] for token in token_alphabet]
    domains = [token['domain'] for token in token_alphabet]

    # target_start_residue, target_end_residue = find_chain_residue_range(pdb_path, target_chain)
    target_binder_range = f"{target_chain}{target_start_residue}-{target_end_residue}"

    masking_charts = generate_combinations(contact_flag, n_samples, p_masking_contact_domain, p_masking_noncontact_domain)
    # print('masking charts', masking_charts)
    full_prompts = []
    for chart in masking_charts:
        prompt = f"{target_binder_range}:"
        for i, mask in enumerate(chart):
            prompt += domains[i] if mask else str(domain_lengths[i])
            if i < len(chart) - 1:
                prompt += "/"
        full_prompts.append(prompt)

    return full_prompts

# # Example usage
# reference_protein_complex = 'summary/UROK_HUMAN_1-133.pdb'
# binder_chain = 'B'
# target_chain = 'A'
# cutoff = 5.0 # distance to define inter-protein contacts (in Angstrom)
# n_samples = 5 # total number of prompts generated
# p_masking_contact_domain = 0.6 # probability of masking a contact domain
# p_masking_noncontact_domain = .1 # probability of masking a non-contact domain
# domain_distance_threshold = 6 # definition of constitutes separate domains (in units of residues)

# full_prompts = generate_full_prompts(reference_protein_complex, binder_chain, target_chain, cutoff, n_samples, p_masking_contact_domain, p_masking_noncontact_domain, domain_distance_threshold)
# for prompt in full_prompts:
#     print(prompt)