
import os
import json
import requests

from tempfile import TemporaryDirectory, NamedTemporaryFile

from plex import CoreTools, ScatteringMethod, plex_init, plex_run, plex_vectorize, plex_mint, plex_upload

plex_dir = os.path.dirname(os.path.dirname(os.getcwd()))
plex_path = os.path.join(plex_dir, "plex")
jobs_dir = os.path.join(plex_dir, "jobs")


def move_from_file_url_to_ipfs(url: str) -> str:
    # Create a temporary directory using the `with` statement, so it's automatically cleaned up when we're done
    with TemporaryDirectory() as tmp_dir:
        # Get the file name from the url
        file_name = url.split('/')[-1]

        # Create the path to save the file
        save_path = os.path.join(tmp_dir, file_name)

        # Send a HTTP request to the url
        response = requests.get(url, stream=True)

        if response.status_code == 200:
            # If the request is successful, open the file in write mode and download the file
            with open(save_path, 'wb') as f:
                for chunk in response.iter_content(chunk_size=1024): 
                    if chunk: 
                        f.write(chunk)
            print(f"File downloaded successfully and saved as {save_path}")
            cid = plex_upload(save_path, plex_path=plex_path)
        else:
            print(f"Failed to download file. HTTP Response Code: {response.status_code}")
            cid = ""
        return cid, file_name

url = 'https://raw.githubusercontent.com/labdao/plex/447-add-rfdiffusion-to-plex/tools/rfdiffusion/6vja_stripped.pdb'
protein_6vja_cid, protein_6vja_file_name = move_from_file_url_to_ipfs(url)

hotspot_rfdiffusion_tool = {
    "class": "CommandLineTool",
    "name": "rfdiffusion_hotspot",
    "description": "design protein binders; generally useful for conditional generation of protein backbones",
    "baseCommand": ["/bin/bash", "-c"],
    "arguments": [
        "source activate SE3nv && python3 /app/scripts/run_inference.py 'contigmap.contigs=[$(inputs.motif.default) $(inputs.binder_length_min.default)-$(inputs.binder_length_max.default)]' $(inputs.hotspot.default) inference.input_pdb=$(inputs.protein.filepath) inference.output_prefix=/outputs/$(inputs.protein.basename)_backbone inference.num_designs=$(inputs.number_of_designs.default) denoiser.noise_scale_ca=0 denoiser.noise_scale_frame=0;"
      ],
    "dockerPull": "public.ecr.aws/p7l9w5o7/rfdiffusion:latest@sha256:0a6ff53004958ee5e770b0b25cd7f270eaf9fc285f6e91f17ad4024d2cc4ea91",
    "gpuBool": True,
    "networkBool": False,
    "inputs": {
      "protein": {
        "type": "File",
        "item": "",
        "glob": ["*.pdb"]
      },
      "motif": {
        "type": "string",
        "item": "",
        "default": "D46-200/0"
      },
      "hotspot": {
        "type": "string",
        "item": "",
        "default": "'ppi.hotspot_res=[D170, D171, D172, D173, D174, D76, D161]'"
      },
      "binder_length_min": {
        "type": "int",
        "item": "",
        "default": "50"
      },
      "binder_length_max": {
        "type": "int",
        "item": "",
        "default": "100"
      },
      "number_of_designs": {
        "type": "int",
        "item": "",
        "default": "10"
      }
    },
    "outputs": {
      "designed_backbones": {
        "type": "Array",
        "item": "File",
        "glob": ["*_backbone_*.pdb"]
      },
      "first_designed_backbone": {
        "type": "File",
        "item": "",
        "glob": ["*_backbone_0.pdb"]
      }
    }
}

with NamedTemporaryFile(suffix=".json", delete=False, mode='w') as temp_file:
    # Use json.dump to write the dictionary to the file
    json.dump(hotspot_rfdiffusion_tool, temp_file)
    print(f"Temporary file saved as {temp_file.name}")

    # flush data to disk
    temp_file.flush()
    hotspot_rfdiffusion_tool_cid = plex_upload(temp_file.name, wrap_file=False, plex_path=plex_path)
    print(f"Tool saved to IPFS as {hotspot_rfdiffusion_tool_cid}")

initial_io_cid = plex_init(
    hotspot_rfdiffusion_tool_cid,
    plex_path=plex_path,
    protein=[f'{protein_6vja_cid}/{protein_6vja_file_name}'],
    )

CACHE = True
if CACHE:
    completed_io_cid = 'QmbHneUZZpNh24is1uTnCpMYKRWStgEqAaYL4aDjN6tNzQ'
    io_cid_file_path = ''
else:
    completed_io_cid, io_cid_file_path = plex_run(initial_io_cid, plex_path=plex_path)

vectors = plex_vectorize(completed_io_cid, hotspot_rfdiffusion_tool_cid, plex_path=plex_path)

print(vectors['first_designed_backbone']['cidPaths'])
