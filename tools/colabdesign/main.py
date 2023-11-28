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
import numpy as np
import subprocess

import pandas as pd
import yaml

import prompt_generator

print("Starting main.py...")

# setup
if not os.path.isdir("params"):
    os.system("apt-get install aria2")
    os.system("mkdir params")
    # send param download into background
    os.system(
        "(\
  aria2c -q -x 16 https://files.ipd.uw.edu/krypton/schedules.zip; \
  aria2c -q -x 16 http://files.ipd.uw.edu/pub/RFdiffusion/6f5902ac237024bdd0c176cb93063dc4/Base_ckpt.pt; \
  aria2c -q -x 16 http://files.ipd.uw.edu/pub/RFdiffusion/e29311f6f1bf1af907f9ef9f44b8328b/Complex_base_ckpt.pt; \
  aria2c -q -x 16 http://files.ipd.uw.edu/pub/RFdiffusion/f572d396fae9206628714fb2ce00f72e/Complex_beta_ckpt.pt; \
  aria2c -q -x 16 https://storage.googleapis.com/alphafold/alphafold_params_2022-12-06.tar; \
  tar -xf alphafold_params_2022-12-06.tar -C params; \
  touch params/done.txt) &"
    )

if not os.path.isdir("RFdiffusion"):
    print("installing RFdiffusion...")
    os.system("git clone https://github.com/sokrypton/RFdiffusion.git")
    os.system("pip -q install jedi omegaconf hydra-core icecream pyrsistent")
    os.system(
        "pip install dgl==1.0.2+cu116 -f https://data.dgl.ai/wheels/cu116/repo.html"
    )
    os.system(
        "cd RFdiffusion/env/SE3Transformer; pip -q install --no-cache-dir -r requirements.txt; pip -q install ."
    )
    os.system("wget -qnc https://files.ipd.uw.edu/krypton/ananas")
    os.system("chmod +x ananas")

if not os.path.isdir("colabdesign"):
    print("installing ColabDesign...")
    os.system("pip -q install git+https://github.com/sokrypton/ColabDesign.git")
    os.system("ln -s /usr/local/lib/python3.*/dist-packages/colabdesign colabdesign")

if not os.path.isdir("RFdiffusion/models"):
    print("downloading RFdiffusion params...")
    os.system("mkdir RFdiffusion/models")
    models = ["Base_ckpt.pt", "Complex_base_ckpt.pt", "Complex_beta_ckpt.pt"]
    for m in models:
        while os.path.isfile(f"{m}.aria2"):
            time.sleep(5)
    os.system(f"mv {' '.join(models)} RFdiffusion/models")
    os.system("unzip schedules.zip; rm schedules.zip")

if "RFdiffusion" not in sys.path:
    os.environ["DGLBACKEND"] = "pytorch"
    sys.path.append("RFdiffusion")


# third party imports
import hydra
from hydra import compose, initialize
from hydra.core.hydra_config import HydraConfig
from omegaconf import DictConfig, OmegaConf
from inference.utils import parse_pdb
from colabdesign.rf.utils import get_ca
from colabdesign.rf.utils import fix_contigs, fix_partial_contigs, fix_pdb, sym_it
from colabdesign.shared.protein import pdb_to_string
from colabdesign.shared.plot import plot_pseudo_3D
from Bio import PDB

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


def extract_deepest_keys(nested_dict, current_path="", deepest_keys=None):
    if deepest_keys is None:
        deepest_keys = {}

    for key, value in nested_dict.items():
        path = f"{current_path}.{key}" if current_path else key

        if isinstance(value, dict):
            extract_deepest_keys(value, path, deepest_keys)
        else:
            deepest_keys[path] = value

    return deepest_keys


def add_deepest_keys_to_dataframe(deepest_keys_values, df_results):
    for path, value in deepest_keys_values.items():
        # Handle special cases for certain keys
        if (
            "target_protein_directory" in path
            or "binder_protein_template_directory" in path
        ):
            value = os.path.splitext(os.path.basename(value))[0]

        # Add the values as a new column in df_results
        df_results[path] = value

    return df_results

