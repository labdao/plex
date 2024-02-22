import subprocess
import os
import logging

class AF2Runner:
    def __init__(self, fasta_file, output_dir):

        self.cache_dir = "/cache/colabfold/params"
        self.input_file = fasta_file
        self.output_dir = output_dir # if output_dir else os.path.join(os.getcwd(), 'outputs')

    def run_prediction(self):
        
        work_dir = os.path.dirname(self.input_file)
        if not work_dir:
            work_dir = os.getcwd()
        
        colabfold_batch_command = "colabfold_batch", f"{self.input_file}", f"{self.output_dir}"

        subprocess.run(colabfold_batch_command, check=True)
        logging.info(f"Colabfold prediction complete. Results are in {self.output_dir}")

    def run(self):
        self.run_prediction()
