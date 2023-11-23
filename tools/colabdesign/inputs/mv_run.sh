#!/bin/bash

# Remove the existing main.py file
rm -f main.py

# Copy main.py from ../inputs/code/ to the current directory
cp ../inputs/code/main.py .

# Run main.py with Python 3
python3 main.py

# chmod +x mv_run.sh
# ./mv_run.sh