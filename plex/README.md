# Plex

## Installing the client

First, install the client by running

```
source <(curl -sSL https://raw.githubusercontent.com/labdao/ganglia/main/plex/install.sh)
```

The installer may ask for your password at some point. You will also be prompted to enter a web3.storage API token.

## Running the client

Once the client is installed, you can run the following command in the `plex` folder to run equibind.

```
./plex -app equibind -input-dir ./testdata/pdbbind_processed_size1
```
