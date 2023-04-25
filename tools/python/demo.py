from typing import List
from pydantic import BaseModel, FilePath, Field
import json
from typing import Dict, List, Any
from pydantic import validator
from validators import validate_protein, validate_small_molecule

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

    def update_filepaths(self, **kwargs) -> None:
        for key, value in kwargs.items():
            if hasattr(self.inputs, key):
                setattr(self.inputs, key, self.inputs.__getattribute__(key).copy(update={'filepath': value}))
            else:
                raise ValueError(f"Invalid key: {key}. Cannot update filepath.")
    
    def run(self, **kwargs) -> None:
        # Update the filepaths
        self.update_filepaths(**kwargs)
        
        # Run the tool
        # TODO call plex.run()
        print("I would be running the tool now. Below is the updated IOModel:")
        

def get(json_file: str) -> IOModel:
    # Read the JSON data from the file
    with open(json_file, 'r') as f:
        data = json.load(f)
    
    # Select the first dictionary in the list
    first_item = data[0]
    
    # Validate the dictionary and return an IOModel instance
    io_instance = IOModel.parse_obj(first_item)
    return io_instance

# Example usage
io_example = get(json_file='io_example.json')
print(io_example)

io_example.run(
    protein='/Users/rindtorff/plex/testdata/binding/pdbbind_processed_size1/6d08/6d08_protein_processed.pdb',
    small_molecule='/Users/rindtorff/plex/testdata/binding/abl/ZINC000003986735.sdf'
)
print(io_example)
