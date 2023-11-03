import time
import os
import glob

def read_and_print_file(file_path):
    # Read the file
    with open(file_path, 'r') as file:
        lines = file.readlines()

    # Print each line with a delay
    total_lines = len(lines)
    for i, line in enumerate(lines, start=1):
        print(f"Total Lines: {total_lines}, Current Line: {i}")
        print(line.strip())
        time.sleep(1)

def main():
    # Find the txt file in the specified directory (assumes only one .txt file is present)
    file_pattern = os.path.join('/', 'inputs', 'message', '*.txt')
    files = glob.glob(file_pattern)

    if files:
        # Read and print the first (and only) file found
        read_and_print_file(files[0])
    else:
        print("No .txt file found in the 'inputs/message/' directory.")

if __name__ == "__main__":
    main()


