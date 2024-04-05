import os
import json
import time
import boto3
import random


def get_plex_job_inputs():
    # Retrieve the environment variable
    json_str = os.getenv("PLEX_JOB_INPUTS")
    job_uuid = os.getenv("JOB_UUID")
    flow_uuid = os.getenv("FLOW_UUID")
    checkpoint_compatible = os.getenv("CHECKPOINT_COMPATIBLE", "False")

    # Check if the environment variable is set
    if (job_uuid is None) or (flow_uuid is None) or (checkpoint_compatible is None) or (json_str is None):
        raise ValueError("One/more of the mandatory environment variables are missing: PLEX_JOB_INPUTS/JOB_UUID/FLOW_UUID/CHECKPOINT_COMPATIBLE.")

    # Convert the JSON string to a Python dictionary
    try:
        data = json.loads(json_str)
        return data, job_uuid, flow_uuid, checkpoint_compatible
    except json.JSONDecodeError:
        # Handle the case where the string is not valid JSON
        raise ValueError("PLEX_JOB_INPUTS is not a valid JSON string.")

def upload_to_s3(file_name, bucket_name, object_name=None):
    if object_name is None:
        object_name = file_name

    s3_client = boto3.client('s3')
    try:
        s3_client.upload_file(file_name, bucket_name, object_name)
        print(f"Successfully uploaded {file_name} to {bucket_name}/{object_name}")
    except Exception as e:
        print(f"Failed to upload {file_name} to {bucket_name}/{object_name}: {e}")
        raise e

def create_event_csv(checkpoint_number, job_inputs):
    file_name = f"checkpoint_{checkpoint_number}_summary.csv"
    with open(file_name, 'w') as file:
        file.write("cycle,proposal,plddt,i_pae,dim1,dim2,pdbFileName\n")
        # Hardcoded data lines for each checkpoint
        if checkpoint_number == 0:
            checkpoint_pdb_filepath = job_inputs["pdb_checkpoint_0"]
            data_line = f"1,1,{random.randint(1, 100)},{random.randint(1, 40)},5,5,{os.path.basename(checkpoint_pdb_filepath)}\n"
        elif checkpoint_number == 1:
            checkpoint_pdb_filepath = job_inputs["pdb_checkpoint_1"]
            data_line = f"2,2,{random.randint(1, 100)},{random.randint(1, 40)},11,3,{os.path.basename(checkpoint_pdb_filepath)}\n"
        elif checkpoint_number == 2:            
            checkpoint_pdb_filepath = job_inputs["pdb_checkpoint_2"]
            data_line = f"3,3,{random.randint(1, 100)},{random.randint(1, 40)},9,12,{os.path.basename(checkpoint_pdb_filepath)}\n"
        else:
            data_line = ""
            checkpoint_pdb_filepath = ""
        file.write(data_line)
    return file_name, checkpoint_pdb_filepath

def display_logo(speedup=False):
    logo = """
                                        @
                                 @@@@@@@@@@@@@@@
                               @@@@@@@@@@@@@@@@@@@
                              @@@@@@@@@@@@@@@@@@@@@
             @@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@
           @@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@
         @@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@@@
        *@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@         @@@@@@@@@@@@@
         @@@@@@@@@@        @@@@@@@@@@@@@@@@@@@@@%            &@@@@@@@@@@
           @@@@           @@@@@@@@@@@@@@@@@@&                     @@@@
                        @@@@@@@@
                   @@@@@@@@@
      @@@@@@@@@@@@@@@@@@@@        ,@@@@@@@@@@@                 @@@@@@@@@@@@
   @@@@@@@@@@@@@@@@@@@@@@       @@@@@@@@@@@@@@@@@           @@@@@@@@@@@@@@@@@@
  @@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@       @@@@@@@@@@@@@@@@@@@@@
 @@@@@@@@@@@@@@@@@@@@@@@     @@@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@     @@@@@@@@@@@@@@@@@@@@@@@     @@@@@@@@@@@@@@@@@@@@@@@
 @@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@@     @@@@@@@@@@@@@@@@@@@@@@@
  @@@@@@@@@@@@@@@@@@@@@       @@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@
   @@@@@@@@@@@@@@@@@@           @@@@@@@@@@@@@@@@@       @@@@@@@@@@@@@@@@@@@@@@
      @@@@@@@@@@@@                 @@@@@@@@@@@         @@@@@@@@@@@@@@@@@@@@
                                                     @@@@@@@@@
                                                 @@@@@@@@
           @@@@                     &@@@@@@@@@@@@@@@@@@           @@@@
         @@@@@@@@@@             @@@@@@@@@@@@@@@@@@@@@        &@@@@@@@@@@
        *@@@@@@@@@@@@@        @@@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@@
         @@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@@@
           @@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@
             @@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@
                              @@@@@@@@@@@@@@@@@@@@@
                               @@@@@@@@@@@@@@@@@@@
                                 @@@@@@@@@@@@@@@
                                        @
    """

    for char in logo:
        print(char, end="", flush=True)
        sleep_time = 0.05
        if speedup:
            time.sleep(sleep_time * 0.1)
        else:
            time.sleep(sleep_time)

def main():
    # Get the job inputs from the environment variable
    job_inputs, job_uuid, flow_uuid, checkpoint_compatible = get_plex_job_inputs()
    print("Job Inputs:", job_inputs)
    if not job_uuid:
        raise ValueError("JOB_UUID environment variable is missing.")
    if not flow_uuid:
        raise ValueError("FLOW_UUID environment variable is missing.")

    os.makedirs("/outputs", exist_ok=True)
    display_logo(job_inputs["speedup"])
    if(checkpoint_compatible.lower == "true"):
        bucket_name = "app-checkpoint-bucket"

        # Simulate checkpoint creation and upload to S3
        for checkpoint in range(0, 3): 
            time.sleep(10)
            object_name = f"checkpoints/{flow_uuid}/{job_uuid}/checkpoint_{checkpoint}"
            event_csv_filename, pdb_path = create_event_csv(checkpoint, job_inputs)
            pdb_file_name = os.path.basename(pdb_path)
            upload_to_s3(event_csv_filename, bucket_name, f"{object_name}/{event_csv_filename}")
            upload_to_s3(pdb_path, bucket_name, f"{object_name}/{pdb_file_name}")
            os.remove(event_csv_filename)
            print(f"Checkpoint {checkpoint} event CSV and PDB created and uploaded.")
    else:
        print("Checkpoint creation is not enabled. Moving on to the next step.")
    
    print("job finished")

if __name__ == "__main__":
    main()