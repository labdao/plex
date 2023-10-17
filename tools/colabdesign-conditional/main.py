# #@title **RFdiffusion Conditional Fold, ProteinMPNN, and AF2 -** Inputs and Parameters
# #@markdown **RFDiffusion Parameters**
# name = "test"
# blueprint_mode = "manual" #@param ["manual", "automated"]

# run_mode = "unconditional"
# pdb = "input.pdb"

# #@markdown ---
# #@markdown **Manual conditional fold blueprint** (define number of secondary structure `elements` (SSE))
# elements = 1 #@param ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20"] {type:"raw"}
# buff_length = 0 # @param {type:"number"}
# #@markdown *  define the buffer length between SSEs
# def_elen = 5 # @param {type:"number"}
# #@markdown *  default element lengths
# def_ss = 'sheet' # @param ["sheet", "helix", "coil"]
# #@markdown *  default secondary structure
# sss = {'helix':'H','sheet':'E','coil':'C'}
# def_ss = sss[def_ss]
# def_cont = 'contact' #@param ["no_contact", "contact"]
# #@markdown *  default element contact
# conts = {'no_contact':'0','contact':'1'}
# def_cont = conts[def_cont]

# #@markdown ___
# #@markdown **Automated conditional fold blueprint** (from input PDB)
# pdb = "" #@param {type:"string"}
# chain = "B" #@param {type:"string"}
# trim_loops = False #@param {type:"boolean"}
# if chain == "": chain = None
# #@markdown ---
# run_name = "def" #@param {type:"string"}
# name = run_name
# #@markdown ---

# num_designs = 3 # @param {type:"number"}
# iterations = 25 #@param ["25", "50", "100", "200"]
# mask_loops = False #@param {type:"boolean"}
# mask_contacts = False #@param {type:"boolean"}
# #@markdown **Optional**: specify target info for **binder design**
# use_target = True #@param {type:"boolean"}
# target_pdb = "input.pdb" #@param {type:"string"}
# #@markdown * leave blank to receive a file upload prompt
# #@markdown * enter filepath to a pdb file that is already loaded
# #@markdown * enter the PDB code
# #@markdown * enter the Alphafold-DB code
# target_chain = "D" #@param {type:"string"}
# target_hotspot = "" #@param {type:"string"}
# denoiser_noise_scale_ca = 0 # @param {type:"number"}
# #@markdown * From 0 to 1, default is 0
# denoiser_noise_scale_frame = 0 # @param {type:"number"}
# #@markdown * From 0 to 1, default is 0

# #@markdown ---
# #@markdown **ProteinMPNN + AF2 Parameters**
# #@title run **ProteinMPNN** to generate a sequence and **AlphaFold** to validate
# num_seqs = 8 #@param ["1", "2", "4", "8", "16", "32", "64"] {type:"raw"}
# #@markdown * number of MPNN sequences to produce
# initial_guess = True #@param {type:"boolean"}
# #@markdown * helps for PPI
# num_recycles = 3 #@param ["0", "1", "2", "3", "6", "12"] {type:"raw"}
# use_multimer = True #@param {type:"boolean"}
# rm_aa = "C" #@param {type:"string"}
# mpnn_sampling_temp = 0.1 #@param ["0.0001", "0.1", "0.15", "0.2", "0.25", "0.3", "0.5", "1.0"] {type:"raw"}
# use_solubleMPNN = False #@param {type:"boolean"}
# #@markdown - for **binder** design, we recommend `initial_guess=True, num_recycles=3`


import os, time, sys
# # Create the /inputs directory
# if not os.path.exists('/inputs'):
#   os.mkdir('/inputs')
# # Create the /outputs directory
# if not os.path.exists('/outputs'):
#   os.mkdir('/outputs')
# # Create a symlink to '/inputs' in the current working directory
# if not os.path.exists('inputs'):
#     os.symlink('/inputs', 'inputs')
# # Create a symlink to '/outputs' in the current working directory
# if not os.path.exists('outputs'):
#     os.symlink('/outputs', 'outputs')

######################################################################



# from IPython.display import display
# import ipywidgets as widgets
import torch
import random, string, re
import numpy as np
import subprocess
import matplotlib.pyplot as plt
import py3Dmol
# from google.colab import files, output

from string import ascii_uppercase, ascii_lowercase
alphabet_list = list(ascii_uppercase+ascii_lowercase)

