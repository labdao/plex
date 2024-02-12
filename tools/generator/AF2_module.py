import subprocess
import os
import shutil

class AF2Runner:
    def __init__(self, fasta_file, output_dir):
        self.cache_dir = os.path.join(os.getcwd(), 'cache')
        self.input_file = fasta_file
        self.output_dir = output_dir if output_dir else os.path.join(os.getcwd(), 'output')
        self.create_directories()
        if self.check_whether_weights_are_present()==False:
            self.download_af2_weights()

    def create_directories(self):
        os.makedirs(self.cache_dir, exist_ok=True)
        os.makedirs(self.output_dir, exist_ok=True)
        print(f"Cache directory is {self.cache_dir}")
        print(f"Output directory is {self.output_dir}")

        for item in os.listdir(self.cache_dir):
            item_path = os.path.join(self.cache_dir, item)
            if item != 'colabfold' and os.path.isfile(item_path):
                os.remove(item_path)
            elif item != 'colabfold' and os.path.isdir(item_path):
                shutil.rmtree(item_path)

    def check_whether_weights_are_present(self):
        # Path to the directory where weights should be
        weights_dir = '/cache/colabfold/params'

        # List of expected weight files
        expected_files = [
            "params_model_1.npz", "params_model_1_ptm.npz", "params_model_1_multimer_v3.npz",
            "params_model_2.npz", "params_model_2_ptm.npz", "params_model_2_multimer_v3.npz",
            "params_model_3.npz", "params_model_3_ptm.npz", "params_model_3_multimer_v3.npz",
            "params_model_4.npz", "params_model_4_ptm.npz", "params_model_4_multimer_v3.npz",
            "params_model_5.npz", "params_model_5_ptm.npz", "params_model_5_multimer_v3.npz"
        ]

        # Check if the weights directory exists and contains all the expected files
        if os.path.exists(weights_dir):
            files_in_dir = set(os.listdir(weights_dir))
            return all(file in files_in_dir for file in expected_files)

        return False

    def download_af2_weights(self):
        # downloads weights
        if not os.listdir(self.cache_dir):
            print("Downloading AlphaFold2 weights...")
            subprocess.run(["python3", "-m", "colabfold.download"], check=True)

    def run_prediction(self):
        
        print("Running prediction job...")
        work_dir = os.path.dirname(self.input_file)
        if not work_dir:
            work_dir = os.getcwd()  # Default to current directory if no directory is part of the input file path
        
        # colabfold_batch_command = "colabfold_batch", f"/inputs/{os.path.basename(self.input_file)}", "/work/output"
        colabfold_batch_command = "colabfold_batch", f"{self.input_file}", f"{self.output_dir}"

        subprocess.run(colabfold_batch_command, check=True)
        print(f"Prediction job complete. Results are in {self.output_dir}")

    def run(self):
        self.run_prediction()
