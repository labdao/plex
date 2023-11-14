import sys
import os
import glob

def main(nonfileparamstring, nonfileparamint):
    # Print the non-file parameters
    print(nonfileparamstring)
    print(nonfileparamint)

    # Find the text file in the specified directory
    file_list = glob.glob('/inputs/fileparam/*.txt')
    if file_list:
        file_path = file_list[0]  # There will only be one file in this directory
        with open(file_path, 'r') as file:
            file_contents = file.read()
            print(file_contents)  # Print the contents of the file

        # Write the results to a new file
        with open('/outputs/results.txt', 'w') as result_file:
            result_file.write(nonfileparamstring + '\n')
            result_file.write(str(nonfileparamint) + '\n')
            result_file.write(file_contents)
    else:
        print("No file found in the directory /inputs/fileparam/")

if __name__ == "__main__":
    # Expects two arguments: one string and one integer
    if len(sys.argv) != 3:
        print("This script expects two arguments: a string and an integer.")
    else:
        nonfileparamstring = sys.argv[1]
        try:
            nonfileparamint = int(sys.argv[2])
            main(nonfileparamstring, nonfileparamint)
        except ValueError:
            print("The second argument must be an integer.")
