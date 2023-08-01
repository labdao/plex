---
title: plex_create
description: Reference for plex_create
sidebar_label: create
sidebar_position: 2
---


### Purpose

The `plex_create` function creates and pins an IO JSON file, which contains a detailed representation of the inputs, outputs, and operations of a specific tool execution within the Plex application. This Python function serves as a wrapper around the corresponding Golang CLI command `plex create`.


### Syntax

```python
def plex_create(tool_path: str, input_dir: str, layers=2, output_dir="", verbose=False, show_animation=False, concurrency="1", annotations=[], plex_path="plex")
```

### Arguments

| Parameter       | Type    | Description   | Default     | Required |
|-----------------|---------|---------------|-------------|----------|
| `tool_path`     | str     | Path to the tool JSON file. | - | Yes |
| `input_dir`     | str     | Directory containing input files. | - | Yes |
| `layers`        | int     | Number of layers to search input directory. | 2 | No |
| `output_dir`    | str     | Output directory. If not specified, the current directory will be used. | "" | No |
| `verbose`       | bool    | If set to True, enables verbose output. | False | No |
| `show_animation`| bool    | If set to True, shows job processing animation. | False | No |
| `concurrency`   | str     | Number of concurrent operations. | "1" | No |
| `annotations`   | list    | Annotations to add to the job. | [] | No |
| `plex_path`     | str     | Path to the Plex executable. | "plex" | No |

### Return value

The function returns the content identifier (CID) of the created and pinned IO JSON file.

### Example usage

```
plex_create(
    tool_path="path/to/tool.json",
    input_dir="path/to/inputs",
    layers=3,
    output_dir="path/to/output",
    verbose=True,
    show_animation=False,
    concurrency="4",
    annotations=["annotation1", "annotation2"],
    plex_path="/usr/local/bin/plex"
)
```