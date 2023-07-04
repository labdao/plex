# Colabfold Tool for PLEX

The `colabfold` tool directory contains all the necessary files to run protein folding predictions using Colabfold within the PLEX environment.

## Requirements
The PLEX environment should be correctly installed and configured as per the main README instructions. The following resources are needed to run Colabfold:

- Docker: Pulls the `public.ecr.aws/p7l9w5o7/colabfold:latest` Docker image to create the container environment.
- GPU: The tool requires GPU resources.
- Network: Internet connection is needed.

## Dockerfile

The `Dockerfile` sets up a Docker image for running the Colabfold protein folding tool. It uses a base image, installs necessary Linux packages, checks for installed tools, and sets up a directory for storing parameters.

## Colabfold Configuration JSON

The `_colabfoldls.json` file is a CWL (Common Workflow Language) specification for running Colabfold. It specifies the required Docker image, hardware requirements, input parameters, and the output.

The JSON file contains:
- `baseCommand`: Runs the command `ls -lah /inputs > /outputs/inputs.txt` to list the files in the inputs directory.
- `inputs`: Accepts a FASTA file and a recycle value as input.
- `outputs`: Outputs a file `inputs.txt`.

## Test FASTA File

The `test.fasta` file is a test protein sequence file in FASTA format, used to ensure the tool is working as expected.

## Running Colabfold

To run the Colabfold tool, navigate to the main directory of the PLEX tool and use the following command:

```
./plex create -t tools/colabfold/_colabfoldls.json -i testdata/yourfile.fasta --autoRun=True
```
Replace `yourfile.fasta` with the path to your own FASTA file.

## Output

The output of the tool will be a file `inputs.txt` listing the input files.

## Contributing

We welcome contributions to improve this tool. Please see the main repository's contribution guidelines for more information.

## Issues

If you encounter any issues while using this tool, please open an issue in the main repository. We appreciate your feedback.