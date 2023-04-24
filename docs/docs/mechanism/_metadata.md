# Why does metadata matter?
At the heart of LabDAO is the idea that a community operated marketplace protocol for laboratory services could generate a [dynamic knowledge graph](https://niklasrindtorff.substack.com/p/building-a-knowledge-graph-for-biological) for biomedicine. The receipts of past transactions among scientists on the lab-exchange serve as its building blocks.

To create the marketplace protocol and the subsequent knowledge graph, we need to work out a structured way to describe scientific services and the resulting data. We need standards for composable metadata.

## Scientific data should be FAIR
Scientific data is ideally stored according to [FAIR](https://www.go-fair.org/fair-principles/) principles, developed by [Mark Wilkinson](https://www.nature.com/articles/sdata201618). These include: 
* **F**indability
* **A**ccessibility
* **I**nteroperability
* **R**euse 

Put differently, data is ideally stored at a location that is known to others, accessible to others, adheres to public standards and thus invites reuse. 

Web3 introduces new methods to share, control and distribute data according to FAIR principles. These methods include:
* direct content addressing 
* on-chain provenance tracking
* robust access control

With these methods, data can be pinned visibly and addressable for everyone. Ownership of data can be traced and access to encrypted information can be tied to token ownership. Web3 increases the findability and accessibility of data - the "F" and "A" of FAIR.

To make data more interoperable ("I") and reusable ("R") its formatting needs to be standardized and it needs to be presented with sufficient metadata to provide relevant context. LabDAO is an online community of scientists and engineers that, among many things, adopts and defines standards for data formats and metadata structure. This article introduces a proposed simple standard for tokenized scientific data. 

## Web3 x Bio forms of metadata 
Let's take a look at common metadata formats seen in both web3 and bio. 

### NFTs
A common case where metadata in web3 is important, is for NFTs following the ERC721 standard. An object, often a digital artwork, is referenced within a metadata JSON together with additional properties. Often these properties give us more details about the artwork. 

![You can see the artwork with its properties below](https://hackmd.io/_uploads/BylIvxsrc.png)

This is how such a metadata JSON looks like. The picture of the penguin itself is referenced under *image* using a universal resource identifier (URI). 

```` pudgy.json
{
  "attributes": [
    {
      "trait_type": "Background",
      "value": "Blue"
    },
    {
      "trait_type": "Skin",
      "value": "Normal"
    },
    {
      "trait_type": "Body",
      "value": "Lab Coat"
    },
    {
      "trait_type": "Face",
      "value": "Cross Eyed"
    },
    {
      "trait_type": "Head",
      "value": "Cowboy Hat"
    }
  ],
  "description": "A collection 8888 Cute Chubby Pudgy Penguins sliding around on the freezing ETH blockchain.",
  "image": "https://ipfs.io/ipfs/QmNf1UsmdGaMbpatQ6toXSkzDpizaGmC9zfunCyoz1enD5/penguin/3020.png",
  "name": "Pudgy Penguin #3020"
}
````

This existing NFT standard as also been adapted by our friends at [Molecule](https://www.molecule.to/) to represent legal sublicensese for intellectual property. This form of NFT is called an IP-NFTs. 

### Biocompute Objects
In bio, specifically computational biology, an international standard for metadata exists: [biocompute objects](https://www.biocomputeobject.org/). If this is absolutely new to you, you are in great company. A lot of biologists, especially in academia, do not know this standard exists. 

To adhere to this standard (which rarely happens), scientists have to provide extensive information about multiple aspects: 
* who created the data (provenance domain)
* a description of the data (usability domain)
* a step-by-step computational pipeline of how the data was processed (description domain)
* information about the environment in which the processing was run (execution domain)
* links to input and output data (io domain)
* parameters that were used in the computational pipeline (parameter domain)

Biocompute objects like this [JSON](https://biocomputeobject.org/objects/view/https/biocomputeobject.org/BCO_014961/1.0) can easily take up more than 700 lines. They lack any web3 primitives, such as decentralized identifiers and stable content addressing via IPFS for referenced data. Could there be a simplified structure for biomedical metadata that makes use of all the tools decentralized storage and records give us?


## lab-NFTs
By building on the basic structure of the biocompute object and simplifying it with web3 primitives, we have come up with a proposed standard for lab-NFT metadata. To reduce the complexity of any given metadata object and allow for reusability of components -a prerequisite for a knowledge graph- we factorize the metadata into 7 components:
1. object information is at top-level and includes a description (useability domain)
2. user: user object
3. parameters: parameter object
4. input: data object that includes context and link to a file
5. provider: provider object
6. execution: execution object (contains information about the runtime and the computational graph that was executed. The execution is often pointing to a standardized docker container)
7. output: data object that includes context and link to a file

Components 2-7 are all pinned as separate JSON objects on IPFS and referenced here. Their URI can be reused in other transactions, leading to new lab-NFTs, and a branched tree of references - our knowledge graph.

A very simple example for metadata can be seen below. In this example a user is requesting the marketplace to generate the reverse complement of a set of DNA sequences. 

To initiate the transaction, the user is generating an incomplete lab-NFT metadata object, which we refer to as a *job object*. A *job object* includes the top-level descriptors, including the name of the service. In addition the *job object* includes a reference to the 1) user, 2) parameters and 3) input data. 

Once the transaction request is posted and funds have been deposited in the exchange escrow, providers in the network can claim the job. After completing the service, the provider adds references to 4) their identity, 5) information about the job execution and 6) output data. 