def enrich_and_collect(multirun_path, path, cfg):
    # Find the scores file in multirun_path
    results_csv_path = None
    for file_name in os.listdir(f"{multirun_path}/"):
        # if file_name.endswith("_scores.csv"):
        if file_name.endswith(f"{path}_scores.csv"):
            results_csv_path = os.path.join(multirun_path, file_name)
            break

    # Check if the scores.csv file is found
    if results_csv_path is None:
        print("No scores file found in the specified directory.")
        return

    # Read the scores csv into a DataFrame
    df_new_results = pd.read_csv(results_csv_path)

    # Extract and add the deepest level keys and values to df_new_results
    deepest_keys_values = extract_deepest_keys(OmegaConf.to_container(cfg))
    df_new_results = add_deepest_keys_to_dataframe(deepest_keys_values, df_new_results)

    # Check if the existing output CSV file exists
    # output_csv_path = os.path.join(multirun_path, f"{path}_samples_table.csv")
    output_csv_path = os.path.join(multirun_path, f"summary_samples_table.csv")
    if os.path.exists(output_csv_path):
        # Read existing data
        df_existing_results = pd.read_csv(output_csv_path)
        # Concatenate new data with existing data
        df_combined_results = pd.concat([df_existing_results, df_new_results], ignore_index=True)
    else:
        df_combined_results = df_new_results

    # Write the combined data to the CSV file
    df_combined_results.to_csv(output_csv_path, index=False)
    print("Enriched results and conditions tables written to", results_csv_path)

def create_conditions_table(multirun_path, path, cfg, user_inputs):
    # Read the combined results CSV into a DataFrame
    # combined_csv_path = os.path.join(multirun_path, f"{path}_samples_table.csv")
    combined_csv_path = os.path.join(multirun_path, f"summary_samples_table.csv")
    df_combined_results = pd.read_csv(combined_csv_path)

    # Columns to process
    columns = ['plddt', 'ptm', 'pae', 'rmsd', 'affinity']

    # Creating a dictionary to store min and max values
    condensed_data = {}
    for col in columns:
        condensed_data[f"{col}_min"] = [df_combined_results[col].min()]
        condensed_data[f"{col}_max"] = [df_combined_results[col].max()]

    # Create a new DataFrame for condensed results
    df_condensed = pd.DataFrame(condensed_data)

    # Insert the n_samples column as the first column
    df_condensed.insert(0, 'n_designs', len(df_combined_results))

    # Write the condensed data to the CSV file
    # condensed_csv_path = os.path.join(multirun_path, f"{path}_conditions_table.csv")
    condensed_csv_path = os.path.join(multirun_path, f"summary_conditions_table.csv")
    df_condensed.to_csv(condensed_csv_path, index=False)
    print("Conditions table written to", condensed_csv_path)


def get_files_from_directory(root_dir, extension, max_depth=3):
    pdb_files = []

    for root, dirs, files in os.walk(root_dir):
        depth = root[len(root_dir) :].count(os.path.sep)

        if depth <= max_depth:
            for f in files:
                if f.endswith(extension):
                    pdb_files.append(os.path.join(root, f))

            # Prune the directory list if we are at max_depth
            if depth == max_depth:
                del dirs[:]
    print(
        "Found {} files with extension {} in directory {}".format(
            len(pdb_files), extension, root_dir
        )
    )
    return pdb_files


def get_pdb(pdb_code=None):
    print("Getting PDB...", pdb_code)
    if pdb_code is None or pdb_code == "":
        # upload_dict = files.upload()
        # pdb_string = upload_dict[list(upload_dict.keys())[0]]
        # with open("tmp.pdb","wb") as out: out.write(pdb_string)
        print("Warning: no target pdb file.")
        return "tmp.pdb"
    elif os.path.isfile(pdb_code):
        return pdb_code
    elif len(pdb_code) == 4:
        if not os.path.isfile(f"{pdb_code}.pdb1"):
            os.system(f"wget -qnc https://files.rcsb.org/download/{pdb_code}.pdb1.gz")
            os.system(f"gunzip {pdb_code}.pdb1.gz")
        return f"{pdb_code}.pdb1"
    else:
        os.system(
            f"wget -qnc https://alphafold.ebi.ac.uk/files/AF-{pdb_code}-F1-model_v3.pdb"
        )
        return f"AF-{pdb_code}-F1-model_v3.pdb"


