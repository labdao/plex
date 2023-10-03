import os

from plex import CoreTools, ScatteringMethod, plex_init, plex_run, plex_vectorize, plex_mint

plex_dir = os.path.dirname(os.path.dirname(os.getcwd()))
plex_path = os.path.join(plex_dir, "plex")
jobs_dir = os.path.join(plex_dir, "jobs")
test_data_dir = os.path.join(plex_dir, "testdata")

print(f"Using plex_path, {plex_path}, if this looks incorrect then make sure you are running from the python/dev directory")

small_molecules = [f"{test_data_dir}/binding/abl/ZINC000003986735.sdf", f"{test_data_dir}/binding/abl/ZINC000019632618.sdf"]
proteins = [f"{test_data_dir}/binding/abl/7n9g.pdb"]

# Custom annotations for testing
custom_annotations = ["python_example", "test"]

print(f"Testing plex_run with max_time set to 60 seconds...")
print(f"Expected behavior is failure with diffdock timeout...")

test_io_cid = plex_init(
    CoreTools.DIFFDOCK.value,
    ScatteringMethod.CROSS_PRODUCT.value,
    plex_path=plex_path,
    small_molecule=small_molecules,
    protein=proteins,
)

print(f"Test IO CID: {test_io_cid}")

# max_time set to 1 min to test timeout
completed_io_cid, io_file_path = plex_run(
    test_io_cid, 
    output_dir=jobs_dir,
    max_time="1",
    annotations=custom_annotations, 
    plex_path=plex_path, 
)

print(f"Testing plex_init with auto_run flag set to True")

test_io_cid = plex_init(
    CoreTools.EQUIBIND.value,
    ScatteringMethod.CROSS_PRODUCT.value,
    auto_run=True,
    plex_path=plex_path,
    small_molecule=small_molecules,
    protein=proteins
)

print(f"Testing plex_init with auto_run flag set to False")

initial_io_cid = plex_init(
    CoreTools.EQUIBIND.value,
    ScatteringMethod.CROSS_PRODUCT.value,
    plex_path=plex_path,
    small_molecule=small_molecules,
    protein=proteins
)

# check that environmental variable for recipient wallet is set
if not os.environ.get('RECIPIENT_WALLET'):
    print("RECIPIENT_WALLET environment variable not set")
else:
    print(f"RECIPIENT_WALLET environment variable set to {os.environ.get('RECIPIENT_WALLET')}")

completed_io_cid, io_file_path = plex_run(initial_io_cid, output_dir=jobs_dir, annotations=custom_annotations, plex_path=plex_path)

# Print annotations to verify
print(f"\nAnnotations used in plex_run: {custom_annotations}\n")

vectors = plex_vectorize(io_file_path, CoreTools.EQUIBIND.value, plex_path=plex_path)

print(vectors)
print(vectors['best_docked_small_molecule']['filePaths'])
print(vectors['best_docked_small_molecule']['cids'])

plex_mint(completed_io_cid, plex_path=plex_path)
