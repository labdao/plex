import os
import time
import pandas as pd

import hydra
from omegaconf import DictConfig, OmegaConf

from utils import slash_to_convexity_notation
from utils import user_input_parsing
from utils import replace_invalid_characters

import json

import logging

from sampler import Sampler
from generator_module import Generator
from scorer_module import Scorer
from selector_module import SequenceSelector

def get_plex_job_inputs():
    # Retrieve the environment variable
    json_str = os.getenv("PLEX_JOB_INPUTS")

    # Check if the environment variable is set
    if json_str is None:
        raise ValueError("PLEX_JOB_INPUTS environment variable is missing.")

    # Convert the JSON string to a Python dictionary
    try:
        data = json.loads(json_str)
        return data
    except json.JSONDecodeError:
        # Handle the case where the string is not valid JSON
        raise ValueError("PLEX_JOB_INPUTS is not a valid JSON string.")

def apply_initial_permissibility_vector(seed, permissibility_seed, cfg):

    mod_sequence = []
    seed_list = list(seed)

    for i, char in enumerate(permissibility_seed):
        if char == 'X' or char == '*' or char in cfg.params.basic_settings.alphabet:
            mod_sequence.append(seed_list[i])
        elif char == '-':
            mod_sequence.append('-')      
    
    mod_sequence = ''.join(mod_sequence)

    return mod_sequence

# def load_initial_data(cfg, outputs_directory):
#     sequence_input = cfg.params.basic_settings.sequence_input
#     binder, target = [s.replace(" ", "") for s in sequence_input.split(';')]
#     binder = binder.upper()
#     target = target.upper()

#     binder = replace_invalid_characters(binder, cfg.params.basic_settings.alphabet)

#     sequences = [{
#         't': 0,
#         'sample_number': 0,
#         'seed': binder,
#         'permissibility_seed': '',
#         '(levenshtein-distance, mask)': 'none',
#         'modified_seq': '',
#         'permissibility_modified_seq': '',
#         'acceptance_flag': True  # manual selection of starting sequence
#     }]

#     if 'X' in sequences[-1]['seed'] or '*' in sequences[-1]['seed']:  # sequence completion
#         seed = sequences[-1]['seed']
#         generator = Generator(cfg, outputs_directory)
#         seed, _, _ = generator.run(0, 1, seed, '', '', None)
#         del generator

#     contig_in_convexity_notation = ''
#     if all(char in cfg.params.basic_settings.alphabet for char in seed):        
#         if cfg.params.basic_settings.init_permissibility_vec == "":
#             contig_in_convexity_notation = sequences[-1]['seed']
#         else:
#             logging.info(f"converting to convexity notation")
#             contig_in_convexity_notation = slash_to_convexity_notation(seed, cfg.params.basic_settings.init_permissibility_vec)

#         # OmegaConf.update(cfg, "params.RFdiffusion_settings.hotspots", "[]", merge=False)

#     else:
#         contig = f"x1:{len(seed)}"
#         OmegaConf.update(cfg, "params.basic_settings.init_permissibility_vec", contig, merge=False)
#         contig_in_convexity_notation = slash_to_convexity_notation(seed, cfg.params.basic_settings.init_permissibility_vec)

#     sequences[-1]['seed'] = seed

#     logging.info(f"contig_in_convexity_notation, {contig_in_convexity_notation}")
#     sequences[-1]['modified_seq'] += apply_initial_permissibility_vector(sequences[-1]['seed'], contig_in_convexity_notation, cfg)
#     logging.info(f"modified sequence, {sequences[-1]['modified_seq']}")
#     sequences[-1]['permissibility_seed'] += contig_in_convexity_notation
#     sequences[-1]['permissibility_modified_seq'] += contig_in_convexity_notation

#     OmegaConf.update(cfg, "params.basic_settings.target_seq", target, merge=False)

#     return pd.DataFrame(sequences), cfg