def get_pdb(pdb_code=None):
  if pdb_code is None or pdb_code == "":
    upload_dict = files.upload()
    pdb_string = upload_dict[list(upload_dict.keys())[0]]
    with open("tmp.pdb","wb") as out: out.write(pdb_string)
    return "tmp.pdb"
  elif os.path.isfile(pdb_code):
    return pdb_code
  elif len(pdb_code) == 4:
    os.system(f"wget -qnc https://files.rcsb.org/view/{pdb_code}.pdb")
    return f"{pdb_code}.pdb"
  else:
    os.system(f"wget -qnc https://alphafold.ebi.ac.uk/files/AF-{pdb_code}-F1-model_v3.pdb")
    return f"AF-{pdb_code}-F1-model_v3.pdb"



def pdb_to_string(pdb_file, chains=None, models=[1]):
  '''read pdb file and return as string'''

  MODRES = {'MSE':'MET','MLY':'LYS','FME':'MET','HYP':'PRO',
            'TPO':'THR','CSO':'CYS','SEP':'SER','M3L':'LYS',
            'HSK':'HIS','SAC':'SER','PCA':'GLU','DAL':'ALA',
            'CME':'CYS','CSD':'CYS','OCS':'CYS','DPR':'PRO',
            'B3K':'LYS','ALY':'LYS','YCM':'CYS','MLZ':'LYS',
            '4BF':'TYR','KCX':'LYS','B3E':'GLU','B3D':'ASP',
            'HZP':'PRO','CSX':'CYS','BAL':'ALA','HIC':'HIS',
            'DBZ':'ALA','DCY':'CYS','DVA':'VAL','NLE':'LEU',
            'SMC':'CYS','AGM':'ARG','B3A':'ALA','DAS':'ASP',
            'DLY':'LYS','DSN':'SER','DTH':'THR','GL3':'GLY',
            'HY3':'PRO','LLP':'LYS','MGN':'GLN','MHS':'HIS',
            'TRQ':'TRP','B3Y':'TYR','PHI':'PHE','PTR':'TYR',
            'TYS':'TYR','IAS':'ASP','GPL':'LYS','KYN':'TRP',
            'CSD':'CYS','SEC':'CYS'}
  restype_1to3 = {'A': 'ALA','R': 'ARG','N': 'ASN',
                  'D': 'ASP','C': 'CYS','Q': 'GLN',
                  'E': 'GLU','G': 'GLY','H': 'HIS',
                  'I': 'ILE','L': 'LEU','K': 'LYS',
                  'M': 'MET','F': 'PHE','P': 'PRO',
                  'S': 'SER','T': 'THR','W': 'TRP',
                  'Y': 'TYR','V': 'VAL'}

  restype_3to1 = {v: k for k, v in restype_1to3.items()}

  if chains is not None:
    if "," in chains: chains = chains.split(",")
    if not isinstance(chains,list): chains = [chains]
  if models is not None:
    if not isinstance(models,list): models = [models]

  modres = {**MODRES}
  lines = []
  seen = []
  model = 1
  for line in open(pdb_file,"rb"):
    line = line.decode("utf-8","ignore").rstrip()
    if line[:5] == "MODEL":
      model = int(line[5:])
    if models is None or model in models:
      if line[:6] == "MODRES":
        k = line[12:15]
        v = line[24:27]
        if k not in modres and v in restype_3to1:
          modres[k] = v
      if line[:6] == "HETATM":
        k = line[17:20]
        if k in modres:
          line = "ATOM  "+line[6:17]+modres[k]+line[20:]
      if line[:4] == "ATOM":
        chain = line[21:22]
        if chains is None or chain in chains:
          atom = line[12:12+4].strip()
          resi = line[17:17+3]
          resn = line[22:22+5].strip()
          if resn[-1].isalpha(): # alternative atom
            resn = resn[:-1]
            line = line[:26]+" "+line[27:]
          key = f"{model}_{chain}_{resn}_{resi}_{atom}"
          if key not in seen: # skip alternative placements
            lines.append(line)
            seen.append(key)
      if line[:5] == "MODEL" or line[:3] == "TER" or line[:6] == "ENDMDL":
        lines.append(line)
  return "\n".join(lines)

