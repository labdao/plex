---
title: Tokens
sidebar_position: 5
sidebar_label: Tokens
---

[**Records**](https://medium.com/@labdao/introducing-records-tokens-for-scientific-creation-d938fbf553e4) represent plex's unique approach to preserving, acknowledging, and ensuring the reproducibility of scientific computations. By leveraging the power of blockchain, each computation in plex can be minted as an [ERC-1155](https://ethereum.org/en/developers/docs/standards/tokens/erc-1155/) NFT.

## Minting with `plex_mint`

Once a computation concludes and its results are recorded in a completed `io.json`, the `plex_mint` command can be invoked. This process transforms the results into a tangible, traceable, and verifiable ProofOfScience NFT.

## Metadata Preservation

Within the NFT's metadata, the `graph` key contains the `io.json` content. All completed job runs are visible. All input and output data are accessible.

By providing this level of transparency and detail, others can validate, reproduce, or build upon the work.

```json
{
  "description": "Research, Reimagined. All Scientists Welcome.",
  "graph": [
    {
      "errMsg": "",
      "inputs": {
        "protein": {
          "class": "File",
          "filepath": "7n9g.pdb",
          "ipfs": "QmUWCBTqbRaKkPXQ3M14NkUuM4TEwfhVfrqLNoBB7syyyd"
        },
        "small_molecule": {
          "class": "File",
          "filepath": "ZINC000003986735.sdf",
          "ipfs": "QmV6qVzdQLNM6SyEDB3rJ5R5BYJsQwQTn1fjmPzvCCkCYz"
        }
      },
      "outputs": {
        "best_docked_small_molecule": {
          "class": "File",
          "filepath": "7n9g_ZINC000003986735_docked.sdf",
          "ipfs": "QmZdoaKEGtESnLoHFMb9bvqdwXjyUuRK6DbEoYz8PYpZ8W"
        },
        "protein": {
          "class": "File",
          "filepath": "7n9g.pdb",
          "ipfs": "QmUWCBTqbRaKkPXQ3M14NkUuM4TEwfhVfrqLNoBB7syyyd"
        }
      },
      "state": "completed",
      "tool": {
        "ipfs": "QmZ2HarAgwZGjc3LBx9mWNwAQkPWiHMignqKup1ckp8NhB",
        "name": "equibind"
      }
    }
  ],
  "image": "ipfs://bafybeiba666bzbff5vu6rayvp5st2tk7tdltqnwjppzyvpljcycfhshdhq",
  "name": "yielding hubble proteins"
}
```

## Reproducibility and Acknowledgement

Storing computations as Records on-chain sets a gold standard for scientific reproducibility. It becomes an immutable record of achievement, open to scrutiny and validation by peers.

## Gasless Transactions

Plex employs an [OpenZeppelin Defender Relayer](https://docs.openzeppelin.com/defender/relay) so users don't have to pay [gas fees](https://ethereum.org/en/developers/docs/gas/) to mint ProofOfScience tokens.

:::warning

Please only interact with the official [**smart contract**](https://goerli-optimism.etherscan.io/address/0xda70C0709d4213eE8441E4731A5F662C0406ed7e#code). The only blockchain we are on is the Optimism Goerli testnet. We are **NOT** on mainnet.

**Official address:** 0xda70C0709d4213eE8441E4731A5F662C0406ed7e

:::