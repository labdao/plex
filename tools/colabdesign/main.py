# standard library imports
import glob
import os
import time
import signal
import sys
import random
import string
import re
import json
import numpy as np
import matplotlib.pyplot as plt
import ipywidgets as widgets
import py3Dmol
import yaml
from google.colab import files
from IPython.display import display, HTML

print("Starting main.py...")

# setup
if not os.path.isdir("params"):
  os.system("apt-get install aria2")
  os.system("mkdir params")
  # send param download into background
  os.system("(\
  aria2c -q -x 16 https://files.ipd.uw.edu/krypton/schedules.zip; \
  aria2c -q -x 16 http://files.ipd.uw.edu/pub/RFdiffusion/6f5902ac237024bdd0c176cb93063dc4/Base_ckpt.pt; \
  aria2c -q -x 16 http://files.ipd.uw.edu/pub/RFdiffusion/e29311f6f1bf1af907f9ef9f44b8328b/Complex_base_ckpt.pt; \
  aria2c -q -x 16 http://files.ipd.uw.edu/pub/RFdiffusion/f572d396fae9206628714fb2ce00f72e/Complex_beta_ckpt.pt; \
  aria2c -q -x 16 https://storage.googleapis.com/alphafold/alphafold_params_2022-12-06.tar; \
  tar -xf alphafold_params_2022-12-06.tar -C params; \
  touch params/done.txt) &")

if not os.path.isdir("RFdiffusion"):
  print("installing RFdiffusion...")
  os.system("git clone https://github.com/sokrypton/RFdiffusion.git")
  os.system("pip -q install jedi omegaconf hydra-core icecream pyrsistent")
  os.system("pip install dgl==1.0.2+cu116 -f https://data.dgl.ai/wheels/cu116/repo.html")
  os.system("cd RFdiffusion/env/SE3Transformer; pip -q install --no-cache-dir -r requirements.txt; pip -q install .")
  os.system("wget -qnc https://files.ipd.uw.edu/krypton/ananas")
  os.system("chmod +x ananas")

if not os.path.isdir("colabdesign"):
  print("installing ColabDesign...")
  os.system("pip -q install git+https://github.com/sokrypton/ColabDesign.git")
  os.system("ln -s /usr/local/lib/python3.*/dist-packages/colabdesign colabdesign")

if not os.path.isdir("RFdiffusion/models"):
  print("downloading RFdiffusion params...")
  os.system("mkdir RFdiffusion/models")
  models = ["Base_ckpt.pt","Complex_base_ckpt.pt","Complex_beta_ckpt.pt"]
  for m in models:
    while os.path.isfile(f"{m}.aria2"):
      time.sleep(5)
  os.system(f"mv {' '.join(models)} RFdiffusion/models")
  os.system("unzip schedules.zip; rm schedules.zip")

if 'RFdiffusion' not in sys.path:
  os.environ["DGLBACKEND"] = "pytorch"
  sys.path.append('RFdiffusion')

import subprocess
import pandas as pd

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
        if 'target_protein_directory' in path or 'binder_protein_template_directory' in path:
            value = os.path.splitext(os.path.basename(value))[0]

        # Add the values as a new column in df_results
        df_results[path] = value
    
    return df_results

import yaml

def enricher(multirun_path, cfg):

    # Find the scores file in multirun_path
    results_csv_path = None
    for file_name in os.listdir(f"{multirun_path}/"):
        if file_name.endswith("_scores.csv"):
            results_csv_path = os.path.join(multirun_path, file_name)
            break

    # Check if the mpnn_results.csv file is found
    if results_csv_path is None:
        print("No scores file found in the specified directory.")
        return

    # Read the scores csv into a DataFrame
    df_results = pd.read_csv(results_csv_path)

    # Extract and add the deepest level keys and values to df_results
    deepest_keys_values = extract_deepest_keys(OmegaConf.to_container(cfg))

    df_results = add_deepest_keys_to_dataframe(deepest_keys_values, df_results)

    # print('enriched results', df_results)

    df_results.to_csv(results_csv_path, index=False)