def from_pdb(pdb_code=None, chains=None, trim_loops=False,
             mask_contacts=False, return_pdb_str=False):

  import pydssp
  def process(secondary_structure, contact_map):
    secondary_structure = np.array(secondary_structure)
    # Find the start and end indices of the continuous secondary structure elements
    sse_start,sse_end = [],[]
    for i, current_element in enumerate(secondary_structure):
      if current_element in ["H", "E", "C"]:
        if i == 0 or secondary_structure[i-1] != current_element:
          sse_start.append(i)
        if i == len(secondary_structure) - 1 or secondary_structure[i+1] != current_element:
          sse_end.append(i)

    sse_types = secondary_structure[sse_start]
    sse_lengths = np.array(sse_end) - np.array(sse_start) + 1
    num_sse = len(sse_lengths)
    reduced_contact_map = np.full((num_sse, num_sse), '0', dtype=object)
    np.fill_diagonal(reduced_contact_map, sse_types)

    for i in range(num_sse):
      for j in range(num_sse):
        if i != j and sse_types[i] != "C" and sse_types[j] != "C":
          interaction_mask = np.any(contact_map[sse_start[i]:sse_end[i]+1, sse_start[j]:sse_end[j]+1])
          reduced_contact_map[i, j] = str(interaction_mask.astype(int))
          if mask_contacts and reduced_contact_map[i, j] == "1":
            reduced_contact_map[i, j] = "?"


    return {"txt":sse_lengths, "adj":reduced_contact_map}

  def coord_2_cb(coord):
    N,Ca,C = coord[:,0],coord[:,1],coord[:,2]
    # recreate Cb given N,Ca,C
    b = Ca - N
    c = C - Ca
    a = np.cross(b, c)
    Cb = -0.57910144*a + 0.5689693*b - 0.5441217*c + Ca
    return Cb
  pdb_filename = get_pdb(pdb_code)
  pdb_str = pdb_to_string(pdb_filename, chains=chains)
  coord = pydssp.read_pdbtext(pdb_str)

  ss = pydssp.assign(coord)

  # filter single length sse
  for i in range(len(ss)):
    if ss[i] in ["H","E"]:
      if (i == (len(ss)-1) or ss[i] != ss[i+1]) and (i == 0 or ss[i] != ss[i-1]):
        ss[i] = "-"

  if not trim_loops:
    ss = [("C" if s == "-" else s) for s in ss]
  cb = coord_2_cb(coord)
  con = np.sqrt(np.square(cb[:,None] - cb[None,:]).sum(-1)) < 6.0
  out = process(ss, con)
  if return_pdb_str:
    out["pdb_str"] = pdb_str
  return out

def get_adj_ss(adj, txt, buff=0, mask_contacts=False):
  # select non-zero elements
  idx = []
  for i in range(len(adj)):
    if txt[i] > 0:
      idx.append(i)

  L = (len(idx) + 1) * buff + sum(txt)
  full_adj = np.full((L,L),2)
  full_sse = np.full((L,),3)
  n = buff
  for i in idx:
    ss = {"H":0, "E":1, "C":2, "?":3}[adj[i][i]]
    full_sse[n:n+txt[i]] = ss
    m = buff
    for j in idx:
      k = str(adj[i][j])
      if i == j:
        val = {"H":0,"E":0,"C":0,"?":2}[k]
      else:
        if mask_contacts and k == "1": k = "?"
        val = {"0":0,"1":1,"?":2}[k]
      full_adj[n:n+txt[i],m:m+txt[j]] = val
      m += txt[j] + buff
    n += txt[i] + buff
  return {"adj":full_adj,"sse":full_sse}

class blueprint_gui:

  def _toggle_callback(self, row, col):
    if row == col:
      new_value = {"H":"E","E":"C","C":"?","?":"H"}[self.adj[row][col]]
      self.txt[row] = {"H": def_elen, "E": def_elen, "C": def_elen, "?": 0}[new_value]
      #self.txt[row] = {"H": 19, "E": 5, "C": 3, "?": 0}[new_value]
      self.adj[row][col] = new_value
      for i in range(self.elements):
        if i != row:
          if new_value == "?":
            self.adj[row][i] = self.adj[i][col] = "?"
          elif self.adj[i][i] != "?" and new_value in ["C","H"]:
            self.adj[row][i] = self.adj[i][col] = '0'
    else:
      if self.adj[row][row] not in ["C","?"] and self.adj[col][col] not in ["C","?"]:
        new_value = {"0":"1","1":"?","?":"0"}[self.adj[row][col]]
        self.adj[row][col] = self.adj[col][row] = new_value

  def _text_callback(self, row, new_value):
    self.txt[row] = int(new_value)

  def _update_callback(self, position, add):
    if position < 0: position = self.elements
    self.elements = self.elements + 1 if add else self.elements - 1
    if self.elements < 0: self.elements = 0
    adj = [['' for _ in range(self.elements)] for _ in range(self.elements)]
    txt = ['' for _ in range(self.elements)]
    for row in range(self.elements):
      old_row = row if row < position else row - 1 if add else row + 1
      if add and row == position:
        txt[row] = def_elen
        #####################################
      else:
        txt[row] = self.txt[old_row]
      for col in range(self.elements):
        old_col = col if col < position else col - 1 if add else col + 1
        if add and (row == position or col == position):
          if row == col:
            adj[row][col] = def_ss
            #adj[row][col] = 'H'
            ###################################
          else:
            cell = self.adj[old_row][old_row] if col == position else self.adj[old_col][old_col]
            adj[row][col] = cell if cell == "?" else def_cont
            #adj[row][col] = cell if cell == "?" else '0'
            ########################################
        else:
          adj[row][col] = self.adj[old_row][old_col]
    self.adj = adj
    self.txt = txt
