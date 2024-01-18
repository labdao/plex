# Lab Exchange 🧫×🧬→💊
⚡ **Run highly reproducible computational biology applications on top of a decentralised compute and storage network.** ⚡

<p align="left">
    <a href="https://github.com/labdao/plex/blob/main/LICENSE.md" alt="License">
        <img src="https://img.shields.io/badge/license-MIT-green" />
    </a>
    <a href="https://github.com/labdao/plex/releases/" alt="Release">
        <img src="https://img.shields.io/github/v/release/labdao/plex?display_name=tag" />
    </a>
    <a href="https://github.com/labdao/plex/pulse" alt="Activity">
        <img src="https://img.shields.io/github/commit-activity/m/labdao/plex" />
    </a>
    <a href="https://img.shields.io/github/downloads/labdao/plex/total">
        <img src="https://img.shields.io/github/downloads/labdao/plex/total" alt="total download">
    </a>
    <a href="https://labdao.xyz/">
        <img alt="LabDAO website" src="https://img.shields.io/badge/website-labdao.xyz-red">
    </a>
    <a href="https://twitter.com/intent/follow?screen_name=lab_dao">
        <img src="https://img.shields.io/twitter/follow/lab_dao?style=social&logo=twitter" alt="follow on Twitter">
    </a>
    <a href="https://discord.gg/labdao" alt="Discord">
        <img src="https://dcbadge.vercel.app/api/server/labdao?compact=true&style=flat-square" />
    </a>
</p>

The Lab Exchange is a full web stack for distributed computational biology.
* 🌎 **Build once, run anywhere:** The Lab Exchange is using distributed compute and storage to run containers on a public network. Need GPUs? We got you covered.
* 🔍 **Content-addressed by default:** Every file processed by plex has a deterministic address based on its content. Keep track of your files and always share the right results with other scientists.
* 🪙 **Records: authorship tracking built-in** Every compute event on Lab Exchange is mintable as an on-chain token that grants the holder rights over the newly generated data.
* 🔗 **Strictly composable:** Every tool in plex has declared inputs and outputs. Plugging together tools by other authors should be easy.

Plex is built with [Bacalhau](https://www.bacalhau.org/) and [IPFS](https://ipfs.tech/).

## Running the App Locally

We have `docker-compose` files available to bring up the stack locally.

Note:
* Only `amd64` architecture is currently supported.

# Build and bring up stack (CPU)
```
docker compose up -d --wait
```

# Build and bring up stack (GPU)
```
docker compose -f docker-compose.yml -f docker-compose-gpu.yml up -d --wait --build
```

To run `bacalhau` cmds against local environment simply set `BACALHAU_API_HOST=127.0.0.1` in your terminal