def get_files_from_directory(root_dir, extension, max_depth=3):
    pdb_files = []
    
    for root, dirs, files in os.walk(root_dir):
        depth = root[len(root_dir):].count(os.path.sep)
        
        if depth <= max_depth:
            for f in files:
                if f.endswith(extension):
                    pdb_files.append(os.path.join(root, f))
                    
            # Prune the directory list if we are at max_depth
            if depth == max_depth:
                del dirs[:]
    print("Found {} files with extension {} in directory {}".format(len(pdb_files), extension, root_dir))
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
    os.system(f"wget -qnc https://alphafold.ebi.ac.uk/files/AF-{pdb_code}-F1-model_v3.pdb")
    return f"AF-{pdb_code}-F1-model_v3.pdb"

def run_ananas(pdb_str, path, outputs_directory, sym=None):
  print("Running AnAnaS...")
  pdb_filename = f"{outputs_directory}/{path}/ananas_input.pdb"
  out_filename = f"{outputs_directory}/{path}/ananas.json"
  with open(pdb_filename,"w") as handle:
    handle.write(pdb_str)

  cmd = f"./ananas {pdb_filename} -u -j {out_filename}"
  if sym is None: os.system(cmd)
  else: os.system(f"{cmd} {sym}")

  # parse results
  try:
    out = json.loads(open(out_filename,"r").read())
    results,AU = out[0], out[-1]["AU"]
    group = AU["group"]
    chains = AU["chain names"]
    rmsd = results["Average_RMSD"]
    print(f"AnAnaS detected {group} symmetry at RMSD:{rmsd:.3}")

    C = np.array(results['transforms'][0]['CENTER'])
    A = [np.array(t["AXIS"]) for t in results['transforms']]

    # apply symmetry and filter to the asymmetric unit
    new_lines = []
    for line in pdb_str.split("\n"):
      if line.startswith("ATOM"):
        chain = line[21:22]
        if chain in chains:
          x = np.array([float(line[i:(i+8)]) for i in [30,38,46]])
          if group[0] == "c":
            x = sym_it(x,C,A[0])
          if group[0] == "d":
            x = sym_it(x,C,A[1],A[0])
          coord_str = "".join(["{:8.3f}".format(a) for a in x])
          new_lines.append(line[:30]+coord_str+line[54:])
      else:
        new_lines.append(line)
    return results, "\n".join(new_lines)

  except:
    return None, pdb_str

def run(command, steps, num_designs=1, visual="none"):
  print("Running command:", command)

  def run_command_and_get_pid(command):
    pid_file = '/dev/shm/pid'
    os.system(f'nohup {command} & echo $! > {pid_file}')
    with open(pid_file, 'r') as f:
      pid = int(f.read().strip())
    os.remove(pid_file)
    return pid
  
  def is_process_running(pid):
    try:
      os.kill(pid, 0)
    except OSError:
      return False
    else:
      return True

  run_output = widgets.Output()
  progress = widgets.FloatProgress(min=0, max=1, description='running', bar_style='info')
  display(widgets.VBox([progress, run_output]))

  # clear previous run
  for n in range(steps):
    if os.path.isfile(f"/dev/shm/{n}.pdb"):
      os.remove(f"/dev/shm/{n}.pdb")

  pid = run_command_and_get_pid(command)
  try:
    fail = False
    for _ in range(num_designs):

      # for each step check if output generated
      for n in range(steps):
        wait = True
        while wait and not fail:
          time.sleep(0.1)
          if os.path.isfile(f"/dev/shm/{n}.pdb"):
            pdb_str = open(f"/dev/shm/{n}.pdb").read()
            if pdb_str[-3:] == "TER":
              wait = False
            elif not is_process_running(pid):
              fail = True
          elif not is_process_running(pid):
            fail = True

        if fail:
          progress.bar_style = 'danger'
          progress.description = "failed"
          break

        else:
          progress.value = (n+1) / steps
          if visual != "none":
            with run_output:
              run_output.clear_output(wait=True)
              if visual == "image":
                xyz, bfact = get_ca(f"/dev/shm/{n}.pdb", get_bfact=True)
                fig = plt.figure()
                fig.set_dpi(100);fig.set_figwidth(6);fig.set_figheight(6)
                ax1 = fig.add_subplot(111);ax1.set_xticks([]);ax1.set_yticks([])
                plot_pseudo_3D(xyz, c=bfact, cmin=0.5, cmax=0.9, ax=ax1)
                plt.show()
              if visual == "interactive":
                view = py3Dmol.view(js='https://3dmol.org/build/3Dmol.js')
                view.addModel(pdb_str,'pdb')
                view.setStyle({'cartoon': {'colorscheme': {'prop':'b','gradient': 'roygb','min':0.5,'max':0.9}}})
                view.zoomTo()
                view.show()
        if os.path.exists(f"/dev/shm/{n}.pdb"):
          os.remove(f"/dev/shm/{n}.pdb")
      if fail:
        progress.bar_style = 'danger'
        progress.description = "failed"
        break

    while is_process_running(pid):
      time.sleep(0.1)

  except KeyboardInterrupt:
    os.kill(pid, signal.SIGTERM)
    progress.bar_style = 'danger'
    progress.description = "stopped"

