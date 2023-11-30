import csv
import json
import os
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


def main():
    # Get the job inputs from the environment variable
    try:
        job_inputs = get_plex_job_inputs()
        print("Job Inputs:", job_inputs)
    except ValueError as e:
        print(e)
        sys.exit(1)

    input_csvs = job_inputs['input_csvs']
    all_rows = []
    header = None

    for i, csv_file in enumerate(input_csvs):
        with open(csv_file, 'r') as file:
            csv_reader = csv.reader(file)
            current_header = next(csv_reader)

            if i == 0:
                # Store the header from the first file
                header = current_header
                all_rows.extend(csv_reader)  # Add all rows from the first file
            else:
                # Check if the current header matches the first file's header
                if current_header != header:
                    print("Error: Column headers do not match.")
                    sys.exit(1)

                all_rows.extend(csv_reader)  # Add all rows from subsequent files

    # Write to a single CSV file
    with open('/outputs/aggregated_results.csv', 'w', newline='') as file:
        csv_writer = csv.writer(file)
        csv_writer.writerow(header)  # Write the header
        csv_writer.writerows(all_rows)  # Write all combined rows

    print("Aggregated results saved to /outputs/aggregated_results.csv")


if __name__ == "__main__":
    main()
