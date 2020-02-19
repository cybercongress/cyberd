
# Cyber - Knowledge consensus computer for The Great Web

<div align="center">
 <img src="docs/img/header.png" width="320" />
</div>

### Code
[![version](https://img.shields.io/github/release/cybercongress/cyberd.svg?style=flat-square)](https://github.com/cybercongress/cyberd/releases/latest)
[![CircleCI](https://img.shields.io/circleci/project/github/cybercongress/cyberd.svg?style=flat-square)](https://circleci.com/gh/cybercongress/cyberd/tree/master)
[![license](https://img.shields.io/badge/License-Cyber-brightgreen.svg?style=flat-square)](https://github.com/cybercongress/cyberd/blob/master/LICENSE)
![Cosmos-SDK](https://img.shields.io/static/v1.svg?label=cosmos-sdk&message=0.37.3&color=blue&style=flat-square)
![Tendermint](https://img.shields.io/static/v1.svg?label=tendermint&message=0.32.7&color=blue&style=flat-square)
[![LoC](https://tokei.rs/b1/github/cybercongress/cyberd?style=flat)](https://github.com/cybercongress/cyberd)
[![contributors](https://img.shields.io/github/contributors/cybercongress/cyberd.svg?style=flat-square)](https://github.com/cybercongress/cyberd/graphs/contributors)
[![Coverage Status](https://img.shields.io/coveralls/github/cybercongress/cyberd/master?style=flat-square)](https://coveralls.io/github/cybercongress/cyberd?branch=master)

### Blockchain
[![chain](https://img.shields.io/badge/Chain-Euler--Dev-success.svg?style=flat-square)](https://github.com/cybercongress/cyberd/blob/master/docs/run_validator.md)
[![block](https://img.shields.io/badge/dynamic/json?color=blue&label=Block%20Height&query=%24.result.height&url=https://titan.cybernode.ai/api/index_stats&style=flat-square)]()
[![cyberlinks](https://img.shields.io/badge/dynamic/json?color=blue&label=Cyberlinks&query=%24.result.linksCount&url=https://titan.cybernode.ai/api/index_stats&style=flat-square)]()
[![cids](https://img.shields.io/badge/dynamic/json?color=blue&label=CIDs&query=%24.result.cidsCount&url=https://titan.cybernode.ai/api/index_stats&style=flat-square)]()
[![agents](https://img.shields.io/badge/dynamic/json?color=blue&label=Web3%20Agents&query=%24.result.accsCount&url=https://titan.cybernode.ai/api/index_stats&style=flat-square)]()
[![validators](https://img.shields.io/badge/dynamic/json?label=Validators&query=%24.result.validators.length&url=https://titan.cybernode.ai/api/validators%3F&style=flat-square)]()

### Community
[![telegram](https://img.shields.io/badge/Join%20Us%20On-Telegram-2599D2.svg?style=flat-square)](https://t.me/fuckgoogle)
[![gitcoin](https://img.shields.io/badge/Join%20Us%20On-Gitcoin-2599D2.svg?style=flat-square)](https://t.me/fuckgoogle)
[![twitter](https://img.shields.io/twitter/follow/cyber_devs?label=Follow)](https://twitter.com/@cyber_devs)
[![reddit](https://img.shields.io/reddit/subreddit-subscribers/cybercongress?style=social)](https://www.reddit.com/r/cybercongress)

## Why

Google is the most powerful religion ever, and now is the time to abandon it. The Great Web is coming. A search system that will drive its growth, needs to be censorless, open and relevant. 

<div align="center">
 <img src="./docs/img/cyber.png"/>
</div>

## What is Cyber

Cyber is a knowledge consensus computer for answers and a search engine. Cyber is defined by web3 agents. It is a distributed consensus supercomputer that runs on top of [IPFS](). It can compute relevance for [web3]() data via the use of [cyberlinks](), that are created from IPFS hashes. Essentialy, those hashes represent [Content IDentificators (CIDs)] weighted by [Page Rank]():

```
CID1 -----> CID2
```

Cyberlinks are committed by Web3 agents. They are links between two CIDs. In its current implementation, a CID is an IPFS hash of CIDv0 or of CIDv1. A web3 agent can link a CID of any keyword, app, etc - with another CID. Then, the agent creates a link between the two hashes with a weight corresponding to a stake.

All the cyberlinks with a given weight are stored within the knowledge graph. The graph is re-computed by the validators every given number of blocks. For calculations, we've implemented `the proof of relevance` root-hash, which is computed from CID rank values (which are computed on CUDA GPUs every round). 

Cyberd is the first implementation of the cyber protocol. It is based on [cosmos-SDK](https://github.com/cosmos/cosmos-sdk) and [tendermint BFT Consensus](https://github.com/tendermint/tendermint). 

This implementation uses a very simple bandwidth model. The main goal of the model is to reduce the daily networks growth to a given constant.

Thus, here we introduce [resource credits(RS)](). Each message, of a transaction type - for example, a `link` or a `send` are assigned an RS cost. We call this `Bandwidth`. A users bandwidth depends on their balance and is equal to the sum of their liquid and staked tokens. The users bandwidth is a recoverable value. Full recovery of the bandwidths, from 0 to maximum value - takes `RecoveryWindow` blocks.

There is a period `AdjustPricePeriod`, summing how much RS or bandwidth was spent during that period (`RecoveryPeriod`). The `SpentBandwidth/DesirableBandwidth` ratio defines the so-called current `price multiplier`. If the network usage is low, the `price multiplier` adjusts the message cost (simply by multiplying it) to allow a user with a lower stake to do more transactions. If the demand for resources increases, the `price multiplier` goes `>1` thus, increasing message cost, and limiting final tx count for a longer-term period (RC recovery will be `<` then RC spending).


## For validators

Each validator participates in the Tendermint consensus and computes/validates cyber~Rank within the knowledge graph.


## For rank providers

Rank providers crawl/index the web. They then, cyberlinks CIDs of any given data to cyberd by consuming their bandwidth.


## For search users

A valuable, censorship-resistant and a provable search of the web for any kind of intellegent species. The search is transactionally-based and only possible if the agent has enough bandwidth.


## For developers

The chance to create a new and a decentralized Google with affiliated services, for example: SEO, crawlers, web indexers, decentralized platforms and so on. You have the chance of changing the web. 

## For data/content producers

The opportunity to move their content to web3 and save it from any type of censorship. Your content is yours. Make sure others will be able to see it.

## For miners / GPUs holders

With the growth of the network, we will need cards, a lot of cards. Join us.


## cyber•Rank
```
0. cyber•Rank - a token weighted Page Rank (initial implementation).
1. A knowledge graph consists of CIDs, which are connected by cyberlinks. 
2. A cyberlink may be cast only once, between any given CIDs. 
3. The weight of a cyberlink is the token balance of the agent. 
4. The rank is computed within the current rank calculation window, on CUDA kernel.
5. After cyber~Ranks is computed, each CID recieves a given rank.
6. [Very important] Rank computation is based on the current (computation window) balances of the agents .
```

## Bandwidth
```
0. The network has a desirable network bandwidth (max bandwidth).
1. An agent's bandwidth is proportional to the stake that he owns against the total network supply.
2. Linking and any other chain operations - consume bandwidth.
3. It takes RecoveryPeriod blocks for a full recovery of an agent's bandwidth.
4. The network gives a discount of up to 100X of the operational costs for incentivized loads.
```

## Search index and proofs
```
0. A node can be launched with an "ALLOW SEARCH" mode. This allows to search the knowledge graph, using this node .
1. A node in search mode, also constructs a full Merkle tree for cyberlinks and calculates their ranks.
2. Proof of rank of any given CID - provides a Merkle path. This allows the client to validate the returned rank of any given CID and the existence of a given cyberlink.
3. The Merkle root of a rank, and the Merke roots of links Merkle trees - are used for calculating the app hash for each block and are part of the protocol/consensus.
```

## Technologies
```
1. Tendermint
2. Cosmos-SDK
3. NVIDIA CUDA
4. IPFS
```

## Paper
```
```

## Status

**WARNING**: Project status: **testnet**. We are at research state at the moment.
Please, read our [whitepaper](./docs/cyberd.md)

## cyberd Public Testnets

To run a full-node or a validator node, on the latest public testnet of cyberd, please follow [this guide](https://cybercongress.ai/docs/cyberd/run_validator/).

## Explorers

### Euler-Dev Testnet
The [cyberd](https://callisto.cybernode.ai/) explorer is based on [bigDipper](https://cosmos.bigdipper.live) by [Forbole](https://www.forbole.com/)

## Peers & Seeds

### Euler-Dev Testnet
```
d0a148810b8b0e6e5bd16ea3ede1e3a7851208b9@titan.cybernode.ai:26656  
```

## Build cyberd and cyberdcli (Go 1.13)
```
git clone github.com/cybercongress/cyberd
cd cyberd
make

or 

go build -tags cuda,netgo,ledger -o cyberd ./cmd/cyberd
go build -tags netgo,ledger -o cyberdcli ./cmd/cyberdcli
```
Build of a cyber node, requires CUDA toolkit with the latest Nvidia drivers installed on your machine. 

You may build just the CLI to communicate with dedicated cyber~Nodes.

## Build GPU kernels
```
cd x/rank/cuda
make
```
It will build libcbdrank.so and copy them with cbdrank.h to /usr/lib

## Use CLI with dedicated node
```
cyberdcli q account <account> --chain-id euler-dev --node https://titan.cybernode.ai:26657
```

##  The 10 min Devenv Setup for Developers and Gitcoiners 

Take the "10 minutes challenge" with [set up dev environment](https://cybercongress.ai/docs/cyberd/setup_dev_env/).

## Ledger and Multisig guides
1. [Ledger with CLI]()
2. [Ledger for validator setup]()
3. [Multisig guide]()

## Game of Links

The "game of links" is a game between cyber~Congress and between Cosmos stakeholders for a place in Genesis. It should bootstrap and load the network during the Euler-5 testnet. The greatest project will come on top with a significant number of followers. The game is finished if both of the following criteria are met:

**\>** 146 validators are in consensus for a period of 10k blocks

**\>** 600000 ATOM donated or 90 days have passed

#### Goals:
```
0. Intentiavize professional and long-term committed validators
1. Bootstrap a full validator set, which will drive and grow the network and the knowledge graph. This will help to build a strong community that is backed by delegators
2. Initialize the knowledge graph with valuable knowledge domains
3. Develop tools for crawling and indexing specific knowledge domains
4. Start to form responsibility and activity in relation to the governance of the network and to decision-making mechanisms.
5. Distribute tokens in the most valuable form
```

#### Rewards for:
```
0. Summary uptime of every validator
1. The load of an agent (the consumed networks bandwidth)
2. Amount of delegated tokens to a validator (Game of Stakes)
3. The relevance of links submitted (TOP-1000 CIDs)
```

## Distribution
```
```

## Docs

Explore the docs in our [knowledge base](https://cybercongress.ai/docs/cyberd/cyberd/).

## IBC
```
Waiting for Game of Zones and Cosmos-SDK v0.38.0 release
```

## Inter Knowledge Protocol for Relevance Machines (over IBC)
```
Protocol prototyping on:
1. Subgraph Interchange
2. CIDs Interchange
3. Ranks Interchange
```

## Community

**\>** [devChat](https://t.me/fuckgoogle) for web3 agnets

**\>** [Twitter](https://twitter.com/cyber_devs) for updates and memes

**\>** [Own blog](https://cybercongress.ai/post/) with RSS and useful articles


## Research and Development
```
- cyber~Rank scaling
— Onchain upgrades
— IBC
- IKP
— Universal oracle
— WASM VM for gas
— CUDA VM for gas
— Privacy
- PoRep/PoST
- Autonomous onchain agents
```

## Let's #fuckgoogle together
```
Bring your PRs and text us to take on board.
```

## Issues

If you have any problems with, or questions, about search - please contact us via
a [GitHub issue](https://github.com/cybercongress/cyberd/issues).

## Contribute

You are invited to contribute new features, fixes, or updates - large or small. We are always thrilled to receive pull requests and do our best to process them as fast as we can. You can find detailed information in our
[contribution guide](./docs/contributing/contributing.md).

## Gitcoin program

We want to reward you for your contributions! We constantly fund our issues on [gitcoin](https://gitcoin.co/profile/cybercongress) and attach good descriptions to them. Along with the current state of the project  and with user stories. We try to answer comments regularly in issues and in our [devChat](https://t.me/fuckgoogle).

<a href="https://gitcoin.co/explorer?q=cyberd">
 <img src="https://gitcoin.co/funding/embed?repo=https://github.com/cybercongress/cyberd">
</a>

## Changelog

Stay tuned with our [Changelog](./CHANGELOG.md).

## Team
<table>
  <tr>
    <td align="center"><a href=https://github.com/xhipster><img src="https://avatars0.githubusercontent.com/u/410789?s=400&v=4" width="100px;" alt="xhipster"/><br /><sub><b>Dima Starodubcev</b></sub></a><br /><a href="https://github.com/cybercongress/cyberd/commits?author=xhipster" title="Documentation">📖</a> <a href="#maintenance-xhipster" title="Maintenance">🚧</a></td><td align="center"><a href="https://github.com/litvintech"><img src="https://avatars2.githubusercontent.com/u/1690657?v=4" width="100px;" alt="Valery Litvin"/><br /><sub><b>Valery Litvin</b></sub></a><br /><a href="https://github.com/cybercongress/cyberd/commits?author=litvintech" title="Code">💻</a> <a href="#projectManagement-litvintech" title="Project Management">📆</a> <a href="https://github.com/cybercongress/cyberd/commits?author=litvintech" title="Documentation">📖</a></td><td align="center"><a href="https://github.com/SaveTheAles"><img src="https://avatars0.githubusercontent.com/u/36516972?v=4" width="100px;" alt="Ales Puchilo"/><br /><sub><b>Ales Puchilo</b></sub></a><br /><a href="https://github.com/cybercongress/cyberd/commits?author=SaveTheAles" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/mrlp4"><img src="https://avatars2.githubusercontent.com/u/38989909?s=400&v=4" width="100px;" alt="Kiryl Laptanovich"/><br /><sub><b>Kiryl Laptanovich</b></sub></a><br /><a href="https://github.com/cybercongress/cyberd/commits?author=mrlp4" title="Tests">⚠️</a></td>
   <td align="center"><a href="https://github.com/serejandmyself"><img src="https://avatars2.githubusercontent.com/u/7550961?s=400&v=4" width="100px;" alt="Sergey Simanovsky"/><br /><sub><b>Sergey Simanovsky</b></sub></a><br /><a href="https://github.com/cybercongress/cyberd/commits?author=serejandmyself" title="Ecosystem">:milky_way:</a></td>
  </tr>
</table>

## Linked Projects

**\>** [Cosmos](https://github.com/cosmos/)

**\>** [Tendermint](https://github.com/tendermint/)

**\>** [IPFS](https://github.com/ipfs/)

## License

Cyber License - Don’t believe, don’t fear, don’t ask.

We will be happy if you fork and launch your own network and set up a knowledge graph. Eventually, we will meet each other somewhere in the metaverse (with the help of IBC, of course).

## The End of Google

<div align="center">
 <img src="docs/img/End-of-Google.jpg" width="600" />
</div>

<div align="center">
 <sub>Built by
 <a href="https://twitter.com/cyber_devs">cyber•Congress</a> and
 <a href="https://github.com/cybercongress/cyberd/graphs/contributors">contributors</a>
</div>
