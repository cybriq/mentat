# sapho

Tendermint ABCI compatible consensus protocol with hybrid proofs

> “It is by will alone I set my mind in motion. It is by the juice of sapho that thoughts acquire speed, the lips acquire stains, the stains become a warning. It is by will alone I set my mind in motion.” – Frank Herbert’s Dune

## Overview

Sapho aims to be a drop in replacement for Tendermint consensus, and will have a full ABCI API implementation so it can be used in place of TM.

Sapho uses a hybrid proof of work/proof of stake mining protocol with difficulty adjustment that is derived from timestamps - that is, it will reduce difficulty to 1 eventually and progress the chain no matter what, as the longer the interval the lower the difficulty drops. This ensures that the chain always progresses and that its longest block interval is less than twice the target.

In order to implement Cosmos IBC protocol, to provide finality, the chain will have masternode tokens and an on-chain market for novating these and an initial token sale denominated in Saph tokens that will be acquired by miners and prior to this they will be operated entirely by the Sapho development team. Masternodes will receive new block proposals, vote on which proposal should become canonical, compare votes and then create a Schnorr multisignature on the block that is added as the sealing transaction and whose presence indicates to nodes that no block older can be mined on, and then the nodes will accept this block once and for all, implementing the required immediate finality for IBC.

It uses the account model, as this allows a more compact ledger and the 
basis for network identities.

There is no scripting engine, transactions are purely created by extensions 
to the API, and have monotonic version numbers and aim for backward 
compatibility so nodes that are not up-to-date still help with the consensus 
as much as possible.