#############################################################################
  def _create_html(self):
    # HTML for initial grid
    html_grid = f'<div class="pos"></div>'
    for row in range(self.elements): html_grid += f'<div class="pos">{row}</div>'
    html_grid += f'<div class="pos"></div>'
    for row in range(self.elements):
      html_grid += f'<div class="pos">{row}</div>'
      for col in range(self.elements):
        value = self.adj[row][col]
        # print(value)
        bgcolor = {"H":"red","E":"yellow","C":"lime","?":"lightgray","0":"white","1":"lightblue"}[value]
        if row != col and (self.adj[row][row] in ["?","C"] or self.adj[col][col] in ["?","C"]):
          opacity = 0.1
        else:
          opacity = 1.0
        html_grid += f'<div class="grid-item" id="cell_{row}_{col}" style="background-color:{bgcolor};opacity:{opacity}">{value}</div>'
      html_grid += f'<div><input class="text" type="number" id="cell_{row}" min="0" value="{self.txt[row]}" onchange="textFieldChanged({row}, this)"></div>'

    self.html_code = f"""
    <style>
    {self._CSS}
    .grid-container {{
      display: grid;
      grid-template-columns: repeat({self.elements+2}, 30px);
      gap: 2px;
    }}
    </style>
    <script>{self._JS}</script>
    <label>resize:</label>
    <button id="add" style="width:25px" class="button" onclick="updateGrid(true)">+</button>
    <button id="remove" style="width:25px" class="button" onclick="updateGrid(false)">-</button>
    <input type="number" id="position" min="-1" max="{self.elements}" value="0" class="text">
    <label>(indicate where to +/- an element)</label>
    <div class="grid-container">{html_grid}</div>
    """

