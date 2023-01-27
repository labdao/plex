import client
import process
import os

if __name__ == "__main__":
    casf_list = os.listdir("/home/ubuntu/casf-2016")
    # print("generating docking instructions")
    # docking_instruct = client.generate_diffdock_instructions()
    # print("running docking")
    # process.main(docking_instruct)
    instruction_list = []
    print("generating scoring instructions")
    for i in casf_list:
        protein = i + '/' + i + '_protein.pdb'
        ligand = i + '/' + i + '_ligand.sdf'
        output = i + '/' + i + '_scored_vina.sdf.gz'
        instruction_list.append(client.generate_vina_instructions(
            protein=protein,
            ligand=ligand,
            output=output
        ))
    print("running scoring")
    for instruction in instruction_list:
        process.main(instruction)