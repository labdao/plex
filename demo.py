import client
import process

if __name__ == "__main__":
    print("generating docking instructions")
    docking_instruct = client.generate_diffdock_instructions()
    print("running docking")
    process.main(docking_instruct)
    print("generating scoring instructions")
    scoring_instruct = client.generate_vina_instructions()
    print("running scoring")
    process.main(scoring_instruct)
    #TODO: #30 add casf scoring