def run_diffusion(contigs, path, pdb=None, iterations=50,
                  symmetry="none", order=1, hotspot=None,
                  chains=None, add_potential=False,
                  num_designs=1, use_beta_model=False, visual="none", outputs_directory="outputs"):

  print("Running diffusion with contigs:", contigs, "and path:", path)
  full_path = f"{outputs_directory}/{path}"
  os.makedirs(full_path, exist_ok=True)
  opts = [f"inference.output_prefix={full_path}",
          f"inference.num_designs={num_designs}"]

  if chains == "": chains = None

  # determine symmetry type
  if symmetry in ["auto","cyclic","dihedral"]:
    if symmetry == "auto":
      sym, copies = None, 1
    else:
      sym, copies = {"cyclic":(f"c{order}",order),
                      "dihedral":(f"d{order}",order*2)}[symmetry]
  else:
    symmetry = None
    sym, copies = None, 1

  # determine mode
  contigs = contigs.replace(","," ").replace(":"," ").split()
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
  if mode in ["partial","fixed"]:
    pdb_str = pdb_to_string(get_pdb(pdb), chains=chains)
    if symmetry == "auto":
      a, pdb_str = run_ananas(pdb_str, path, outputs_directory)
      if a is None:
        print(f'ERROR: no symmetry detected')
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
          print(f'ERROR: the detected symmetry ({a["group"]}) not currently supported')
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
      sym_opts += ["'potentials.guiding_potentials=[\"type:olig_contacts,weight_intra:1,weight_inter:0.1\"]'",
                    "potentials.olig_intra_all=True","potentials.olig_inter_all=True",
                    "potentials.guide_scale=2","potentials.guide_decay=quadratic"]
    opts = sym_opts + opts
    contigs = sum([contigs] * copies,[])

  opts.append(f"'contigmap.contigs=[{' '.join(contigs)}]'")
  opts += ["inference.dump_pdb=True","inference.dump_pdb_path='/dev/shm'"]
  if use_beta_model:
    opts += ["inference.ckpt_override_path=./RFdiffusion/models/Complex_beta_ckpt.pt"]

  print("mode:", mode)
  print("output:", full_path)
  print("contigs:", contigs)

  opts_str = " ".join(opts)
  cmd = f"./RFdiffusion/run_inference.py {opts_str}"
  print(cmd)

  # RUN
  run(cmd, iterations, num_designs, visual=visual)

  # fix pdbs
  for n in range(num_designs):
    pdbs = [f"{outputs_directory}/traj/{path}_{n}_pX0_traj.pdb",
            f"{outputs_directory}/traj/{path}_{n}_Xt-1_traj.pdb",
            f"{full_path}_{n}.pdb"]
    for pdb in pdbs:
      with open(pdb,"r") as handle: pdb_str = handle.read()
      with open(pdb,"w") as handle: handle.write(fix_pdb(pdb_str, contigs))

  return contigs, copies

def prodigy_run(csv_path, pdb_path):
    df = pd.read_csv(csv_path)
    for i,r in df.iterrows():
        design = r['design']
        n = r['n']
        file_path = f"{pdb_path}/design{design}_n{n}.pdb"
        try:
            subprocess.run(["prodigy", "-q", file_path], stdout=open('temp.txt', 'w'), check=True)
            with open('temp.txt', 'r') as f:
                lines = f.readlines()
                if lines:  # Check if lines is not empty
                    affinity = float(lines[0].split(' ')[-1].split('/')[0])
                    df.loc[i,'affinity'] = affinity
                else:
                    # print(f"No output from prodigy for {r['path']}")
                    print(f"No output from prodigy for {file_path}")
                    # Handle the case where prodigy did not produce output
        except subprocess.CalledProcessError:
            # print(f"Prodigy command failed for {r['path']}")
            print(f"Prodigy command failed for {file_path}")

    # export results
    df.to_csv(f"{csv_path}",index=None)

