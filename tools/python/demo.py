import json
from typing import Dict, List, Any
from pydantic import BaseModel, FilePath, Field
from pydantic import validator
from validators import validate_protein, validate_small_molecule

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
    protein: File
    small_molecule: File

    @validator('protein', pre=True)
    def validate_protein(cls, file):
        print("Validating protein")
        return validate_protein(file)

    @validator('small_molecule', pre=True)
    def validate_small_molecule(cls, file):
        print("Validating small_molecule")
        return validate_small_molecule(file)

class IOModel(BaseModel):
    inputs: Inputs  # Use the Inputs model
    outputs: Dict[str, Any]
    tool: str
    state: str
    errMsg: str

# Validate the entire IOModel
for item in data:
    io_instance = IOModel.parse_obj(item)
    print(io_instance)
