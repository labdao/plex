from pydantic import BaseModel, Field, FilePath, validator
from typing import Dict, Optional, List, Union, Any
import os
import json
from fastapi.encoders import jsonable_encoder
from tools import Tool

class File(BaseModel):
    class_: str = Field("File", alias='class')
    filepath: FilePath

# TODO #300 add class Array
# class Array(BaseModel):


class Int(BaseModel):
    class_: str = Field("Int", alias='class')
    value: int

class String(BaseModel):
    class_: str = Field("String", alias='class')
    value: str

class IO(BaseModel):
    inputs: Dict[str, Any]
    outputs: Dict[str, Any]
    tool: str
    state: str
    errMsg: str

predefined_classes = {
    "File": File,
    "Int": Int,
    "String": String
}

def generate_dynamic_models(tool_config: Tool):
    class CustomInputs(BaseModel):
        pass

    class CustomOutputs(BaseModel):
        pass

    for key, config in tool_config.inputs.items():
        cls = predefined_classes.get(config["type"], Any)
        default = config.get("default", None)
        setattr(CustomInputs, key, Field(default, __orig_bases__=(cls,)))

    for key, config in tool_config.outputs.items():
        cls = predefined_classes.get(config["type"], Any)
        default = config.get("default", None)
        setattr(CustomOutputs, key, Field(default, __orig_bases__=(cls,)))

    return CustomInputs, CustomOutputs

if __name__ == "__main__":
    with open("tools/equibind.json", "r") as f:
        tool_config_data = json.load(f)

    tool_config = Tool(**tool_config_data)
    print("################# Tool:")
    print(tool_config)

    print("################# Custum Input:")
    CustomInputs, CustomOutputs = generate_dynamic_models(tool_config)
    print(CustomInputs)

    print("Attributes and methods of CustomInputs:")
    custom_inputs = CustomInputs(protein=File(filepath='/Users/rindtorff/plex/testdata/binding/pdbbind_processed_size1/6d08/6d08_protein_processed.pdb'), 
                                 small_molecule=File(filepath='/Users/rindtorff/plex/testdata/binding/abl/ZINC000003986735.sdf'))

    
    print("Dynamic Inputs instance:")
    print(custom_inputs)
    print("Dynamic Inputs instance, serialized for protein:")
    print(custom_inputs.protein)

    iomodel = IO(
        inputs=jsonable_encoder(custom_inputs),
        outputs={},
        tool="tools/equibind.json",
        state="processing",
        errMsg="",
    )

    print(iomodel)
