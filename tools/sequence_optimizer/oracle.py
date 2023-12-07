import os
from AF2_module import AF2Runner

def write_dataframe_to_fastas(dataframe, cfg):
    input_dir = os.path.join(cfg.inputs.directory, 'current_sequences')
    if os.path.exists(input_dir):
        # If the folder already exists, empty the folder of all files
        for file_name in os.listdir(input_dir):
            file_path = os.path.join(input_dir, file_name)
            if os.path.isfile(file_path):
                os.remove(file_path)
    else:
        os.makedirs(input_dir, exist_ok=True)

    for index, row in dataframe.iterrows():
        file_path = os.path.join(input_dir, f"seq_{row['sequence_number']}.fasta")
        with open(file_path, 'w') as file:
            file.write(f">{row['sequence_number']}\n{row['seq']}\n")
    return os.path.abspath(input_dir)

class Oracle:
    def __init__(self, df, outputs_directory, cfg):

        self.df = df
        self.outputs_directory = outputs_directory
        self.cfg = cfg

    def run(self):

        # prepare input sequences as fastas and run AF2 K-times
        seq_input_dir = write_dataframe_to_fastas(self.df, self.cfg)

        K = self.cfg.params.basic_settings.AF2_repeats_per_seq
        for n in range(K):
            print("starting repeat number ", n)
            af2_runner = AF2Runner(seq_input_dir, self.outputs_directory)
            af2_runner.run()
        
        return None