@hydra.main(version_base=None, config_path="conf", config_name="config")
def my_app(cfg : DictConfig) -> None:
    print(OmegaConf.to_yaml(cfg))
    print(f"Working directory : {os.getcwd()}")

    # defining output directory
    if cfg.outputs.directory is None:
        outputs_directory = hydra.core.hydra_config.HydraConfig.get().runtime.output_dir
    else:
        outputs_directory = cfg.outputs.directory
    print(f"Output directory  : {outputs_directory}")

    # defining input files
    input_target_path = get_files_from_directory(cfg.inputs.target_directory, ".pdb")
    if cfg.inputs.target_pattern is not None:
        input_target_path = [file for file in input_target_path if cfg.inputs.target_pattern in file]
    if not isinstance(input_target_path, list):
        input_target_path = [input_target_path]
    print("Identified Targets : ", input_target_path)
    
    # running design for every input target file
    for target_path in input_target_path:
        # for binder_path in input_binder_path:
        
            
        start_time = time.time()

        name = cfg.params.basic_settings.experiment_name
        pdb = target_path # cfg.basic_settings.pdb
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
        contigs_override = cfg.params.expert_settings.RFDiffusion_Binder.contigs_override

        ## binder length
        if min_binder_length != None and max_binder_length != None:
            binder_length_constructed = str(min_binder_length) + "-" + str(max_binder_length)
        else:
            binder_length_constructed = str(binder_length) + "-" + str(binder_length)

        ## residue start
        if pdb_start_residue != None and pdb_end_residue != None:
            residue_constructed = str(pdb_start_residue) + "-" + str(pdb_end_residue)
        else:
            residue_constructed = ""

        ## contig assembly
        contigs_constructed = pdb_chain + residue_constructed + "/0: " + binder_length_constructed
        if contigs_override == "":
            contigs = contigs_constructed
        else:
            contigs = contigs_override

        # determine where to save
        path = name
        while os.path.exists(f"{outputs_directory}/{path}_0.pdb"):
          path = name + "_" + ''.join(random.choices(string.ascii_lowercase + string.digits, k=5))

        flags = {"contigs":contigs,
                "pdb":pdb,
                "order":order,
                "iterations":iterations,
                "symmetry":symmetry,
                "hotspot":hotspot,
                "path":path,
                "chains":chains,
                "add_potential":add_potential,
                "num_designs":num_designs,
                "use_beta_model":use_beta_model,
                "visual":visual,
                "outputs_directory":outputs_directory}

        for k,v in flags.items():
          if isinstance(v,str):
            flags[k] = v.replace("'","").replace('"','')

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
        opts = [f"--pdb={outputs_directory}/{path}_0.pdb",
                f"--loc={outputs_directory}/{path}",
                f"--contig={contigs_str}",
                f"--copies={copies}",
                f"--num_seqs={num_seqs}",
                f"--num_recycles={num_recycles}",
                f"--rm_aa={rm_aa}",
                f"--mpnn_sampling_temp={mpnn_sampling_temp}",
                f"--num_designs={num_designs}"]
        if initial_guess: opts.append("--initial_guess")
        if use_multimer: opts.append("--use_multimer")
        if use_solubleMPNN: opts.append("--use_soluble")
        opts = ' '.join(opts)

        command_design = f"python colabdesign/rf/designability_test.py {opts}"
        os.system(command_design)

        print("running Prodigy")
        prodigy_run(f"{outputs_directory}/{path}/mpnn_results.csv", f"{outputs_directory}/{path}/all_pdb")

        command_mv = f"mkdir {outputs_directory}/{path}/traj && mv {outputs_directory}/traj/{path}* {outputs_directory}/{path}/traj && mv {outputs_directory}/{path}_* {outputs_directory}/{path}"
        command_zip = f"zip -r {path}.result.zip {outputs_directory}/{path}*"
        command_collect = f"mv {path}.result.zip /{outputs_directory} && mv {outputs_directory}/{path}/best.pdb /{outputs_directory}/{path}_best.pdb && mv {outputs_directory}/{path}/mpnn_results.csv /{outputs_directory}/{path}_scores.csv"
        os.system(command_mv)
        os.system(command_zip)
        os.system(command_collect)

        # enrich and summarise run and results information and write to csv file
        print("running enricher")
        enricher(outputs_directory, cfg)

        print("design complete...")
        end_time = time.time()
        duration = end_time - start_time
        print(f"executed in {duration:.2f} seconds.")

if __name__ == "__main__":
    my_app()
