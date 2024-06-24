from pydantic import BaseModel, validator, FilePath, PositiveFloat, StrictBool, FiniteFloat, conint
from Bio.PDB import PDBParser
import re
from typing import Optional, List
import logging
import selfies as sf

class ProteinSequence(BaseModel):
    sequence: str
    hotspots: Optional[List[conint(ge=0)]] = None

    @validator('hotspots')
    def validate_hotspots(cls, v, values, **kwargs):
        if 'sequence' in values and v is not None:
            sequence_length = len(values['sequence'])
            v = list(set(v))  # remove duplicates
            if not all(0 <= hotspot < sequence_length for hotspot in v):
                raise ValueError('All hotspots must point to valid residues within the bounds of the sequence.')
        return v

    @validator('sequence')
    def replace_non_aminoacid_residues(cls, v):
        return re.sub('[^ACDEFGHIKLMNPQRSTVWY]', 'X', v.upper())
    

class SmallMoleculeSequence(BaseModel):
    smiles: Optional[str] = None
    selfies: Optional[str] = None
    enamine_id: Optional[str] = None

    #TODO implement
    # @validator('selfies')
    # def validate_selfies(cls, v):
    #     if v is not None and not re.match('^(\[[^\]]+\])+$', v):
    #         raise ValueError('selfies is expected to be a sequence of characters enclosed in square brackets')
    #     return v

    # @validator('smiles')
    # def validate_smiles(cls, v):
    #     if v is not None and not re.match('^[A-Za-z0-9@+\-\#\%\.=\(\)\[\]]+$', v):
    #         raise ValueError('smiles is expected to contain only alphanumeric characters, @, +, -, #, %, ., =, (, ), [, and ]')
    #     return v

    @validator('smiles', pre=True, always=True)
    def check_smiles_non_empty(cls, v, values):
        if 'selfies' in values and values['selfies'] is None and v is None:
            raise ValueError("Both 'smiles' and 'selfies' cannot be None. At least one must be non-empty.")
        return v

    @validator('selfies', pre=True, always=True)
    def check_selfies_non_empty(cls, v, values):
        if 'smiles' in values and values['smiles'] is None and v is None:
            raise ValueError("Both 'smiles' and 'selfies' cannot be None. At least one must be non-empty.")
        return v
    
    @property
    def to_smiles(self):
        if self.smiles is None and self.selfies is not None:
            self.smiles = sf.decoder(self.selfies)
        return self.smiles
    
    @property
    def to_selfies(self):
        if self.selfies is None and self.smiles is not None:
            self.selfies = sf.encoder(self.smiles)
        return self.selfies

    def to_dict(self):
        return {
            "smiles": self.smiles,
            "selfies": self.selfies,
            "enamine_id": self.enamine_id
        } 

#TODO: absorb into ProteinStructure
class BoundingBox(BaseModel):
    center_x: FiniteFloat
    center_y: FiniteFloat
    center_z: FiniteFloat
    size_x: PositiveFloat
    size_y: PositiveFloat
    size_z: PositiveFloat
    
    #TODO add hotspots to the ProteinSequence objects