# class RFdiff_gui(blueprint_gui):
class RFdiff_gui(blueprint_gui):

  # def __init__(self, elements=elements, adj=None, txt=None, buff_length=buff_length, name=name):
  def __init__(self, blueprint_mode, elements, name, use_target, target_chain, denoiser_noise_scale_ca, denoiser_noise_scale_frame, target_hotspot, iterations, mask_loops, mask_contacts, use_solubleMPNN, 
         use_multimer, initial_guess, num_designs, mpnn_sampling_temp, rm_aa, num_recycles, num_seqs, target_pdb, buff_length, def_ss, def_cont, def_elen, outputs, adj=None, txt=None):
    self.path = self.name = name
    # self.input = widgets.Output()
    # self.output = widgets.Output()
    self.buff_length = buff_length

    # small_button_style = widgets.Layout(width='30px', height='30px', border='2px solid black')
    # button_style = widgets.Layout(width='84px', height='35px', border='2px solid black')
    # self.buttons = {
    #     "buff_length": widgets.BoundedIntText(description='buff_length', value=self.buff_length, min=0, max=20),
    #     "reset":       widgets.Button(description='reset',     layout=button_style),
    #     "animate":     widgets.Button(description='animate',   layout=button_style),
    #     "freeze":      widgets.Button(description='freeze',    layout=button_style),
    #     "download":    widgets.Button(description='download',  layout=button_style),
    #     "color":       widgets.Dropdown(
    #                     options=['SSE','pLDDT'],
    #                     value='SSE',
    #                     description='color',
    #                     disabled=False)
    # }
    # self.buttons["animate"].on_click(self._plot_pdb)
    # self.buttons["freeze"].on_click(self._plot_pdb)
    # self.buttons["download"].on_click(self._download)
    # self.buttons["color"].observe(self._plot_pdb)
    # self._plot = {"mode":"freeze","color":"SSE"}

    # prep inputs
    if adj is not None and txt is not None:
      self.elements = len(adj)
      self.adj, self.txt = adj,txt
    else:
      self.elements = elements
      self.adj = [[def_ss if row == col else def_cont for col in range(self.elements)] for row in range(self.elements)]
      #self.adj = [["H" if row == col else "0" for col in range(self.elements)] for row in range(self.elements)]
      self.txt = [def_elen for _ in range(self.elements)]
      #self.txt = [19 for _ in range(self.elements)]

    # output.register_callback("update_callback", self._update_callback)
    # output.register_callback("toggle_callback", self._toggle_callback)
    # output.register_callback("text_callback",   self._text_callback)
    self._CSS = open("blueprint.css","r").read()
    self._JS = open("blueprint.js","r").read()

  # def _redraw(self):
  #   with self.input:
  #     self._create_html()
  #     self.input.clear_output(wait=True)
  #     display(
  #         widgets.VBox([
  #         widgets.HTML(self.html_code),
  #         widgets.Label("Options"),
  #         self.buttons["buff_length"],
  #       ])
  #     )

  # def display_input(self):
  #   self._redraw()
  #   display(self.input)

  # def display_output(self):
  #   display(self.output)

  def _download(self, button):
    os.system(f"zip -r {self.path}.result.zip inputs/{self.path}* inputs/traj/{self.path}*")
    files.download(f"{self.path}.result.zip")

  # def _plot_pdb(self, change):
  #   update = False
  #   if isinstance(change, widgets.Button):
  #     self._plot["mode"] = change.description
  #     update = True
  #   elif isinstance(change, dict) and change['name'] == 'value':
  #     widget = change['owner']
  #     if isinstance(widget, widgets.Dropdown):
  #       self._plot["color"] = change["new"]
  #       update = True
  #   if update:
  #     view = py3Dmol.view()
  #     if self._plot["mode"] == "animate":
  #       pdb = f"inputs/traj/{self.path}_0_pX0_traj.pdb"
  #       pdb_str = open(pdb,'r').read()
  #       view.addModelsAsFrames(pdb_str,'pdb')
  #     else:
  #       pdb = f"outputs/{self.path}_0.pdb"
  #       pdb_str = open(pdb,'r').read()
  #       view.addModel(pdb_str,'pdb')
  #     if self._plot["color"] == "SSE":
  #       view.setStyle({"ss":"h","chain":"A"},{'cartoon': {'color':'red'}})
  #       view.setStyle({"ss":"c","chain":"A"},{'cartoon': {'color':'lime'}})
  #       view.setStyle({"ss":"s","chain":"A"},{'cartoon': {'color':'yellow'}})
  #       if self.use_target:
  #         view.setStyle({"chain":"B"},{'cartoon': {'color':'white'}})
  #     else:
  #       view.setStyle({'cartoon': {'colorscheme': {'prop':'b','gradient': 'roygb','min':0.5,'max':0.9}}})
  #     view.zoomTo()
  #     if self._plot["mode"] == "animate":
  #       view.animate({'loop': 'backAndForth'})
  #     out = widgets.Output()
  #     with out: view.show()
  #     toggle = self.buttons["freeze"] if self._plot["mode"] == "animate" else self.buttons["animate"]
  #     with self.output:
  #       self.output.clear_output(wait=True)
  #       display(widgets.VBox([out, widgets.HBox([toggle, self.buttons["download"], self.buttons["color"]])]))

  def _make_path(self):
    os.makedirs(f"inputs/{self.path}", exist_ok=True)
    while os.path.exists(f"inputs/{self.path}_0.pdb"):
      self.path = self.name + "_" + ''.join(random.choices(string.ascii_lowercase + string.digits, k=5))
      os.makedirs(f"inputs/{self.path}", exist_ok=True)

  def _get_adj_ss(self, mask_contacts=False):
    # get unique path
    full = get_adj_ss(adj=self.adj,
                      txt=self.txt,
                      # buff=self.buttons["buff_length"].value,
                      buff=self.buff_length,
                      mask_contacts=mask_contacts)
    self._sse = full["sse"]
    self._adj = full["adj"]

    # save results
    loc = [f"inputs/{self.path}/tmp_ss.pt",
           f"inputs/{self.path}/tmp_adj.pt"]
    torch.save(torch.from_numpy(self._sse).float(),loc[0])
    torch.save(torch.from_numpy(self._adj).float(),loc[1])

  # def diffuse(self, use_target, iterations=50,
  #            mask_loops=True,
  #            mask_contacts=False,
  #            extra_cmd=None):
  def diffuse(self, use_target, num_designs, outputs, iterations=50,
              mask_loops=True,
              mask_contacts=False,
              extra_cmd=None):
    self.use_target = use_target
    # self._redraw()
    self._make_path()
    self._get_adj_ss(mask_contacts=mask_contacts)
    # # run
    # with self.output: ## commented this out
    # self.output.clear_output()
    cmd = ["./RFdiffusion/run_inference.py",
          f"inference.num_designs={num_designs}",
          f"inference.output_prefix={outputs}/{self.path}",
          "scaffoldguided.scaffoldguided=True",
          f"scaffoldguided.scaffold_dir=inputs/{self.path}",
          f"diffuser.T={iterations}",
          f"scaffoldguided.mask_loops={mask_loops}",
          "inference.dump_pdb=True",
          "inference.dump_pdb_path=/dev/shm"]

    if extra_cmd is not None:
      cmd += extra_cmd

    self.cmd_str = " ".join(cmd)
    self._run(self.cmd_str, iterations, num_designs)
  # self._plot_pdb(self.buttons["freeze"])    
  # with self.output: ## commented this out
  #     self.output.clear_output()
  #     cmd = ["./RFdiffusion/run_inference.py",
  #            f"inference.num_designs={num_designs}",
  #            f"inference.output_prefix=outputs/{self.path}",
  #            "scaffoldguided.scaffoldguided=True",
  #            f"scaffoldguided.scaffold_dir=inputs/{self.path}",
  #            f"diffuser.T={iterations}",
  #            f"scaffoldguided.mask_loops={mask_loops}",
  #            "inference.dump_pdb=True",
  #            "inference.dump_pdb_path=/dev/shm"]

  #     if extra_cmd is not None:
  #       cmd += extra_cmd

  #     self.cmd_str = " ".join(cmd)
  #     self._run(self.cmd_str, iterations)
  #   self._plot_pdb(self.buttons["freeze"])

  # def _run(self, command, steps, num_designs=num_designs):
  def _run(self, command, steps, num_designs):
    def run_command_and_get_pid(command):
      pid_file = '/dev/shm/pid'
      os.system(f'nohup {command} > /dev/null & echo $! > {pid_file}')
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

    # run_output = widgets.Output()
    # progress = widgets.FloatProgress(min=0, max=1, description='running', bar_style='info')
    # display(widgets.VBox([progress, run_output]))

    # clear previous run
    for n in range(steps):
      if os.path.isfile(f"/dev/shm/{n}.pdb"):
        os.remove(f"/dev/shm/{n}.pdb")

    pid = run_command_and_get_pid(command)
    try:
      fail = False
      for _ in range(num_designs):
        # for each step
        for n in range(steps):
          wait = True
          while wait and not fail:
            time.sleep(0.5)
            # check if output generated
            if os.path.isfile(f"/dev/shm/{n}.pdb"):
              pdb_str = open(f"/dev/shm/{n}.pdb").read()
              if pdb_str[-3:] == "TER":
                wait = False
              elif not is_process_running(pid):
                fail = True
            elif not is_process_running(pid):
              fail = True

          # if fail:
          #   progress.bar_style = 'danger'
          #   progress.description = "failed"
          #   break
          # else:
          #   progress.value = (n+1) / steps
          #   with run_output:
          #     run_output.clear_output(wait=True)
          #     view = py3Dmol.view(js='https://3dmol.org/build/3Dmol.js')
          #     view.addModel(pdb_str,'pdb')
          #     view.setStyle({'cartoon': {'colorscheme': {'prop':'b','gradient': 'roygb','min':0.5,'max':0.9}}})
          #     view.zoomTo()
          #     view.show()
          # if os.path.exists(f"/dev/shm/{n}.pdb"):
          #   os.remove(f"/dev/shm/{n}.pdb")

        # if fail:
        #   progress.bar_style = 'danger'
        #   progress.description = "failed"
        #   break

      while is_process_running(pid):
        time.sleep(0.5)

    except KeyboardInterrupt:
      os.kill(pid, signal.SIGTERM)
      progress.bar_style = 'danger'
      progress.description = "stopped"

