from biopandas.pdb import PandasPdb
from rdkit import Chem
from rdkit.Chem import Draw
from pydantic import validator

def validate_protein(file):
    filepath = file['filepath']
    try:
        # Use BioPandas to read the PDB file
        ppdb = PandasPdb()
        ppdb.read_pdb(filepath)
    except Exception as e:
        raise ValueError(f"Invalid PDB file: {e}")
    return file

def validate_small_molecule(file):
    filepath = file['filepath']
    try:
        # Use RDKit to read the SDF file
        suppl = Chem.SDMolSupplier(filepath)
        # Iterate over the molecules in the SDF file
        for mol in suppl:
            if mol is None:
                raise ValueError("Invalid molecule in SDF file")
            # Optionally, you can perform additional validation on the molecule here
    except Exception as e:
        raise ValueError(f"Invalid SDF file: {e}")
    return file

# Add more validator functions for other file types as needed

@validator('items', pre=True)
def validate_files(cls, items):
    # Create a dictionary of validator functions
    validator_dict = {
        name: globals().get(f"validate_{name}", None)
        for name in items.keys()
    }
    # Apply specific validation rules based on the name of the file
    for name, file in items.items():
        validator_func = validator_dict.get(name)
        if validator_func:
            file = validator_func(file)
    return items