def run_ananas(pdb_str, path, outputs_directory, sym=None):
    print("Running AnAnaS...")
    pdb_filename = f"{outputs_directory}/{path}/ananas_input.pdb"
    out_filename = f"{outputs_directory}/{path}/ananas.json"
    with open(pdb_filename, "w") as handle:
        handle.write(pdb_str)

    cmd = f"./ananas {pdb_filename} -u -j {out_filename}"
    if sym is None:
        os.system(cmd)
    else:
        os.system(f"{cmd} {sym}")

    # parse results
    try:
        out = json.loads(open(out_filename, "r").read())
        results, AU = out[0], out[-1]["AU"]
        group = AU["group"]
        chains = AU["chain names"]
        rmsd = results["Average_RMSD"]
        print(f"AnAnaS detected {group} symmetry at RMSD:{rmsd:.3}")

        C = np.array(results["transforms"][0]["CENTER"])
        A = [np.array(t["AXIS"]) for t in results["transforms"]]

        # apply symmetry and filter to the asymmetric unit
        new_lines = []
        for line in pdb_str.split("\n"):
            if line.startswith("ATOM"):
                chain = line[21:22]
                if chain in chains:
                    x = np.array([float(line[i : (i + 8)]) for i in [30, 38, 46]])
                    if group[0] == "c":
                        x = sym_it(x, C, A[0])
                    if group[0] == "d":
                        x = sym_it(x, C, A[1], A[0])
                    coord_str = "".join(["{:8.3f}".format(a) for a in x])
                    new_lines.append(line[:30] + coord_str + line[54:])
            else:
                new_lines.append(line)
        return results, "\n".join(new_lines)

    except:
        return None, pdb_str