def main(blueprint_mode, elements, name, use_target, target_chain, denoiser_noise_scale_ca, denoiser_noise_scale_frame, target_hotspot, iterations, mask_loops, trim_loops, mask_contacts, use_solubleMPNN, 
         use_multimer, initial_guess, num_designs, mpnn_sampling_temp, rm_aa, num_recycles, num_seqs, target_pdb, buff_length, def_ss, def_cont, def_elen, pdb, outputs, chain):

    if blueprint_mode == "automated":
        pdb_feats = from_pdb(pdb, chains=chain, trim_loops=trim_loops)
        buff_length=(5 if trim_loops else 0)
        rfdiff = RFdiff_gui(**pdb_feats, blueprint_mode=blueprint_mode, elements=elements, name=name, use_target=use_target, target_chain=target_chain, denoiser_noise_scale_ca=denoiser_noise_scale_ca, denoiser_noise_scale_frame=denoiser_noise_scale_frame, target_hotspot=target_hotspot, iterations=iterations, mask_loops=mask_loops, mask_contacts=mask_contacts, use_solubleMPNN=use_solubleMPNN, 
         use_multimer=use_multimer, initial_guess=initial_guess, num_designs=num_designs, mpnn_sampling_temp=mpnn_sampling_temp, rm_aa=rm_aa, num_recycles=num_recycles, num_seqs=num_seqs, target_pdb=target_pdb, buff_length=buff_length, def_ss=def_ss, def_cont=def_cont, def_elen=def_elen, outputs=outputs)
    else:
        rfdiff = RFdiff_gui(blueprint_mode, elements, name, use_target, target_chain, denoiser_noise_scale_ca, denoiser_noise_scale_frame, target_hotspot, iterations, mask_loops, mask_contacts, use_solubleMPNN, 
         use_multimer, initial_guess, num_designs, mpnn_sampling_temp, rm_aa, num_recycles, num_seqs, target_pdb, buff_length, def_ss, def_cont, def_elen, outputs)

    # RFD
    if use_target:
        # prep target features
        rfdiff._make_path()
        # print('rfdiff.path', rfdiff.path)
        # path = f"inputs/{rfdiff.path}/target"
        path = directory_path = os.path.dirname(target_pdb)
        # print('path', path)
        # print('target_pdb', target_pdb)
        os.makedirs(path, exist_ok=True)
        target = from_pdb(target_pdb, target_chain, return_pdb_str=True)
        # target_pdb = f"{path}/input.pdb"
        with open(target_pdb,"w") as handle:
            handle.write(target["pdb_str"])
            full = get_adj_ss(adj=target["adj"], txt=target["txt"])
            torch.save(torch.from_numpy(full["sse"]).float(),f"{path}/ss.pt")
            torch.save(torch.from_numpy(full["adj"]).float(),f"{path}/adj.pt")
        extra_cmd = ["scaffoldguided.target_pdb=True",
                    # f"scaffoldguided.target_path={path}/input.pdb",
                    f"scaffoldguided.target_path={target_pdb}",
                    f"scaffoldguided.target_ss={path}/ss.pt",
                    f"scaffoldguided.target_adj={path}/adj.pt",
                    f"denoiser.noise_scale_ca={denoiser_noise_scale_ca}",
                    f"denoiser.noise_scale_frame={denoiser_noise_scale_frame}"]
        if target_hotspot != "":
            extra_cmd += [f"'ppi.hotspot_res=[{target_hotspot}]'"]
    else:
        extra_cmd = None

    if "rfdiff" in dir():
        # rfdiff.display_output()
        rfdiff.diffuse(use_target, num_designs, outputs, iterations,
                    mask_loops=mask_loops,
                    mask_contacts=mask_contacts,
                    extra_cmd=extra_cmd)
    else:
        print("Error, looks like you didn't run the cell above")

    parser = PDB.PDBParser(QUIET=True)
    pdb1 = f'{outputs}/{name}_0.pdb'
    structure = parser.get_structure("temp", pdb1)
    model = structure[0]

    for model in structure:
        chain_len = dict()
        for chain in model:
            chain_id = str(chain).split('=')[1].split('>')[0]
            res_no = 0
            non_resi = 0
            for r in chain.get_residues():
                if r.id[0] == ' ':
                    res_no +=1
                else:
                    non_resi +=1
            chain_len[chain_id] = res_no
    if len(chain_len) == 1:
        contigs = [f"{chain_len['A']}-{chain_len['A']}"]
    else:
        contigs = [f"{chain_len['A']}-{chain_len['A']}",f"B1-{chain_len['B']}"]

    # PMPNN + AF2
    contigs_str = ":".join(contigs)
    opts = [f"--pdb={outputs}/{name}_0.pdb",
            f"--loc={outputs}",
            f"--contig={contigs_str}",
            f"--copies=1",
            f"--num_seqs={num_seqs}",
            f"--num_recycles={num_recycles}",
            f"--rm_aa={rm_aa}",
            f"--mpnn_sampling_temp={mpnn_sampling_temp}",
            f"--num_designs={num_designs}"]
    if initial_guess: opts.append("--initial_guess")
    if use_multimer: opts.append("--use_multimer")
    if use_solubleMPNN: opts.append("--use_soluble")
    opts = ' '.join(opts)
    print('running PMPNN + AF2')
    subprocess.run(f"python colabdesign/rf/designability_test.py {opts}", shell=True)

    # mpnn results to df
    df = pd.read_csv(f'{outputs}/mpnn_results.csv')
    df['path'] = f'{outputs}/all_pdb/design'+df['design'].astype(str)+'_n'+ df['n'].astype(str)+'.pdb'
    df = df.iloc[:,1:]

    print('running Prodigy')
    for i,r in df.iterrows():
        try:
            subprocess.run(["prodigy", "-q", r['path']], stdout=open('temp.txt', 'w'), check=True)
            with open('temp.txt', 'r') as f:
                lines = f.readlines()
                if lines:  # Check if lines is not empty
                    affinity = float(lines[0].split(' ')[-1].split('/')[0])
                    df.loc[i,'affinity'] = affinity
                else:
                    print(f"No output from prodigy for {r['path']}")
                    # Handle the case where prodigy did not produce output
        except subprocess.CalledProcessError:
            print(f"Prodigy command failed for {r['path']}")
    # export results
    df.to_csv(f'{outputs}/mpnn_results.csv',index=None)

