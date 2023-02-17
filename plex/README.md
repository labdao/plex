# Plex

## Installing the client

First, install the client by running

```
source <(curl -sSL https://raw.githubusercontent.com/labdao/ganglia/main/plex/install.sh)
```

The installer may ask for your password at some point. Next, set your web3.storage API token.

```
export WEB3STORAGE_TOKEN=<your token here>
```

## Running the client

Once the client is installed, you can run the following command in the newly-created `plex` folder to run equibind.

```
./plex -app equibind -input-dir ./testdata/pdbbind_processed_size1
```
