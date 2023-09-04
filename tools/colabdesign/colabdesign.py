#@title Input protein sequence(s), then hit `Runtime` -> `Run all`
# from google.colab import files
import os
# plex_init is used to instantiate a DAG for the network to process. Every DAG is a simple JSON file.
from plex import CoreTools, plex_init
# plex_run is submits the constructed DAG to the network
from plex import plex_run


# wallet information
user_account = '0xcaf6A0c4468087d76e6B2917cea10F0E1aA2f9D4' #@param {type:"string"}
#@markdown  * Put your Metamask wallet address in the form above
tool_filepath = "_colabdesign-dev.ipynb" #@param {type:"string"}
input_protein_filepath = "6vja_stripped.pdb" #@param {type:"string"}
input_config_filepath = "config.yaml" #@param {type:"string"}

# write wallet address
os.environ["RECIPIENT_WALLET"] = user_account

# construct and execute DAG
initial_io_cid = plex_init(
    tool_path = tool_filepath,
    config = input_config_filepath,
    protein = input_protein_filepath, 
    auto_run=True
)



# run plex
completed_io_cid, completed_io_filepath = plex_run(initial_io_cid)

# print results
print("DAG filepath", completed_io_filepath)
print("DAG location", completed_io_cid)

# mint token
from plex import plex_mint
# using the autotask webhook enables gasless minting
os.environ["AUTOTASK_WEBHOOK"] = "https://api.defender.openzeppelin.com/autotasks/e15b3f39-28f8-4d30-9bf3-5d569bdf2e78/runs/webhook/8315d17c-c493-4d04-a257-79209f95bb64/2gmqi9SRRAQMoy1SRdktai"
plex_mint(completed_io_cid)