def get_pdb_file_path(directory):
    for filename in os.listdir(directory):
        if filename.endswith(".pdb"):
            return os.path.join(directory, filename)
    return None

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


import hydra
from omegaconf import DictConfig, OmegaConf
import os
import torch
import pandas as pd
from Bio import PDB
import subprocess

@hydra.main(version_base=None, config_path="conf", config_name="config")
def my_app(cfg : DictConfig) -> None:

    # defining output directory
    if cfg.outputs.directory is None:
        outputs_directory = hydra.core.hydra_config.HydraConfig.get().runtime.output_dir
    else:
        outputs_directory = cfg.outputs.directory
    print(f"Output directory  : {outputs_directory}")

    # # defining receptor file paths
    # input_receptor_path = get_files_from_directory(cfg.inputs.receptors_directory, '.pdb')
    # input_ligand_path = get_files_from_directory(cfg.inputs.ligands_directory, '.pdb')
    target_pdb = get_files_from_directory(cfg.inputs.target_protein_directory, '.pdb')
    binder_template_pdb = get_files_from_directory(cfg.inputs.binder_protein_template_directory, '.pdb')

    # filtering receptor and ligand files if pattern is available
    if cfg.inputs.target_protein_pattern is not None:
        target_pdb = [file for file in target_pdb if cfg.inputs.target_protein_pattern in file]
        print("Retained target: ", target_pdb)
    if cfg.inputs.binder_protein_template_pattern is not None:
        binder_template_pdb = [file for file in binder_template_pdb if cfg.inputs.binder_protein_template_pattern in file]
        print("Retained binder template : ", binder_template_pdb)

    # Check if more than one file is retained
    if len(target_pdb) > 1 or len(binder_template_pdb) > 1:
        print("Error: More than one receptor or ligand file retained. Please check the input patterns.")
        sys.exit()

    target_pdb = target_pdb[0]
    binder_template_pdb = binder_template_pdb[0]

    # target_pdb = get_pdb_file_path(cfg.inputs.target_protein_directory)
    # binder_template_pdb = get_pdb_file_path(cfg.inputs.binder_protein_template_directory)

    main(cfg.params.blueprint_mode, cfg.params.elements, cfg.params.name, cfg.params.binder_design_target.use_target, cfg.params.binder_design_target.target_chain, 
        cfg.params.binder_design_target.denoiser_noise_scale_ca, cfg.params.binder_design_target.denoiser_noise_scale_frame, cfg.params.binder_design_target.target_hotspot,
        cfg.params.automated_conditional_fold_blueprint.iterations, cfg.params.automated_conditional_fold_blueprint.mask_loops, cfg.params.automated_conditional_fold_blueprint.trim_loops, cfg.params.automated_conditional_fold_blueprint.mask_contacts,
        cfg.params.ProteinMPNN.use_solubleMPNN, cfg.params.Alphafold.use_multimer, cfg.params.ProteinMPNN.initial_guess, cfg.params.automated_conditional_fold_blueprint.num_designs,
        cfg.params.ProteinMPNN.mpnn_sampling_temp, cfg.params.ProteinMPNN.rm_aa, cfg.params.Alphafold.num_recycles, cfg.params.ProteinMPNN.num_seqs, target_pdb, cfg.params.buff_length, cfg.params.def_ss, cfg.params.def_cont, cfg.params.def_elen, binder_template_pdb, cfg.outputs.directory, cfg.params.automated_conditional_fold_blueprint.chain)

