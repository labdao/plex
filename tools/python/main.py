import json
from io_model import IOModel

# Load the JSON data from the example file
with open("io_example.json", "r") as f:
    example_data = json.load(f)

# Create an IOModel instance from the loaded data
io_instance = IOModel.parse_obj(example_data[0])

# Change the filepath of the small molecule input
io_instance.inputs["small_molecule"]["filepath"] = "/new/path/to/small_molecule.sdf"

# Update the example data with the modified IOModel instance
example_data[0] = io_instance.dict()

# Save the updated data back to the JSON file
with open("io_example_updated.json", "w") as f:
    json.dump(example_data, f, indent=2)
