# PLEx Lab Exchange

## Installing the client

First, install the client by running

```
source <(curl -sSL https://raw.githubusercontent.com/labdao/ganglia/main/plex/install.sh)
```

The installer may ask for your password at some point. 

When the installer is complete, next set your web3.storage API token.

```
export WEB3STORAGE_TOKEN=<your token here>
```

Finally, update the executable's permissions.

```
chmod +x ./plex
```

## Running the client

Once the client is installed, you can run the following command in the newly-created `plex` folder to run equibind.

```
./plex -app equibind -input-dir ./testdata -gpu false
```

## Running a node
This is a script for setting up a compute instance to run LabDAO jobs. Requires linux OS with Nvidia GPU.

Tested on Ubuntu 20.04 LTS with Nvidia T4, V100, and A10 GPUs (AWS G4, P3, and G5 instance types)

The install script sets up Docker, Nvidia Drivers, Nvidia Container Toolkit, and IPFS
```
curl -sL https://raw.githubusercontent.com/labdao/ganglia/main/sripts/provide_compute.sh | bash && newgrp docker
```

After the script run the following command in a separate terminal to start a Bacalhau server to accept jobs.
```
ipfs daemon
```

Once the daemon is running, configure the Bacalhau node based on the addresses used by the IFPS node.
```
ipfs id

# copy the ip4 tcp output and change port 4001 to 5001 then export
export IPFS_CONNECT=/ip4/127.0.0.1/tcp/5001/p2p/<your id goes here>

# example: export IPFS_CONNECT=/ip4/127.0.0.1/tcp/5001/p2p/12D3KooWPH1BpPfNXwkf778GMP2H5z7pwjKVQFnA5NS3DngU7pxG

LOG_LEVEL=debug bacalhau serve --job-selection-accept-networked --limit-total-gpu 1 --limit-total-memory 12gb --ipfs-connect $IPFS_CONNECT
```

## Notes
To download large bacalhau results the below command may need ran 
```
sudo sysctl -w net.core.rmem_max=2500000
```
