from biopandas.pdb import PandasPdb
from rdkit import Chem
from rdkit.Chem import Draw
from pydantic import validator

def validate_protein(file):
    filepath = file['filepath']
    
    if not filepath.endswith(".pdb"):
        raise ValueError(f"'protein' field: {filepath} is not a PDB file. Please ensure the file has a .pdb extension.")
    
    try:
        # Use BioPandas to read the PDB file
        ppdb = PandasPdb()
        ppdb.read_pdb(filepath)
    except Exception as e:
        raise ValueError(f"Invalid PDB file for 'protein' field: {filepath}. Error: {e}. Please ensure the file is a valid PDB file.")
    
    return file

def validate_small_molecule(file):
    filepath = file['filepath']
    
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
    
    return file
