# io_model.py
from pydantic import BaseModel, FilePath
from typing import Dict, Any
from inputs_model import Inputs, File  # Import the Inputs and File models

class IOModel(BaseModel):
    inputs: Inputs  # Use the Inputs model
    outputs: Dict[str, Any]
    tool: str
    state: str
    errMsg: str

# Function to load JSON data, modify file path, and save updated data
def modify_io_model(json_data, new_file_path):
    # Create an IOModel instance from the loaded data
    io_instance = IOModel.parse_obj(json_data[0])

    # Change the filepath of the small molecule input
    io_instance.inputs.items["small_molecule"]["filepath"] = new_file_path

    # Update the JSON data with the modified IOModel instance
    json_data[0] = io_instance.dict()

    return json_data

if __name__ == "__main__":
    import json

    # Load the JSON data from the example file
    with open("io_example.json", "r") as f:
        example_data = json.load(f)

    # Modify the file path of the small molecule input
    updated_data = modify_io_model(example_data, "/new/path/to/small_molecule.sdf")

    # Save the updated data back to the JSON file
    with open("io_example_updated.json", "w") as f:
        json.dump(updated_data, f, indent=2)
