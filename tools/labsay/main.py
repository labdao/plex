import os
import json
import sys


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


def display_logo():
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
        if char == " ":
            time.sleep(0.0005)


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

        display_logo()
        print("job finished")
        print("results are saved in result.txt")
        print("\U0001F331")


if __name__ == "__main__":
    main()
