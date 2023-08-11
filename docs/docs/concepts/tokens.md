---
title: Tokens
sidebar_position: 5
sidebar_label: Tokens
---

**ProofOfScience** represents Plex's unique approach to preserving, acknowledging, and ensuring the reproducibility of scientific computations. By leveraging the power of blockchain, each computation in Plex can be minted into an ERC 1155 token called a ProofOfScience Non-Fungible Token (NFT).

## Minting with `plex_mint`

Once a computation concludes and its results are recorded in a completed `io.json`, the `plex_mint` command can be invoked. This process transforms the results into a tangible, traceable, and verifiable ProofOfScience NFT.

## Metadata Preservation

The NFT's embedded metadata captures essential details about the completed job runs, ensuring a comprehensive overview of the scientific process undertaken. By providing this level of transparency and detail, it enables others in the scientific community to validate, reproduce, or build upon the work.

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

Storing computations as ProofOfScience tokens on-chain not only serves as a testament to the work done but also sets a gold standard for scientific reproducibility. It becomes an immutable record of achievement, open to scrutiny and validation by peers.

## Gasless Transactions

Recognizing the potential hurdles of transaction costs, Plex employs an OpenZeppelin Defender Relayer. This means users don't bear the brunt of gas fees when minting their ProofOfScience NFTs. We take care of it, making the minting process smooth and cost-efficient.

:::warning

Please only interact with the official [**smart contract**](https://goerli-optimism.etherscan.io/address/0xda70C0709d4213eE8441E4731A5F662C0406ed7e#code). The only blockchain we are on is the Optimism Goerli testnet. We are **NOT** on mainnet.

**Official address:** 0xda70C0709d4213eE8441E4731A5F662C0406ed7e

:::