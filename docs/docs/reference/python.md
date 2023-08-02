---
title: Python API Reference
description: Python API Reference
sidebar_label: Python
sidebar_position: 1
---

The following functions are considered core when using the plex [pip package](https://pypi.org/project/PlexLabExchange/). Additional documentation for helper functions will be added over time.

**[Install plex](/quickstart/installation)**

---

## plex_init

```python
def plex_init(
    tool_path: str, 
    scattering_method=ScatteringMethod.DOT_PRODUCT.value, 
    plex_path="plex", 
    **kwargs
)
```

Initializes an input/output (IO) JSON from the tool configuration and inputs. The IO JSON defines the execution of tasks within plex.

**Example usage:**

```python
plex_init(
    "/path/to/tool/config", 
    ScatteringMethod.DOT_PRODUCT.value, 
    plex_path="plex", 
    param1=value1, 
    param2=value2
)
```

**Arguments:**
* `tool_path` *str*, **required** - path of tool config; can be a local path or IPFS CID
* `scattering_method` *Enum*, *optional* - method for handling multiple input vectors; supports 'dotProduct' and 'crossProduct'
* `plex_path` *str*, *optional* - path pointing to plex binary
* `**kwargs` *keyword arguments*, *optional* - additional parameters in the form of a dictionary, where keys are input names and values are input values

---

## plex_run

```python
def plex_run(
    io_json_cid: str, 
    output_dir="", 
    verbose=False, 
    show_animation=False, 
    concurrency="1", 
    annotations=[], 
    plex_path="plex"
)
```

Runs a job given instructions defined by an IO JSON.

**Example usage:**

```python
io_json_cid, io_json_local_filepath = plex_run(
    io_json_cid=io_json_cid, 
    output_dir="/path/to/output/dir", 
    verbose=True, 
    show_animation=True, 
    concurrency="2", 
    annotations=["annotation1", "annotation2"],
)
```

**Arguments:**
* `io_json_cid` *str*, **required** - input/output JSON CID
* `output_dir` *str*, *optional* - output directory
* `verbose` *bool*, *optional* - verbose mode with more detailed logs
* `show_animation` *bool*, *optional* - emote animation during job runs
* `concurrency` *str*, *optional* - concurrency for processing jobs
* `annotations` *List[str]*, *optional* - list of annotations for jobs
* `plex_path` *str*, *optional* - path pointing to plex binary

---

## plex_mint

```python
def plex_mint(io_json_cid: str, image_cid="", plex_path="plex")
```

Mints a ProofOfScience token on the Optimism Goerli testnet.

**Example usage:**

```python
plex_mint(completed_io_json_cid)
```

**Arguments:**
* `io_json_cid` *str*, **required** - input/output JSON CID
* `image_cid` *str*, *optional* - custom image CID; when not provided, a [default GIF](https://ipfs.io/ipfs/bafybeiba666bzbff5vu6rayvp5st2tk7tdltqnwjppzyvpljcycfhshdhq) is used for the image metadata
* `plex_path` *str*, *optional* - path pointing to plex binary

:::note

To mint an NFT, the following values **must** be provided as environment variables
* **RECIPIENT_WALLET** - a 0x wallet address on [Optimism Goerli testnet](https://goerli-optimism.etherscan.io/)
* **AUTOTASK_WEBHOOK** - necessary for gasless minting

```python
os.environ["RECIPIENT_WALLET"] = "" # enter your wallet address
os.environ["AUTOTASK_WEBHOOK"] = "https://api.defender.openzeppelin.com/autotasks/e15b3f39-28f8-4d30-9bf3-5d569bdf2e78/runs/webhook/8315d17c-c493-4d04-a257-79209f95bb64/2gmqi9SRRAQMoy1SRdktai" 
```

:::