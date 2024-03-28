# Changelog

All notable changes to this tool will be documented in this file, in the order of latest to oldest.

# Versions available on the platform:

## fastcolabdesign v0.8 - 2024-03-28

- Same as colabdesign v0.8, but a lightweight version with the ability to input the number of binders
- This version is available under community models, for internal purpose only
- The code (main.py) still has the logic for colabdesign v0.8. the change to have number of binders to 240 was changed briefly just to get an image, and then reverted back, so future updates can be made on top of the main model colabdesign v0.8

## colabdesign v0.8 - 2024-03-25

- Latest colabdesign model with checkpoints
- Default number of binders set to 240
- Available under protein-binder-design

# Archived older versions:

## colabdesign v0.7

- This model has job level checkpoints.
- Archived after checkpoints are changed to flow level. Newer models upload checkpoints with this folder structure: {flowUUID}/{jobUUID}/...