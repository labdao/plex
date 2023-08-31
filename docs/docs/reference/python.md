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
    auto_run=False,
    **kwargs
)
```

Initializes an input/output (IO) JSON from the tool configuration and inputs. The IO JSON defines the execution of tasks within plex.

**Example usage:**

```python

fasta_filepath = ["home/project/protein.fasta"]

initial_io_cid = plex_init(
    CoreTools.COLABFOLD_MINI.value, 
    sequence=fasta_filepath
)
```

**Arguments:**
* `tool_path` *str*, **required** - path of tool config; can be a local path, CoreTools reference, or IPFS CID
* `scattering_method` *Enum*, *optional* - method for handling multiple input vectors; supports 'dotProduct' and 'crossProduct'; default value is dotProduct
    * **dotProduct:** Pairs corresponding elements from each input vector into subarrays. Requires all input vectors to have the same length. 
    * **crossProduct:** Generates all combinations of elements across input vectors, forming the Cartesian product. Input vector lengths can vary.
* `plex_path` *str*, *optional* - path pointing to plex binary
* `auto_run` *bool*, *optional* - automatically submits the job for computation based on the IO JSON
* `**kwargs` *keyword arguments*, *optional* - additional parameters in the form of a list, where keys are input names and values are input values; see each tool config for specific arguments accepted

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
    io_json_cid=initial_io_cid, 
    output_dir="home/project", 
)
```

**Arguments:**
* `io_json_cid` *str*, **required** - input/output JSON CID
* `output_dir` *str*, *optional* - output directory; plex default creates a jobs directory
* `verbose` *bool*, *optional* - verbose mode with more detailed logs
* `show_animation` *bool*, *optional* - emote animation during job runs
* `concurrency` *str*, *optional* - concurrency for processing jobs
* `annotations` *List[str]*, *optional* - list of annotations for jobs; mostly used for usage metrics
* `plex_path` *str*, *optional* - path pointing to plex binary

---

## plex_mint

```python
def plex_mint(io_json_cid: str, image_cid="", plex_path="plex")
```

Mints a ProofOfScience token on the Optimism Goerli testnet.

**Example usage:**

```python
plex_mint(io_json_cid)
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