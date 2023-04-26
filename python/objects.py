from pydantic import BaseModel, Field, FilePath, validator
from typing import Dict, Optional
import os
# internal imports
# from engine_components import Inputs, IOModel
# domain-specific imports
from biopandas.pdb import PandasPdb
from rdkit import Chem
# demo
import json

class File(BaseModel):
    class_: str = Field("File", alias='class')
    filepath: FilePath

class Protein(File):
    skip_validation: bool = False  # New attribute to control validation

    @validator("filepath", pre=True)
    def validate_protein_graph(cls, filepath, values):
        # Skip validation if the 'skip_validation' attribute is set to True
        if values.get('skip_validation', True):
            return filepath
        
        if not filepath.endswith(".pdb"):
            raise ValueError(f"'protein' field: {filepath} is not a PDB file. Please ensure the file has a .pdb extension.")
        try:
            # Use BioPandas to read the PDB file
            ppdb = PandasPdb()
            ppdb.read_pdb(filepath)
        except Exception as e:
            raise ValueError(f"Invalid PDB file for 'protein' field: {filepath}. Error: {e}. Please ensure the file is a valid PDB file.")
        return filepath

class SmallMolecule(File):
    skip_validation: bool = False  # New attribute to control validation

    @validator("filepath", pre=True)
    def validate_small_molecule(cls, filepath, values):
        # Skip validation if the 'skip_validation' attribute is set to True
        if values.get('skip_validation', True):
            return filepath
        
        if not filepath.endswith(".sdf"):
            raise ValueError(f"'small_molecule' field: {filepath} is not an SDF file. Please ensure the file has a .sdf extension.")
        try:
            # Use RDKit to read the SDF file
            suppl = Chem.SDMolSupplier(filepath)
            # Iterate over the molecules in the SDF file
            for mol in suppl:
                if mol is None:
                    raise ValueError(f"Invalid molecule in SDF file for 'small_molecule' field: {filepath}. Please ensure the file contains valid molecules.")
                # Optionally, you can perform additional validation on the molecule here
        except Exception as e:
            raise ValueError(f"Invalid SDF file for 'small_molecule' field: {filepath}. Error: {e}. Please ensure the file is a valid SDF file.")
        return filepath


# TODO need generalisable composability of inputs, specific to the tool that is being used
class Inputs(BaseModel):
    protein: Protein
    small_molecule: SmallMolecule

if __name__ == "__main__":
    # Example usage
    protein_6d08 = Protein(filepath="/Users/rindtorff/plex/testdata/binding/pdbbind_processed_size1/6d08/6d08_protein_processed.pdb")
    ligand_abl = SmallMolecule(filepath="/Users/rindtorff/plex/testdata/binding/abl/ZINC000003986735.sdf")

    inputs = Inputs(
        protein= protein_6d08,
        small_molecule=ligand_abl
    )

    print("Inputs instance:")
    print(inputs)
    print("Inputs instance, serialized for protein:")
    print(inputs.protein)

    # Define the JSON object
    data = [
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
            "inputs": {},  # Will be replaced with the Inputs instance
            "state": "processing",
            "errMsg": ""
        }
    ]
    print("JSON object:")
    print(data)

    # Replace the "inputs" section with the serialized Inputs instance
    data[0]["inputs"] = json.loads(inputs.json())
    json_data = json.dumps(data, indent=2)
    print("JSON object, edited and serialized:")
    print(json_data)

    # Save the data to a JSON file
    with open("example.json", "w") as outfile:
        json.dump(data, outfile, indent=2)