from pydantic import BaseModel, Field, FilePath, validator
from typing import Dict, Optional, List, Union, Any
import json

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