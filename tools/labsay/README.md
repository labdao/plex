# labsay-seq-only tool

Labsay is an internal tool designed to merge the input functionality of the seq-only generator with the output of Labsay v0.8, to test the unified view frontend implementation and increase the development iteration speed.

## Features

- Accepts seq-only generator inputs.
- Takes in labsay sample checkpoint pdb inputs.
- Uploads simulated checkpoints to s3.

## Points to note

- This tool is checkpoint compatible, so the tool manifest reflects the same, with the `checkpointCompatible` flag set to `True`. 
- If you would like to test certain functionalities without checkpoints, and require checkpointCompatible = False, please set the above flag to false.
- While the tool mirrors labsay checkpoint inputs and seq-only generator inputs together, it does not take in the traditional file_example, number_example values during the experiment submission stage (like labsay v0.8).

## Steps to test locally:

```bash
cd tools/labsay/
chmod +x test.sh
./test.sh
```
(or)

### Note: Please refer to [CHANGELOG.md](./CHANGELOG.md) for steps to test the latest version of the tool.