def run_diffusion(
    contigs,
    path,
    pdb=None,
    iterations=50,
    symmetry="none",
    order=1,
    hotspot=None,
    chains=None,
    add_potential=False,
    num_designs=1,
    use_beta_model=False,
    visual="none",
    outputs_directory="outputs",
):
    print("Running diffusion with contigs:", contigs, "and path:", path)
    full_path = f"{outputs_directory}/{path}"
    os.makedirs(full_path, exist_ok=True)
    opts = [
        f"inference.output_prefix={full_path}",
        f"inference.num_designs={num_designs}",
    ]

    if chains == "":
        chains = None

    # determine symmetry type
    if symmetry in ["auto", "cyclic", "dihedral"]:
        if symmetry == "auto":
            sym, copies = None, 1
        else:
            sym, copies = {
                "cyclic": (f"c{order}", order),
                "dihedral": (f"d{order}", order * 2),
            }[symmetry]
    else:
        symmetry = None
        sym, copies = None, 1

    # determine mode
    contigs = contigs.replace(",", " ").replace(":", " ").split()
    is_fixed, is_free = False, False
    fixed_chains = []
    for contig in contigs:
        for x in contig.split("/"):
            a = x.split("-")[0]
            if a[0].isalpha():
                is_fixed = True
                if a[0] not in fixed_chains:
                    fixed_chains.append(a[0])
            if a.isnumeric():
                is_free = True
    if len(contigs) == 0 or not is_free:
        mode = "partial"
    elif is_fixed:
        mode = "fixed"
    else:
        mode = "free"

    # fix input contigs
    if mode in ["partial", "fixed"]:
        pdb_str = pdb_to_string(get_pdb(pdb), chains=chains)
        if symmetry == "auto":
            a, pdb_str = run_ananas(pdb_str, path, outputs_directory)
            if a is None:
                print(f"ERROR: no symmetry detected")
                symmetry = None
                sym, copies = None, 1
            else:
                if a["group"][0] == "c":
                    symmetry = "cyclic"
                    sym, copies = a["group"], int(a["group"][1:])
                elif a["group"][0] == "d":
                    symmetry = "dihedral"
                    sym, copies = a["group"], 2 * int(a["group"][1:])
                else:
                    print(
                        f'ERROR: the detected symmetry ({a["group"]}) not currently supported'
                    )
                    symmetry = None
                    sym, copies = None, 1

        elif mode == "fixed":
            pdb_str = pdb_to_string(pdb_str, chains=fixed_chains)

        pdb_filename = f"{full_path}/input.pdb"
        with open(pdb_filename, "w") as handle:
            handle.write(pdb_str)

        parsed_pdb = parse_pdb(pdb_filename)
        opts.append(f"inference.input_pdb={pdb_filename}")
        if mode in ["partial"]:
            iterations = int(80 * (iterations / 200))
            opts.append(f"diffuser.partial_T={iterations}")
            contigs = fix_partial_contigs(contigs, parsed_pdb)
        else:
            opts.append(f"diffuser.T={iterations}")
            contigs = fix_contigs(contigs, parsed_pdb)
    else:
        opts.append(f"diffuser.T={iterations}")
        parsed_pdb = None
        contigs = fix_contigs(contigs, parsed_pdb)

    if hotspot is not None and hotspot != "":
        opts.append(f"ppi.hotspot_res=[{hotspot}]")

    # setup symmetry
    if sym is not None:
        sym_opts = ["--config-name symmetry", f"inference.symmetry={sym}"]
        if add_potential:
            sym_opts += [
                "'potentials.guiding_potentials=[\"type:olig_contacts,weight_intra:1,weight_inter:0.1\"]'",
                "potentials.olig_intra_all=True",
                "potentials.olig_inter_all=True",
                "potentials.guide_scale=2",
                "potentials.guide_decay=quadratic",
            ]
        opts = sym_opts + opts
        contigs = sum([contigs] * copies, [])

    opts.append(f"'contigmap.contigs=[{' '.join(contigs)}]'")
    opts += ["inference.dump_pdb=True", "inference.dump_pdb_path='/dev/shm'"]
    if use_beta_model:
        opts += [
            "inference.ckpt_override_path=./RFdiffusion/models/Complex_beta_ckpt.pt"
        ]

    print("mode:", mode)
    print("output:", full_path)
    print("contigs:", contigs)

    opts_str = " ".join(opts)
    command = f"python -u ./RFdiffusion/run_inference.py {opts_str}"
    print(command)

    with subprocess.Popen(
        command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True
    ) as process:
        exit_code = 0
        try:
            # Read output line by line as it is produced
            for line in process.stdout:
                print(line.strip())

            # Wait for the subprocess to finish and get the exit code
            process.communicate()
            exit_code = process.returncode
            print(f"Finished script with return code {exit_code}")

            if exit_code != 0:
                # If subprocess failed, log the stderr
                error_message = process.stderr.read()
                print(f"Command failed with exit code {exit_code}")
                print(f"Error message: {error_message}")

        except Exception as e:
            print(f"Error while running command: {e}")
            exit(exit_code)

    print("Done with the run RFDiffusion script call")

    # fix pdbs
    for n in range(num_designs):
        pdbs = [
            f"{outputs_directory}/traj/{path}_{n}_pX0_traj.pdb",
            f"{outputs_directory}/traj/{path}_{n}_Xt-1_traj.pdb",
            f"{full_path}_{n}.pdb",
        ]
        for pdb in pdbs:
            with open(pdb, "r") as handle:
                pdb_str = handle.read()
            with open(pdb, "w") as handle:
                handle.write(fix_pdb(pdb_str, contigs))

    return contigs, copies


