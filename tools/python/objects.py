from pydantic import BaseModel, FilePath, validator
from typing import Dict, Optional
import os
# domain-specific imports
from biopandas.pdb import PandasPdb
from rdkit import Chem

class BaseObject(BaseModel):
    pass

# Todo rename to ProteinGraph
class Protein(BaseObject):
    filepath: Optional[FilePath] = None

    @validator("filepath", pre=True)
    def validate_protein_graph(cls, filepath):
        if filepath is not None:
            if not os.path.isfile(filepath):
                raise ValueError(f"'protein' field: {filepath} does not exist or is not a file. Please provide a valid file.")
            if not filepath.endswith(".pdb"):
                raise ValueError(f"'protein' field: {filepath} is not a PDB file. Please ensure the file has a .pdb extension.")
            try:
                # Use BioPandas to read the PDB file
                ppdb = PandasPdb()
                ppdb.read_pdb(filepath)
            except Exception as e:
                raise ValueError(f"Invalid PDB file for 'protein' field: {filepath}. Error: {e}. Please ensure the file is a valid PDB file.")
        return filepath

class SmallMolecule(BaseObject):
    filepath: Optional[FilePath] = None
    pass

# Example usage
protein = Protein(filepath="/Users/rindtorff/plex/testdata/binding/pdbbind_processed_size1/6d08/6d08_protein_processed.pdb")
small_molecule = SmallMolecule(filepath="/Users/rindtorff/plex/testdata/binding/abl/ZINC000003986735.sdf")
print(protein)

