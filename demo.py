import client
import process
import os
import json

BASEDIR = "/home/ubuntu/PDBBind_processed"

if __name__ == "__main__":
    task_list = os.listdir(BASEDIR)
    testset_df = pd.read_csv(BASEDIR + '/testset_csv.csv')
    filtered_task_list = [i.split('/')[-2] for i in df['protein_path'].tolist() if i.split('/')[-2] in task_list]
    print("generating docking instructions for diffdock")
    docking_instruct = client.generate_diffdock_instructions()
    for i in filtered_task_list:
        protein = i + '/' + i + '_protein.pdb'
        ligand = i + '/' + i + '_ligand.sdf'
        output = i + '/' + i + '_scored_vina.sdf.gz'
        instruction = client.generate_vina_instructions(
            protein=protein,
            ligand=ligand,
            output=output
        )
        print(instruction)
        instruction_list.append(instruction)
    # print("running docking")
    # process.main(docking_instruct)
    instruction_list = []
    print("generating scoring instructions for vina")
    for i in filtered_task_list:
        protein = i + '/' + i + '_protein.pdb'
        ligand = i + '/' + i + '_ligand.sdf'
        output = i + '/' + i + '_scored_vina.sdf.gz'
        instruction = client.generate_vina_instructions(
            protein=protein,
            ligand=ligand,
            output=output
        )
        print(instruction)
        instruction_list.append(instruction)
    print("generating scoring instructions for gnina")
    for i in filtered_task_list:
        protein = i + '/' + i + '_protein.pdb'
        ligand = i + '/' + i + '_ligand.sdf'
        output = i + '/' + i + '_scored_gnina.sdf.gz'
        instruction = client.generate_gnina_instructions(
            protein=protein,
            ligand=ligand,
            output=output
        )
        print(instruction)
        instruction_list.append(instruction)
    print("running scoring")
    for instruction in instruction_list:
        process.main(json.loads(instruction))