def prodigy_run(csv_path, pdb_path):
    df = pd.read_csv(csv_path)
    for i, r in df.iterrows():
        design = r["design"]
        n = r["n"]
        file_path = f"{pdb_path}/design{design}_n{n}.pdb"
        try:
            subprocess.run(
                ["prodigy", "-q", file_path], stdout=open("temp.txt", "w"), check=True
            )
            with open("temp.txt", "r") as f:
                lines = f.readlines()
                if lines:  # Check if lines is not empty
                    affinity = float(lines[0].split(" ")[-1].split("/")[0])
                    df.loc[i, "affinity"] = affinity
                else:
                    # print(f"No output from prodigy for {r['path']}")
                    print(f"No output from prodigy for {file_path}")
                    # Handle the case where prodigy did not produce output
        except subprocess.CalledProcessError:
            # print(f"Prodigy command failed for {r['path']}")
            print(f"Prodigy command failed for {file_path}")

    # export results
    df.to_csv(f"{csv_path}", index=None)


@hydra.main(version_base=None, config_path="conf", config_name="config")
def my_app(cfg: DictConfig) -> None:
    user_inputs = get_plex_job_inputs()
    print(f"user inputs from plex: {user_inputs}")

    # Override Hydra default params with user supplied params
    OmegaConf.update(cfg, "params.basic_settings.binder_length", "", merge=False)
    OmegaConf.update(cfg, "params.advanced_settings.hotspot", "", merge=False)
    OmegaConf.update(cfg, "params.basic_settings.num_designs", user_inputs["number_of_binder_designs"], merge=False)
    OmegaConf.update(cfg, "params.basic_settings.pdb_chain", user_inputs["target_chain"], merge=False)
    # OmegaConf.update(cfg, "params.advanced_settings.pdb_start_residue", user_inputs["target_start_residue"], merge=False)
    # OmegaConf.update(cfg, "params.advanced_settings.pdb_end_residue", user_inputs["target_end_residue"], merge=False)
    OmegaConf.update(cfg, "params.expert_settings.RFDiffusion_Binder.contigs_override", "", merge=False)

    print(OmegaConf.to_yaml(cfg))
    print(f"Working directory : {os.getcwd()}")

    # defining output directory
    if cfg.outputs.directory is None:
        outputs_directory = hydra.core.hydra_config.HydraConfig.get().runtime.output_dir
    else:
        outputs_directory = cfg.outputs.directory
    print(f"Output directory : {outputs_directory}")

    # defining input files
    if user_inputs.get("protein_complex"):
        input_target_path = user_inputs["protein_complex"]
    else:
        input_target_path = get_files_from_directory(cfg.inputs.target_directory, ".pdb")

        if cfg.inputs.target_pattern is not None:
            input_target_path = [
                file for file in input_target_path if cfg.inputs.target_pattern in file
            ]

    if not isinstance(input_target_path, list):
        input_target_path = [input_target_path]
    print("Identified Targets : ", input_target_path)

    first_residue_target_chain, last_residue_target_chain = find_chain_residue_range(input_target_path[0], user_inputs["target_chain"])
    if isinstance(user_inputs["target_start_residue"], int) and user_inputs["target_start_residue"] > 0:
        OmegaConf.update(cfg, "params.advanced_settings.pdb_start_residue", user_inputs["target_start_residue"], merge=False)
    else: OmegaConf.update(cfg, "params.advanced_settings.pdb_start_residue", first_residue_target_chain, merge=False)
    
    if isinstance(user_inputs["target_end_residue"], int) and user_inputs["target_end_residue"] > 0:
        OmegaConf.update(cfg, "params.advanced_settings.pdb_end_residue", user_inputs["target_end_residue"], merge=False)
    else: OmegaConf.update(cfg, "params.advanced_settings.pdb_end_residue", last_residue_target_chain, merge=False)


    # running design for every input target file
    for target_path in input_target_path:
        start_time = time.time()

        name = cfg.params.basic_settings.experiment_name
        pdb = target_path  # cfg.basic_settings.pdb
        hotspot = cfg.params.advanced_settings.hotspot.replace(" ", "")
        iterations = cfg.params.expert_settings.RFDiffusion_Binder.iterations
        num_designs = cfg.params.basic_settings.num_designs
        use_beta_model = cfg.params.advanced_settings.use_beta_model
        visual = cfg.params.expert_settings.RFDiffusion_Binder.visual

        # symmetry settings
        symmetry = cfg.params.expert_settings.RFDiffusion_Symmetry.symmetry
        order = cfg.params.expert_settings.RFDiffusion_Symmetry.order
        chains = cfg.params.expert_settings.RFDiffusion_Symmetry.chains
        add_potential = cfg.params.expert_settings.RFDiffusion_Symmetry.add_potential

        # contig assembly
        ##TODO: simplify load contig parameters
        binder_length = cfg.params.basic_settings.binder_length
        pdb_chain = cfg.params.basic_settings.pdb_chain
        pdb_start_residue = cfg.params.advanced_settings.pdb_start_residue
        pdb_end_residue = cfg.params.advanced_settings.pdb_end_residue
        min_binder_length = cfg.params.advanced_settings.min_binder_length
        max_binder_length = cfg.params.advanced_settings.max_binder_length
        contigs_override = (
            cfg.params.expert_settings.RFDiffusion_Binder.contigs_override
        )

        binder_chain = user_inputs["binder_chain"]
        target_chain = pdb_chain
        cutoff = user_inputs["cutoff"] # distance to define inter-protein contacts (in Angstrom)
        n_samples = user_inputs["n_prompts"] # total number of prompts generated

        p_masking_contact_domain = user_inputs["p_contact_domain_masking"] # probability of masking a contact domain
        p_masking_noncontact_domain = user_inputs["p_noncontact_domain_masking"] # probability of masking a non-contact domain
        p_masking_contact_domain = float(p_masking_contact_domain)
        p_masking_noncontact_domain = float(p_masking_noncontact_domain)
        if (0.0 <= p_masking_contact_domain <= 1.0) and (0.0 <= p_masking_noncontact_domain <= 1.0): # Check if the the p's are within the interval [0.0, 1.0]
            pass
        else:
            raise ValueError(f"p_masking_contact or p_masking_noncontact is not in the interval [0.0, 1.0].")
        
        domain_distance_threshold = user_inputs["domain_distance_threshold"] # definition of constitutes separate domains (in units of residues)
        full_prompts = prompt_generator.generate_full_prompts(target_path, binder_chain, target_chain, pdb_start_residue, pdb_end_residue, cutoff, n_samples, p_masking_contact_domain, p_masking_noncontact_domain, domain_distance_threshold)
        print(full_prompts)

        prompt_counter = 0
        for prompt in full_prompts:

            contigs_override = prompt
            cfg.params.expert_settings.RFDiffusion_Binder.contigs_override = contigs_override
            prompt_counter += 1
            print("promt: ", contigs_override)
            
            ## binder length
            if min_binder_length != None and max_binder_length != None:
                binder_length_constructed = (
                    str(min_binder_length) + "-" + str(max_binder_length)
                )
            else:
                binder_length_constructed = str(binder_length) + "-" + str(binder_length)

            ## residue start
            if pdb_start_residue != None and pdb_end_residue != None:
                residue_constructed = str(pdb_start_residue) + "-" + str(pdb_end_residue)
            else:
                residue_constructed = ""

            ## contig assembly
            contigs_constructed = (
                pdb_chain + residue_constructed + "/0: " + binder_length_constructed
            )
            if contigs_override == "":
                contigs = contigs_constructed
            else:
                contigs = contigs_override

            # determine where to save
            path = name
            while os.path.exists(f"{outputs_directory}/{path}_0.pdb"):
                path = (
                    name
                    + "_"
                    + "".join(random.choices(string.ascii_lowercase + string.digits, k=5))
                )

            flags = {
                "contigs": contigs,
                "pdb": pdb,
                "order": order,
                "iterations": iterations,
                "symmetry": symmetry,
                "hotspot": hotspot,
                "path": path,
                "chains": chains,
                "add_potential": add_potential,
                "num_designs": num_designs,
                "use_beta_model": use_beta_model,
                "visual": visual,
                "outputs_directory": outputs_directory,
            }

            for k, v in flags.items():
                if isinstance(v, str):
                    flags[k] = v.replace("'", "").replace('"', "")

            contigs, copies = run_diffusion(**flags)

            num_seqs = cfg.params.expert_settings.ProteinMPNN.num_seqs
            initial_guess = cfg.params.expert_settings.ProteinMPNN.initial_guess
            num_recycles = cfg.params.expert_settings.Alphafold.num_recycles
            use_multimer = cfg.params.expert_settings.Alphafold.use_multimer
            rm_aa = cfg.params.expert_settings.ProteinMPNN.rm_aa
            mpnn_sampling_temp = cfg.params.expert_settings.ProteinMPNN.mpnn_sampling_temp
            use_solubleMPNN = cfg.params.expert_settings.ProteinMPNN.use_solubleMPNN

            if not os.path.isfile("params/done.txt"):
                print("downloading AlphaFold params...")
                while not os.path.isfile("params/done.txt"):
                    time.sleep(5)

            contigs_str = ":".join(contigs)
            opts = [
                f"--pdb={outputs_directory}/{path}_0.pdb",
                f"--loc={outputs_directory}/{path}",
                f"--contig={contigs_str}",
                f"--copies={copies}",
                f"--num_seqs={num_seqs}",
                f"--num_recycles={num_recycles}",
                f"--rm_aa={rm_aa}",
                f"--mpnn_sampling_temp={mpnn_sampling_temp}",
                f"--num_designs={num_designs}",
            ]
            if initial_guess:
                opts.append("--initial_guess")
            if use_multimer:
                opts.append("--use_multimer")
            if use_solubleMPNN:
                opts.append("--use_soluble")
            opts = " ".join(opts)

            command_design = f"python -u colabdesign/rf/designability_test.py {opts}"
            os.system(command_design)

            print("running Prodigy")
            prodigy_run(
                f"{outputs_directory}/{path}/mpnn_results.csv",
                f"{outputs_directory}/{path}/all_pdb",
            )

            command_mv = f"mkdir {outputs_directory}/{path}/traj && mv {outputs_directory}/traj/{path}* {outputs_directory}/{path}/traj && mv {outputs_directory}/{path}_* {outputs_directory}/{path}"
            command_zip = f"zip -r {path}.result.zip {outputs_directory}/{path}*"
            command_collect = f"mv {path}.result.zip /{outputs_directory} && mv {outputs_directory}/{path}/best.pdb /{outputs_directory}/{path}_prompt{prompt_counter}_best.pdb && mv {outputs_directory}/{path}/mpnn_results.csv /{outputs_directory}/{path}_scores.csv"
            os.system(command_mv)
            os.system(command_zip)
            os.system(command_collect)

            # enrich and summarise run and results information and write to csv file
            print("running enricher")
            # enricher(outputs_directory, cfg)
            enrich_and_collect(outputs_directory, path, cfg)
        
        create_conditions_table(outputs_directory, path, cfg, user_inputs)

        print("design complete...")
        end_time = time.time()
        duration = end_time - start_time
        print(f"executed in {duration:.2f} seconds.")


if __name__ == "__main__":
    my_app()


# #OLD CODE
# # reference_protein_complex = 'summary/UROK_HUMAN_1-133.pdb'
# binder_chain = 'B' # 'B'
# target_chain = pdb_chain # 'A'
# cutoff = 5.0 # distance to define inter-protein contacts (in Angstrom)
# n_samples = 1 # total number of prompts generated
# p_masking_contact_domain = .6 # probability of masking a contact domain
# p_masking_noncontact_domain = 0.1 # probability of masking a non-contact domain
# domain_distance_threshold = 6 # definition of constitutes separate domains (in units of residues)
