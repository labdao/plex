import json
from typing import Dict, List
from pydantic import BaseModel, FilePath, Field
from pydantic import validator
from validators import *

# Load the JSON data
json_data = """
[
    {
      "outputs": {
        "best_docked_small_molecule": {
          "class": "File",
          "filepath": ""
        },
        "protein": {
          "class": "File",
          "filepath": ""
        }
      },
      "tool": "tools/equibind.json",
      "inputs": {
        "protein": {
          "class": "File",
          "filepath": "/Users/rindtorff/plex/testdata/binding/pdbbind_processed_size1/6d08/6d08_protein_processed.pdb"
        },
        "small_molecule": {
          "class": "File",
          "filepath": "/Users/rindtorff/plex/testdata/binding/pdbbind_processed_size1/6d08/6d08_ligand.sdf"
        }
      },
      "state": "processing",
      "errMsg": ""
    }
]
"""

data = json.loads(json_data)

# Validation classes and functions
class File(BaseModel):
    class_: str = Field(..., alias='class')
    filepath: FilePath

class Inputs(BaseModel):
    items: Dict[str, File]

    @validator('items', pre=True)
    def validate_files(cls, items):
        validator_dict = {
            name: globals().get(f"validate_{name}", None)
            for name in items.keys()
        }
        for name, file in items.items():
            validator_func = validator_dict.get(name)
            if validator_func:
                file = validator_func(file)
        return items

# Validate the inputs section
for item in data:
    inputs = item.get("inputs", None)
    if inputs:
        inputs_instance = Inputs(items=inputs)
        print(inputs_instance)
