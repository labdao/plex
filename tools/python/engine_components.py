from pydantic import BaseModel, Field, FilePath, validator
from typing import Dict, Optional
import os

class IOModel(BaseModel):
    inputs: Inputs  # Use the Inputs model
    outputs: Dict[str, Any]
    tool: str
    state: str
    errMsg: str