Simply put, a user of the protocol is creating a half-baked NFT metadata object and is paying someone else (the provider) to complete the metadata and mint the token.

```` toplevel.json
{
  "name": "reverse-complement", # a standardized name for the performed process
  "description": "the reverse complement of a DNA sequence",
  "image": "ipfs://QmZ9oReVUiNQSc9GaqqTEPUW3XHo6eprVSa9nqbGNotP8B",
  "properties": {
      "user": "0x64BC15E0A5A12dDbe321EEDD832d057775D11F56"
      "parameters": # no parameters in this simple example
      "input": "ipfs://QmaQv91AUmu4wDNiNc8exUrSrDnZ65YBEXuTrxSJSnJbcz",
      "provider": 
      "execution": 
      "output":
  }
}
````

The **user and provider object** represent the identity of the people involved in the transaction. The identity can be a referenced JSON object or a [holonym](https://pulse.opsci.io/provable-and-computable-identity-for-future-proof-scientific-workflows-b020cdea11e3).

A **parameter object** is defined by the LabDAO community and is application specific. Certain simple applications, such as our reverse complement example, do not require parameters.

An example for parameters of a more complex job, an alphafold v2 run, is shown below: 

```` parameter.json
{
  "weights": # URI to model weights
  "mode": "monomer",
  "database": "full",
  "max_template_date": "2022-01-01",
  "is_prokaryote": false
}
````

Returning to our reverse complement example, the **input object** is following a simple structure and is referencing a [fasta file](https://ipfs.io/ipfs/QmfURWZakhnD1Rn3DqC38Eqps1zgZ17dxM4KgV1rmJDhww) with the DNA sequence of interest. The output object added by the provider after completion of the job follows the same structure.

```` input.json
{
  "name": "example_sequence.fasta"
  "type": "fasta"
  "uri": "ipfs://QmfURWZakhnD1Rn3DqC38Eqps1zgZ17dxM4KgV1rmJDhww"
}
````

The **execution object** contains information about the environment in which a job was processed. In our case it is referencing a container image used for the completion of the job.

```` execution.json
{
  "image" # URI of reverse complement application container
  "runtime": "deta microservice, https://02wun6.deta.dev/"
  "logs": # additional information
}
````


## Endgame: encryption and composable NFTs
The lab-NFT concept allows us to exchange FAIR scientific data on the lab-exchange. Scientists can easily work together to perform research and share their results openly. 

To accelerate not only progress in the basic sciences, but also innovation around valuable intellectual property, we will introduce threshold encryption for lab-NFTs. All input and output data is symmetrically encrypted before it is pinned to IPFS. The owner of the lab-NFT can request access to the referenced files using tools like [lit protocol](https://litprotocol.com/).

Most scientific publications are based on data generated through multiple measurements. The same is true for most forms of intellectual property. To enable users to bundle individual lab-NFTs for simplified distribution, composable IP-NFTs will be developed based on existing implementations, such as [aavegotchi](https://github.com/aavegotchi/aavegotchi-contracts) which is inspired by the [EIP998](http://erc998.org/) proposal. Composability will allow authors to bundle their work and distribute it more easily. 
