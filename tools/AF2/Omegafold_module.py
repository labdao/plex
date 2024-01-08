import subprocess
import os

class Omegafold:
    def __init__(self, fasta_file, output_dir):
        self.cache_dir = '/root/.cache/omegafold_ckpt/'
        self.input_file = fasta_file
        self.output_dir = output_dir if output_dir else os.path.join(os.getcwd(), 'output')
        self.create_directories()
        if self.check_whether_weights_are_present()==False:
            print(f"Did not find model parameters in {self.cache_dir}. Will download them.")

    def create_directories(self):
        os.makedirs(self.cache_dir, exist_ok=True)
        os.makedirs(self.output_dir, exist_ok=True)
        print(f"Cache directory is {self.cache_dir}")
        print(f"Output directory is {self.output_dir}")

    def check_whether_weights_are_present(self):
        # Path to the directory where weights should be
        weights_dir = '/root/.cache/omegafold_ckpt/'

        # List of expected weight files
        expected_files = [
            "model.pt"
        ]

        # Check if the weights directory exists and contains all the expected files
        if os.path.exists(weights_dir):
            files_in_dir = set(os.listdir(weights_dir))
            return all(file in files_in_dir for file in expected_files)

        return False

    def run_prediction(self):
        
        print("Running prediction job...")
        work_dir = os.path.dirname(self.input_file)
        if not work_dir:
            work_dir = os.getcwd()  # Default to current directory if no directory is part of the input file path
        
        omegafold_command = "omegafold", f"{self.input_file}", f"{self.output_dir}"

        subprocess.run(omegafold_command, check=True)
        print(f"Prediction job complete. Results are in {self.output_dir}")

    def run(self):
        self.run_prediction()