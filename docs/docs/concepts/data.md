---
title: Content-Addressed Data
sidebar_position: 2
sidebar_label: Data
---

Plex utilizes a decentralized storage protocol, [**IPFS**](https://docs.ipfs.tech/), for managing file storage in its scientific computing workflows. Within IPFS, all data is content-addressed, meaning each file is given a unique content identifier ([**CID**](https://docs.ipfs.tech/concepts/content-addressing/#what-is-a-cid)).

CIDs are derived from a file's content rather than its location.

Using CIDs not only enhances file retrieval but also promotes data integrity since the identifier changes if the content does, making any alterations immediately noticeable.

**Plex [pins](https://docs.ipfs.tech/how-to/pin-files/) all input and output data to IPFS.** See [Input / Output](io.md) for more details.

An example of content-addressed data:

```json
"protein": {
    "class": "File",
    "filepath": "6d08_protein_processed.pdb",
    "ipfs": "QmeTreLhxMmBaRqHemJcStvdyHZThdzi4gTmvTyY1igeCk"
}
```
The CID, **QmeTreLhxMmBaRqHemJcStvdyHZThdzi4gTmvTyY1igeCk**, can be used to access the content in multiple ways.

| Source | Access |
| ------ | ---- |
| IPFS-enabled browser (ie, [Brave](https://brave.com/ipfs-support/)) | ipfs://QmeTreLhxMmBaRqHemJcStvdyHZThdzi4gTmvTyY1igeCk |
| IPFS Desktop | QmeTreLhxMmBaRqHemJcStvdyHZThdzi4gTmvTyY1igeCk |
| IPFS http gateway | http://bacalhau.labdao.xyz:8080/ipfs/QmeTreLhxMmBaRqHemJcStvdyHZThdzi4gTmvTyY1igeCk
