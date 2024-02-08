import os
import json
import sys
import boto3

def upload_to_s3(file_name, bucket_name, object_name=None):
    print("file_name: ", file_name)
    print("bucket_name: ", bucket_name)
    print("object_name: ", object_name)

    if object_name is None:
        object_name = file_name

    s3_client = boto3.client('s3')
    try:
        response = s3_client.upload_file(file_name, bucket_name, object_name)
        print(f"Successfully uploaded {file_name} to {bucket_name}/{object_name}")
    except Exception as e:
        print(f"Failed to upload {file_name} to {bucket_name}/{object_name}: {e}")
        raise e

def get_plex_job_inputs():
    # Retrieve the environment variable
    json_str = os.getenv("PLEX_JOB_INPUTS")

    # Check if the environment variable is set
    if json_str is None:
        raise ValueError("PLEX_JOB_INPUTS environment variable is missing.")

    # Convert the JSON string to a Python dictionary
    try:
        data = json.loads(json_str)
        return data
    except json.JSONDecodeError:
        # Handle the case where the string is not valid JSON
        raise ValueError("PLEX_JOB_INPUTS is not a valid JSON string.")


import time


def main():
    # Get the job inputs from the environment variable
    try:
        job_inputs = get_plex_job_inputs()
        print("Job Inputs:", job_inputs)
    except ValueError as e:
        print(e)
        sys.exit(1)

    # Create /outputs directory if it doesn't exist
    os.makedirs("/outputs", exist_ok=True)

    # Create or overwrite the result.txt file
    with open("/outputs/result.txt", "w") as result_file:
        # Read contents of file_example and write to result.txt
        try:
            with open(job_inputs["file_example"], "r") as file_example:
                result_file.write(file_example.read() + "\n")
        except FileNotFoundError:
            print(f"File {job_inputs['file_example']} not found.")
            sys.exit(1)

        # Append string_example to result.txt
        result_file.write(job_inputs["string_example"] + "\n")

        # Calculate product and append to result.txt
        product = job_inputs["number_example"] * len(job_inputs["string_example"])
        result_file.write(
            f"Product of number_example and length of string_example: {product}\n"
        )

        # display_logo()
        job_uuid = os.getenv("JOB_UUID")
        print("job uuid: ", job_uuid)
        print("job finished")
        print("results are saved in result.txt")
        print("\U0001F331")

    bucket_name = "app-checkpoint-bucket"

    # Simulate checkpoint creation and upload to S3
    for checkpoint in range(1, 4):  # Example: Create 3 checkpoints
        object_name = f"checkpoints/{job_uuid}/checkpoint_{checkpoint}"
        checkpoint_filename = f"checkpoint_{checkpoint}_intermediatefile.txt"
        with open(checkpoint_filename, "w") as f:
            f.write(f"Checkpoint {checkpoint} data\n")
        print(f"Checkpoint {checkpoint} file created.")

        # Upload the checkpoint file to S3
        upload_to_s3(checkpoint_filename, bucket_name, f"{object_name}/{checkpoint_filename}")

        # Simulate some computation time
        time.sleep(1)  # Sleep for 1 second to simulate computation


if __name__ == "__main__":
    main()