class ProteinStructure(ProteinSequence):
    pdb: FilePath
    pdbqt: Optional[FilePath] = None
    #chain: Optional[str] #TODO check this
    bounding_box: Optional[BoundingBox] = None

    @classmethod
    def from_pdb(cls, pdb_path: str):
        sequence = cls._extract_sequence_from_pdb(pdb_path)
        return cls(sequence=sequence, pdb=pdb_path)

    @staticmethod
    def _extract_sequence_from_pdb(pdb_path: str) -> str:
        # This is a simplified example. You might need to adjust this to suit your needs.
        parser = PDBParser()
        structure = parser.get_structure('protein', pdb_path)
        sequence = ""
        for model in structure:
            for chain in model:
                for residue in chain:
                    sequence += residue.get_resname()
        return sequence

    #TODO create protein_structure from pdbqt and other file formats
    
    def create_bounding_box_from_hotspots(self) -> None:
        parser = PDBParser()
        structure = parser.get_structure('protein', self.pdb)
        
        coords = []
        for model in structure:
            for chain in model:
                for residue in chain:
                    if residue.id[1] in self.hotspots:
                        for atom in residue:
                            coords.append(atom.coord)
        
        coords = np.array(coords)
        min_coords = np.min(coords, axis=0)
        max_coords = np.max(coords, axis=0)
        center = (min_coords + max_coords) / 2
        size = max_coords - min_coords
        
        self.bounding_box = BoundingBox(
            center_x=center[0],
            center_y=center[1],
            center_z=center[2],
            size_x=size[0],
            size_y=size[1],
            size_z=size[2]
            )
    
    # @validator('pdb')
    # def validate_pdb_chains(cls, v, values):
    #     parser = PDBParser()
    #     structure = parser.get_structure('PDB', v)
    #     chains = list(structure.get_chains())
    #     if len(chains) > 1 and not values.get('chain'):
    #         logging.warning("PDB file contains more than one chain. A specific chain should be specified.")
    #     return v

    # @validator('pdb')
    # def validate_pdb_file(cls, v):
    #     if not v.suffix == '.pdb':
    #         raise ValueError('File must be a .pdb file')
    #     return v

    @validator('pdb')
    def validate_pdb(cls, v, values):
        # Check if the file has the correct suffix
        if not v.suffix == '.pdb':
            raise ValueError('File must be a .pdb file')

        # Parse the PDB file and check the number of chains
        parser = PDBParser()
        structure = parser.get_structure('PDB', v)
        chains = list(structure.get_chains())
        if len(chains) > 1 and not values.get('chain'):
            logging.warning("PDB file contains more than one chain. A specific chain should be specified.")

        return v

    @validator('pdbqt')
    def validate_pdb_file(cls, v):
        if v is not None and not v.suffix == '.pdbqt':
            raise ValueError('If provided, file must be a .pdbqt file')
        return v

class SmallMoleculeStructure(SmallMoleculeSequence):
    small_molecule_structure_sdf: FilePath

    @validator('small_molecule_structure_sdf')
    def validate_sdf_file(cls, v):
        if not v.suffix == '.sdf':
            raise ValueError('File must be an .sdf file')
        return v

    def to_dict(self):
        return {
            "smiles": self.smiles,
            "selfies": self.selfies,
            "small_molecule_structure_sdf": str(self.small_molecule_structure_sdf)
        }
    
class SmallMoleculeProteinComplexSequence(BaseModel):
    small_molecule_sequence: SmallMoleculeSequence
    protein_sequence: ProteinSequence

    def to_dict(self):
        return {
            "small_molecule_sequence": self.small_molecule_sequence,
            "protein_sequence": self.protein_sequence
        }

class SmallMoleculeProteinComplexStructure(BaseModel):
    isolated_small_molecule: SmallMoleculeStructure
    isolated_protein: ProteinStructure
    complex_small_molecule_structure: Optional[FilePath] = None

    @validator('complex_small_molecule_structure')
    def validate_complex_small_molecule_structure(cls, v):
        if v is not None and not v.suffix == '.sdf':
            raise ValueError('If provided, file must be a .sdf file')
        return v

    def to_dict(self):
        return {
            "isolated_small_molecule": self.isolated_small_molecule,
            "isolated_protein": self.isolated_protein,
            "complex_small_molecule_structure": self.complex_small_molecule_structure
        }


class ProteinBinderTargetSequence(BaseModel):
    binder_sequence: str
    target_sequence: str


    @classmethod
    def replace_non_aminoacid_residues(cls, v):
        return re.sub('[^ACDEFGHIKLMNPQRSTVWY]', 'X', v.upper())

    @validator('binder_sequence', 'target_sequence')
    def validate_sequences(cls, v):
        return cls.replace_non_aminoacid_residues(v)

    @property
    def is_binder_sequence_fully_undetermined(self) -> StrictBool:
        return StrictBool(self.binder_sequence == 'X' * len(self.binder_sequence))

class ProteinBinderTargetStructure(ProteinBinderTargetSequence):
    pdb: FilePath
    # binder_sequence: str
    # target_sequence: str
    # binder_sequence_complete: Optional[StrictBool] = None
    # target_sequence_complete: Optional[StrictBool] = None
    
    # @validator('binder_chain', 'target_chain')
    # def validate_single_capital_letter(cls, v):
    #     if not re.match('^[A-Z]$', v):
    #         raise ValueError('Chain must be a single capital letter')
    #     return v

    @validator('pdb')
    def validate_pdb_chains(cls, v, values, **kwargs):
        parser = PDBParser()
        structure = parser.get_structure('PDB', v)
        chains = list(structure.get_chains())
        if len(chains) < 2:
            raise ValueError("PDB file must contain at least two chains.")
        elif len(chains) >= 2:
            if not values.get('binder_chain') or not values.get('target_chain'):
                logging.warning("The pdb file contains multiple chains. You may want to specify a target and a binder chain.")
        return v