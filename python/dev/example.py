import os

from plex import CoreTools, ScatteringMethod, plex_init, plex_run, plex_vectorize, plex_mint

plex_dir = os.path.dirname(os.path.dirname(os.getcwd()))
plex_path = os.path.join(plex_dir, "plex")
jobs_dir = os.path.join(plex_dir, "jobs")
test_data_dir = os.path.join(plex_dir, "testdata")

print(f"Using plex_path, {plex_path}, if this looks incorrect then make sure you are running from the python/dev directory")

small_molecules = [f"{test_data_dir}/binding/abl/ZINC000003986735.sdf", f"{test_data_dir}/binding/abl/ZINC000019632618.sdf"]
proteins = [f"{test_data_dir}/binding/abl/7n9g.pdb"]

initial_io_cid = plex_init(
    CoreTools.EQUIBIND.value,
    ScatteringMethod.CROSS_PRODUCT.value,
    plex_path=plex_path,
    small_molecule=small_molecules,
    protein=proteins)

# Custom annotations for testing
custom_annotations = ["python_example", "test"]

completed_io_cid, io_file_path = plex_run(initial_io_cid, output_dir=jobs_dir, annotations=custom_annotations, plex_path=plex_path)

# Print annotations to verify
print(f"\nAnnotations used in plex_run: {custom_annotations}\n")

vectors = plex_vectorize(io_file_path, CoreTools.EQUIBIND.value, plex_path=plex_path)

print(vectors)
print(vectors['best_docked_small_molecule']['filePaths'])
print(vectors['best_docked_small_molecule']['cids'])

plex_mint(completed_io_cid, plex_path=plex_path)
