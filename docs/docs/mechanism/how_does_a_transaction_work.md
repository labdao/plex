# How does a transaction on the lab-exchange work? 
A transaction on the lab-exchange is a multi-step process. To break down the required steps, we illustrate an application call that includes the basic *lab-reverse-complement* service. 

**TODO - feature in development!**

## phase 1: transaction definition on the [client] side
1. the client has installed the [openlab CLI] and created a [local wallet]
2. the client lists the available [applications] using ``` openlab app list ```. In our example the service we focus on ```reverse-complement```
3. the client submits a job using ``` openlab job submit reverse-complement```. During this step the command line is asking for input information interactively. In case the user prefers to use the tool without an interactive component, a JSON template can be exported, edited and submitted with the commands below:
```
# non-interactive submission of jobs on openlab:
# export the job instruction JSON object to a example_directory
openlab app example reverse-complement /example_directory
# edit the job instruction object locally
nano /target_directory/job-reverse_complement.json
# submit a prepared job instruction instead of an interactive process
openlab job submit -t /target_directory/job-reverse_complement.json
```


## phase 2: transaction submission on the [client] side
underneath the hood of ```openlab job submit``` multiple functions are called: 
1. the job object, a JSON with the instructions, is generated. In our case ```job-reverse_complement_20220408184322.json``` is stored in a local directory controlled by the client. An example for ```job-reverse_complement_20220408184322.json``` is below: 

````
{
    "service": 0,
    "service_name": "lab-reverse_complement",
    "input":{
        "sequence": "ctatataaataaataaataaatattatatatatag"
    }
}
````

2. the job object is pinned on IPFS, giving us a ```jobURI``` which will be needed when interacting with the smart contracts of the openlab exchange. The function call is ``` openlab file push job-reverse_complement_20220408184322.json```. 
3. the client interacts with the lab-exchange [exchange contract](https://mumbai.polygonscan.com/address/0xfcF2b192c888d411827fDa1884C6FE2438C15Ad0#writeContract) and calls the ```submitJob``` function. The ```jobURI``` of the job object is an argument of this function. The job is created and enters the ```open``` state.
4. optional: to facilitate the execution of services, the client can also share the job object directly with an [index service] maintained by the DAO via http. The service broadcasts all requests from clients to a collection of community-registered servers. 

![openlab_state](https://github.com/labdao/assets/blob/main/openlab_exchange/state_transition.png?raw=true)

## phase 3: transaction verifcation on the [provider] side
1. the server checks the state of the openlab exchange contract by querying the [subgraph](https://thegraph.com/hosted-service/subgraph/tohrnii/openlab-exchange-mumbai-c) for jobs within the ```open``` state.
2. the server pulls job objects via their ``` jobURI``` from IPFS and checks wether the requested service is within the repertoire of services that can be provided.
3. optional: to facilitate the execution of services, the server can receive job objects from the [index service] and filters incoming requests before verifying them (step 1) and their job objects (step 2).
4. the server interacts with the openlab [exchange contract](https://mumbai.polygonscan.com/address/0xfcF2b192c888d411827fDa1884C6FE2438C15Ad0#writeContract) and calls the ```acceptJob``` function, with the transaction's ```_jobId``` being the only input argument.
5. the server validates that the ```acceptJob``` function call was successful. 

![](https://github.com/labdao/assets/blob/main/openlab_exchange/Group%203.png?raw=true)

## phase 4: transaction processing on the [provider] site
once the transaction is verified and claimed the server can get to work and process the request.
1. the server accesses required information from job object, including the pulling of input data from IPFS
2. the server takes the input information and starts the service. In our case the service is accessible via https://02wun6.deta.dev/ and the request from the server would look like ``` curl https://02wun6.deta.dev/ctatataaataaataaataaatattatatatatag```. In a live deployment, where running the service was not for free, the server would protect the endpoint. 
3. the server collects all output data generated during job processing, and pins files to IPFS together with a metadata JSON object reffered to as the token object. The token object is referenced with a ```tokenURI``` which is an input argument for the ```closeJob``` function. An example token object is displayed below: 

````
{
    "service": 0,
    "service_name": "lab-reverse_complement",
    "input":{
        "sequence": "ctatataaataaataaataaatattatatatatag"
    },
    "output":{
        "sequence": "ctatatatataatatttatttatttatttatatag"
    },
    "metadata":{
        "repository": "https://github.com/openlab-apps/lab-reverse_complement/tree/2-serverless-stage",
        "commit": "ddf80e0e953fde3d8cbe9f17fc3f1108cac7bc37",
        "deployment": "https://02wun6.deta.dev/"
    },
    "job_uri": "ipfs/QmdQP9D44Hgp8697FpGJrkTiQYSiUH3xEsvumKX3jJmV58",
    "image": "https://gateway.pinata.cloud/ipfs/QmZ9oReVUiNQSc9GaqqTEPUW3XHo6eprVSa9nqbGNotP8B"
}
````
3. the server calls the ```swap```function, which includes both the minting of a LAB-NFT, the transfer of that LAB-NFT to the client, the calling of ```closeJob```, and the claiming of funds held within the escrow. 

## phase 5: transaction end on the [client] side
In the near future, clients will be able to flag a transaction between receiving the LAB-NFT and release of funds from escrow.
