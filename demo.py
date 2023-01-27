import client
import process

if __name__ == "__main__":
    print("generating diffdock instructions")
    diffdock_instruct = client.generate_diffdock_instructions()
    print("running diffdock")
    process.main(diffdock_instruct)
    #TODO: #29 add vina scoring
    #TODO: #30 add casf scoring