---
title: plex_init
description: Reference for plex_init
sidebar_label: init
sidebar_position: 1
---

### Purpose

The `plex_init` function initializes an IO JSON from the tool configuration and inputs. These parameters guide the execution of tasks within the Plex application. 


### Syntax

```python
plex_init(tool_path: str, scattering_method=ScatteringMethod.DOT_PRODUCT.value, plex_path="plex", **kwargs)
```

### Arguments

| Name               | Type                               | Description                                                               | Default                               | Required |
|--------------------|------------------------------------|---------------------------------------------------------------------------|---------------------------------------|----------|
| `tool_path`        | str                                | The path of the tool configuration. This can be a local path or an IPFS CID. | N/A                                   | Yes      |
| `scattering_method`| ScatteringMethod Enum value (string) | The method for handling multiple input vectors. It supports 'dotProduct' and 'crossProduct'. | `ScatteringMethod.DOT_PRODUCT.value` | No       |
| `plex_path`        | str                                | The path where the `plex` command is located.                            | `'plex'`                              | No       |
| `**kwargs`         | keyword arguments                  | Additional parameters in the form of a dictionary, where keys are input names and values are input values. | N/A                                   | No       |

### Return value

A string containing the IO JSON CID.

### Exceptions

Raises a `PlexError` if the initialization of the IO JSON CID fails.

### Example usage

```python
plex_init("/path/to/tool/config", ScatteringMethod.DOT_PRODUCT.value, plex_path="plex", param1=value1, param2=value2)
```

### Note

The `plex_init` function uses the `subprocess` module to call an external Go command. The standard output and error of this command are printed to the console. The command is expected to print a line containing the string "Pinned IO JSON CID:", followed by the CID. This CID is then returned by the `plex_init` function.