def load_initial_data_and_determine_logic(cfg, outputs_directory):
    sequence_input = cfg.params.basic_settings.sequence_input
    binder, target = [s.replace(" ", "") for s in sequence_input.split(';')]
    binder = binder.upper()
    target = target.upper()

    binder = replace_invalid_characters(binder, cfg.params.basic_settings.alphabet)
    if len(binder) > 300:
        binder = binder[:300]
        OmegaConf.update(cfg, "params.basic_settings.init_permissibility_vec", cfg.params.basic_settings.init_permissibility_vec[:300], merge=False)
        logging.info("The binder sequence is longer than 300 characters and has been cropped.")

    if len(target) > 300:
        target = target[:300]
        logging.info("The target sequence is longer than 300 characters and has been cropped.")

    sequences = [{
        't': 0,
        'sample_number': 0,
        'seed': binder,
        'permissibility_seed': '',
        '(levenshtein-distance, mask)': 'none',
        'modified_seq': '',
        'permissibility_modified_seq': '',
        'acceptance_flag': True  # manual selection of starting sequence
    }]

    if all(char in ['X', '*'] for char in sequences[-1]['seed']): # completely undetermined binder sequence
        logging.info(f"running algorithm for completely undetermined sequence")

        seed = sequences[-1]['seed'].replace('*', 'X')  # replace all '*' characters with 'X' in seed
        contig = f"x1:{len(seed)}"
        OmegaConf.update(cfg, "params.basic_settings.init_permissibility_vec", contig, merge=False)
        contig_in_convexity_notation = slash_to_convexity_notation(seed, cfg.params.basic_settings.init_permissibility_vec)
        sequences[-1]['seed'] = seed

        OmegaConf.update(cfg, "params.basic_settings.generator", 'RFdiff+ProteinMPNN+ESM2', merge=False)
        if cfg.params.basic_settings.high_fidelity:
            OmegaConf.update(cfg, "params.basic_settings.scorers", 'colabfold,prodigy', merge=False)
        else:
            OmegaConf.update(cfg, "params.basic_settings.scorers", 'initial_omegafold,omegafold+diffdock', merge=False)
    
    elif all(char in cfg.params.basic_settings.alphabet for char in sequences[-1]['seed']): # completely determined sequence
        logging.info(f"running algorithm for completely determined sequence")
        
        OmegaConf.update(cfg, "params.basic_settings.generator", '[]', merge=False)
        if cfg.params.basic_settings.high_fidelity:
            OmegaConf.update(cfg, "params.basic_settings.scorers", 'colabfold,prodigy', merge=False)
        else:
            OmegaConf.update(cfg, "params.basic_settings.scorers", 'omegafold+diffdock,prodigy', merge=False)


    elif 'X' in sequences[-1]['seed'] or '*' in sequences[-1]['seed']:  # partially determined sequence
        seed = sequences[-1]['seed']
        generator = Generator(cfg, outputs_directory)
        seed, _, _ = generator.run(0, 1, seed, '', '', None)
        del generator

        contig_in_convexity_notation = ''
        if all(char in cfg.params.basic_settings.alphabet for char in seed):
            logging.info(f"running algorithm for partially determined sequence")        
            if cfg.params.basic_settings.init_permissibility_vec == "":
                contig_in_convexity_notation = sequences[-1]['seed']
            else:
                logging.info(f"converting to convexity notation")
                contig_in_convexity_notation = slash_to_convexity_notation(seed, cfg.params.basic_settings.init_permissibility_vec)

            OmegaConf.update(cfg, "params.basic_settings.generator", 'RFdiff+ProteinMPNN', merge=False)
            if cfg.params.basic_settings.high_fidelity:
                OmegaConf.update(cfg, "params.basic_settings.scorers", 'colabfold,prodigy', merge=False)
            else:
                OmegaConf.update(cfg, "params.basic_settings.scorers", 'omegafold+diffdock,prodigy', merge=False)

        else:
            logging.info(f"running algorithm for completely undetermined sequence")
            contig = f"x1:{len(seed)}"
            OmegaConf.update(cfg, "params.basic_settings.init_permissibility_vec", contig, merge=False)
            contig_in_convexity_notation = slash_to_convexity_notation(seed, cfg.params.basic_settings.init_permissibility_vec)
            OmegaConf.update(cfg, "params.basic_settings.generator", 'RFdiff+ProteinMPNN+ESM2', merge=False)
            if cfg.params.basic_settings.high_fidelity:
                OmegaConf.update(cfg, "params.basic_settings.scorers", 'colabfold,prodigy', merge=False)
            else:
                OmegaConf.update(cfg, "params.basic_settings.scorers", 'omegafold+diffdock,prodigy', merge=False)

        sequences[-1]['seed'] = seed

    logging.info(f"contig_in_convexity_notation, {contig_in_convexity_notation}")
    sequences[-1]['modified_seq'] += apply_initial_permissibility_vector(sequences[-1]['seed'], contig_in_convexity_notation, cfg)
    logging.info(f"modified sequence, {sequences[-1]['modified_seq']}")
    sequences[-1]['permissibility_seed'] += contig_in_convexity_notation
    sequences[-1]['permissibility_modified_seq'] += contig_in_convexity_notation

    OmegaConf.update(cfg, "params.basic_settings.target_seq", target, merge=False)

    return pd.DataFrame(sequences), cfg

