from pydantic import BaseModel, Field, FilePath, validator
from typing import Dict, Optional, List, Union, Any
import os
import json
from fastapi.encoders import jsonable_encoder

# Import the Protein and SmallMolecule classes from the objects module
from objects import Protein, SmallMolecule, File

# Define the Tool class
class Tool(BaseModel):
    class_: str = Field(alias="class")
    name: str
    description: str
    baseCommand: List[str]
    arguments: List[str]
    dockerPull: str
    gpuBool: bool
    networkBool: bool
    inputs: Dict[str, Dict[str, Union[str, List[str]]]]
    outputs: Dict[str, Dict[str, Union[str, List[str]]]]

class IOModel(BaseModel):
    inputs: Dict[str, Any]  # Use a dictionary to store dynamic inputs
    outputs: Dict[str, Any]
    tool: str
    state: str
    errMsg: str

# Define a dictionary containing the predefined classes
predefined_classes = {
    "protein": Protein,
    "small_molecule": SmallMolecule,
    "File": File,  # Add a mapping for the "File" type
}


def generate_fields(field_definitions: Dict[str, Dict[str, Any]]) -> Dict[str, Any]:
    fields = {}
    for key, value in field_definitions.items():
        if value["type"] in predefined_classes:
            fields[key] = (Field(..., type_=predefined_classes[value["type"]]))
        else:
            fields[key] = (Field(..., type_=eval(value["type"])))
    return fields

def generate_dynamic_io_models(tool_config: Tool):
    # Generate input fields
    input_fields = generate_fields(tool_config.inputs)
    # Generate the CustomInputs class dynamically
    CustomInputs = type("CustomInputs", (BaseModel,), input_fields)

    # Generate output fields
    output_fields = generate_fields(tool_config.outputs)
    # Generate the CustomOutputs class dynamically
    CustomOutputs = type("CustomOutputs", (BaseModel,), output_fields)

    return CustomInputs, CustomOutputs




# NOT RUN
if __name__ == "__main__":
    # Load the tool configuration from a JSON file
    with open("../equibind.json", "r") as f:
        tool_config_data = json.load(f)

    # Create an instance of the Tool class
    tool_config = Tool(**tool_config_data)
    print(tool_config)

    # Generate the custom input and output models
    CustomInputs, CustomOutputs = generate_dynamic_io_models(tool_config)

    # Define inputs 
    protein_file = Protein(filepath="/Users/rindtorff/plex/testdata/binding/pdbbind_processed_size1/6d08/6d08_protein_processed.pdb")
    ligand_file = SmallMolecule(filepath="/Users/rindtorff/plex/testdata/binding/abl/ZINC000003986735.sdf")

    print("protein:")
    print(protein_file)

    print("Attributes and methods of CustomInputs:")
    #print(CustomInputs.schema())

    # Create instances of the custom input and output models
    custom_inputs = CustomInputs(
        protein=protein_file,
        small_molecule=ligand_file,
    )
    
    print("Dynamic Inputs instance:")
    print(custom_inputs)
    print("Dynamic Inputs instance, serialized for protein:")
    print(custom_inputs.protein)

    custom_outputs = CustomOutputs(
        best_docked_small_molecule={"class": "File", "filepath": ""},
        protein={"class": "File", "filepath": ""},
    )

    # Create an instance of the IOModel class
    iomodel = IOModel(
        inputs=jsonable_encoder(custom_inputs),
        outputs=jsonable_encoder(custom_outputs),
        tool="tools/equibind.json",
        state="processing",
        errMsg="",
    )

    print(iomodel)
