from pydantic import BaseModel, FilePath
from validators import validate_files
from typing import Dict 

class File(BaseModel):
    class_: str
    filepath: FilePath

class Inputs(BaseModel):
    items: Dict[str, File]  # Define a field to hold the dictionary of File models

    validate_files = validate_files

if __name__ == "__main__":
    # Example usage
    example_inputs = {
        "protein": {
            "class": "File",
            "filepath": "/Users/rindtorff/plex/8ae8b6c2-77cc-4201-a5b6-0c0ede451acd/entry-0/inputs/6d08_protein_processed.pdb"
        },
        "small_molecule": {
            "class": "File",
            "filepath": "/Users/rindtorff/plex/8ae8b6c2-77cc-4201-a5b6-0c0ede451acd/entry-0/inputs/6d08_ligand.sdf"
        }
    }

    # Create an Inputs instance from the example data
    inputs_instance = Inputs(items=example_inputs)  # Pass the value of the "inputs" key directly

    # The inputs_instance contains the validated data
    print(inputs_instance)