if __name__ == "__main__":
    my_app()

    # print('running Prodigy')
    # for i,r in df.iterrows():
    #     subprocess.run(["prodigy", "-q", r['path']], stdout=open('temp.txt', 'w'), check=True)
    #     with open('temp.txt', 'r') as f:
    #         lines = f.readlines()
    #         if lines:  # Check if lines is not empty
    #             affinity = float(lines[0].split(' ')[-1].split('/')[0])
    #             df.loc[i,'affinity'] = affinity
    #         else:
    #             print(f"No output from prodigy for {r['path']}")
    #             # Handle the case where prodigy did not produce output

    # # prodigy # my original edit
    # print('running Prodigy')
    # for i,r in df.iterrows():
    #   subprocess.run(f"prodigy -q {r['path']} > temp.txt", shell=True)
    #   affinity = float(open('temp.txt','r').readlines()[0].split(' ')[-1].split('/')[0])
    #   df.loc[i,'affinity'] = affinity

  #   # prodigy # alternative edit
  #   for i,r in df.iterrows():
  # #   !prodigy -q {r['path']} > temp.txt
  #     with open('temp.txt', 'w') as f:
  #         subprocess.run(["prodigy", "-q", r['path']], stdout=f, check=True)
  #     affinity = float(open('temp.txt','r').readlines()[0].split(' ')[-1].split('/')[0])
  #     df.loc[i,'affinity'] = affinity
