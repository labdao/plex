import client
import process
import os
import json
import pandas as pd

BASEDIR = "/home/ubuntu/inputs/"
TASKDIR = "PDBBind_processed"

if __name__ == "__main__":
    task_list = os.listdir(BASEDIR + TASKDIR)
    testset_df = pd.read_csv(BASEDIR + TASKDIR + '/testset_csv.csv')
    filtered_task_list = [i.split('/')[-2] for i in testset_df['protein_path'].tolist() if i.split('/')[-2] in task_list]
    instruction_list = []
    print("generating docking instructions for gnina")
    for i in filtered_task_list:
        protein = TASKDIR + '/' + i + '/' + i + '_protein_processed.pdb'
        ligand = TASKDIR + '/' + i + '/' + i + '_ligand.sdf'
        output = TASKDIR + '/' + i + '/' + i + '_scored_gnina.sdf.gz'
        instruction = client.generate_gnina_instructions(
            protein=protein,
            ligand=ligand,
            output=output,
        )
        print(instruction)
        instruction_list.append(instruction)
    print("generating docking instructions for diffdock")
    for i in filtered_task_list:
        protein = TASKDIR + '/' + i + '/' + i + '_protein_processed.pdb'
        ligand = TASKDIR + '/' + i + '/' + i + '_ligand.sdf'
        output = TASKDIR + '/' + i
        instruction = client.generate_diffdock_instructions(
            protein=protein,
            ligand=ligand,
            output=output
        )
        print(instruction)
        instruction_list.append(instruction)
    print("generating scoring instructions for vina")
    for i in filtered_task_list:
        protein = TASKDIR + '/' + i + '/' + i + '_protein_processed.pdb'
        ligand = TASKDIR + '/' + i + '/' + i + '_ligand.sdf'
        output = TASKDIR + '/' + i + '/' + i + '_scored_vina.sdf.gz'
        instruction = client.generate_vina_instructions(
            protein=protein,
            ligand=ligand,
            output=output,
            modifier='--score_only'
        )
        print(instruction)
        instruction_list.append(instruction)
    print("generating scoring instructions for gnina")
    for i in filtered_task_list:
        protein = TASKDIR + '/' + i + '/' + i + '_protein_processed.pdb'
        ligand = TASKDIR + '/' + i + '/' + i + '_ligand.sdf'
        output = TASKDIR + '/' + i + '/' + i + '_scored_gnina.sdf.gz'
        instruction = client.generate_gnina_instructions(
            protein=protein,
            ligand=ligand,
            output=output,
            modifier='--score_only'
        )
        print(instruction)
        instruction_list.append(instruction)
    print("running instructions")
    for instruction in instruction_list:
        process.main(json.loads(instruction))