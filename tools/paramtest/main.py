import argparse
import ast
import glob
import json
import os
import sys



# Function to append text to the result file
def append_to_result(text):
    with open('/outputs/result.txt', 'a') as f:
        f.write(text + '\n')  # Adding newline for readability

# Function to read content from a single file and append to result
def process_single_file(file_path_pattern):
    files = glob.glob(file_path_pattern)
    if files:  # Check if the file exists
        with open(files[0], 'r') as file:
            append_to_result(file.read())

# Function to handle the numeric operations
def multiply_and_append(number, numbers_array):
    result = number * sum(numbers_array)
    append_to_result(str(result))

def bool_and_append(single_bool, bools_array):
    for bool_value in bools_array:
        append_to_result(str(single_bool and bool_value))

def main():
    args = sys.argv[1:]  # Skip the script name
    if len(args) != 6:
        print("Usage: main.py first_string 'string_array' first_number 'number_array' first_bool 'bool_array'")
        sys.exit(1)

    # Extract arguments knowing their fixed positions
    first_string = args[0]
    string_array = json.loads(args[1])
    first_number = int(args[2])
    number_array = json.loads(args[3])
    first_bool = args[4].lower() == 'true'  # Simple string to boolean conversion
    bool_array = json.loads(args[5])

    # Debug print to confirm argument values
    print(f"First String: {first_string}")
    print(f"String Array: {string_array}")
    print(f"First Number: {first_number}")
    print(f"Number Array: {number_array}")
    print(f"First Bool: {first_bool}")
    print(f"Bool Array: {bool_array}")

    if not os.path.exists('/outputs'):
        os.makedirs('/outputs')

    # Clear the result.txt file or create it if it doesn't exist
    open('/outputs/result.txt', 'w').close()

    # Process the first file
    process_single_file('/inputs/first_file/*.txt')

    # Process the files in file_array directories
    for i in range(100):  # assuming 'n' goes from 0 to 99 for the example
        process_single_file(f'/inputs/file_array/{i}/*.txt')

    # Append the input string
    append_to_result(first_string)

    # Append all strings from the array
    for s in string_array:
        append_to_result(s)

    # Perform the numeric operation and append the result
    multiply_and_append(first_number, number_array)

    # Perform the boolean 'and' operation and append the results
    bool_and_append(first_bool, bool_array)

if __name__ == "__main__":
    print("hi")
    print("Raw command-line arguments:", sys.argv)
    main()