@hydra.main(version_base=None, config_path="conf", config_name="config")
def my_app(cfg: DictConfig) -> None:

    # defining output directory
    if cfg.outputs.directory is None:
        outputs_directory = hydra.core.hydra_config.HydraConfig.get().runtime.output_dir
    else:
        outputs_directory = cfg.outputs.directory

    ## plex user inputs # some of these are currently not used!
    user_inputs = get_plex_job_inputs()
    permissibility_seed = user_inputs["init_permissibility_vec"]
    logging.info(f"user inputs from plex: {user_inputs}")

    cfg = user_input_parsing(cfg, user_inputs)

    logging.info(f"{OmegaConf.to_yaml(cfg)}")
    logging.info(f"Working directory : {os.getcwd()}")

    logging.info(f"inputs directory: {cfg.inputs.directory}")

    start_time = time.time()

    # df, cfg = load_initial_data(cfg, outputs_directory)
    df, cfg = load_initial_data_and_determine_logic(cfg, outputs_directory)


    seed_row = df[(df['t']==0) & (df['acceptance_flag'] == True)]
    seed = seed_row['modified_seq'].values[0]
    permissibility_seed = seed_row['permissibility_modified_seq'].values[0]
    logging.info(f"target sequence {cfg.params.basic_settings.target_seq}")
    logging.info(f"initial seed sequence {seed}")

    generator = Generator(cfg, outputs_directory)
    scorer = Scorer(cfg, outputs_directory)
    selector = SequenceSelector(cfg)
    sampler = Sampler(cfg, outputs_directory, generator, selector, scorer, cfg.params.basic_settings.evolve, cfg.params.basic_settings.n_samples)


    for t in range(cfg.params.basic_settings.number_of_binders):

        logging.info(f"starting evolution step, {t+1}")
        logging.info(f"seed sequence, {seed}")

        mod_seq, modified_permissibility_seq, df = sampler.run(t+1, seed, permissibility_seed, df)

        logging.info(f"modified sequence, {mod_seq}")
        logging.info(f"modified permissibility vector, {modified_permissibility_seq}")

        df.to_csv(f"{outputs_directory}/summary.csv", index=False)

        if cfg.params.basic_settings.evolve:
            seed = mod_seq
            permissibility_seed = modified_permissibility_seq

        print('\n')

    end_time = time.time()
    duration = end_time - start_time

    logging.info("sequence to structure complete...")
    logging.info(f"executed in {duration:.2f} seconds.")

if __name__ == "__main__":
    my_app()