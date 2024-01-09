# process_fasta.py
import sys

input_file = sys.argv[1]  # Take the input file path as a command-line argument

with open(input_file, 'r') as file:
    lines = file.readlines()

# Replace '>1' with '>binder' in the first line
lines[0] = lines[0].replace('>1', '>omegafold_binder')

# Remove everything up to and including ':' in all lines
lines_to_write = [line[line.find(':')+1:] if ':' in line else line for line in lines]

with open(input_file, 'w') as file:
    file.writelines(lines_to_write)


### OLD CODE ###

# # process_fasta.py
# import sys

# input_file = sys.argv[1]  # Take the input file path as a command-line argument

# with open(input_file, 'r') as file:
#     lines = file.readlines()

# new_lines = []

# for line in lines:
#     # Check if the line contains a colon
#     if ':' in line:
#         # Split the line at the colon
#         before_colon, after_colon = line.split(':', 1)
#         # Add the part before the colon to new_lines
#         new_lines.append(before_colon + '\n')
#         # Add a new chain header for the part after the colon
#         new_lines.append('>' + '\n')
#         # Add the part after the colon as a new line
#         new_lines.append(after_colon)
#     else:
#         # If there is no colon, just add the line as it is
#         new_lines.append(line)

# # Replace '>1' with '>' in the first line if present
# new_lines[0] = new_lines[0].replace('>1', '>')

# with open(input_file, 'w') as file:
#     file.writelines(new_lines)

# # process_fasta.py
# import sys

# input_file = sys.argv[1]  # Take the input file path as a command-line argument

# with open(input_file, 'r') as file:
#     lines = file.readlines()

# # Replace '>1' with '>' in the first line and ':' with '' in all lines
# lines[0] = lines[0].replace('>1', '>')
# lines_to_write = [line.replace(':', '/') for line in lines]

# with open(input_file, 'w') as file:
#     file.writelines(lines_to_write)