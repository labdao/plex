import json
from typing import List
from objects import Protein, SmallMolecule, Inputs  # Import the necessary classes

def equibind(path_to_json: str, protein_filepaths: List[str], small_molecule_filepaths: List[str]) -> List[dict]:
    # Load the JSON template
    with open(path_to_json, "r") as f:
        data_template = json.load(f)
    
    # Initialize an empty list to store the generated JSON dictionaries
    json_data_list = []
    
    # Iterate over the protein and small molecule filepaths
    for protein_filepath, small_molecule_filepath in zip(protein_filepaths, small_molecule_filepaths):
        # Create instances of the Protein and SmallMolecule classes
        protein = Protein(filepath=protein_filepath)
        small_molecule = SmallMolecule(filepath=small_molecule_filepath)
        
        # Create an instance of the Inputs class
        inputs = Inputs(
            protein=protein,
            small_molecule=small_molecule
        )
        
        # Create a copy of the data template
        data = data_template.copy()
        
        # Replace the "inputs" section with the serialized Inputs instance
        data[0]["inputs"] = json.loads(inputs.json())
        
        # Append the generated JSON dictionary to the list
        json_data_list.append(data)
    
    # Save the generated JSON data to a file named 'equibind-io-output.json'
    with open("equibind-io-output.json", "w") as outfile:
        json.dump(json_data_list, outfile, indent=2)
    
    return json_data_list

if __name__ == "__main__":
    # Example usage
    protein_filepaths = [
        "/Users/rindtorff/plex/testdata/binding/pdbbind_processed_size1/6d08/6d08_protein_processed.pdb",
        "/Users/rindtorff/plex/testdata/binding/abl/7n9g.pdb",
        # Additional protein filepaths...
    ]
    small_molecule_filepaths = [
        "/Users/rindtorff/plex/testdata/binding/abl/ZINC000003986735.sdf",
        "/Users/rindtorff/plex/testdata/binding/abl/ZINC000019632618.sdf",
        # Additional small molecule filepaths...
    ]
    
    # Call the equibind function
    json_data_list = equibind("equibind_io.json", protein_filepaths, small_molecule_filepaths)
    
    # Print the generated JSON data
    for json_data in json_data_list:
        print(json.dumps(json_data, indent=2))
