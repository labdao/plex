# Currently hardcoded to work on GPU
# Run functions individually to set up on other instances

```bash -c "source ./provide-compute.sh; setup"```

```export PLEX_ENV=<stage or prod>```

```bash -c "source ./provide-compute.sh; start"```

You may visit the instance at:

```http://<public IPV4 address>:8888/lab```
