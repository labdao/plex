---
title: Content-addressed Data
sidebar_position: 2
sidebar_label: Data
---

Plex utilizes a decentralized storage system, [IPFS](https://docs.ipfs.tech/), for managing file storage in its scientific computing workflows. Within this system, each file is content-addressed, meaning it's given a unique content identifier ([CID](https://docs.ipfs.tech/concepts/content-addressing/#what-is-a-cid)) derived from its actual content rather than its location or name. 

This ensures that even if a file is moved, renamed, or distributed across the network, it can be easily and precisely located using its CID. This method not only enhances file retrieval but also promotes data integrity since the identifier changes if the content does, making any alterations immediately noticeable.

**All input and output data used in plex gets pinned to IPFS.** See [Input / Output](io.md) for more details.

An example of content-addressed data:

```json
"protein": {
    "class": "File",
    "filepath": "6d08_protein_processed.pdb",
    "ipfs": "QmeTreLhxMmBaRqHemJcStvdyHZThdzi4gTmvTyY1igeCk"
}
```
The CID, QmeTreLhxMmBaRqHemJcStvdyHZThdzi4gTmvTyY1igeCk, can be used to access the content in multiple ways.

| Source | Access |
| ------ | ---- |
| IPFS-enabled Browser (for example, Brave) | ipfs://QmeTreLhxMmBaRqHemJcStvdyHZThdzi4gTmvTyY1igeCk |
| IPFS Desktop | QmeTreLhxMmBaRqHemJcStvdyHZThdzi4gTmvTyY1igeCk |
| IPFS http gateway | https://ipfs.io/ipfs/QmeTreLhxMmBaRqHemJcStvdyHZThdzi4gTmvTyY1igeCk
