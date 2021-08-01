# sapho

Tendermint ABCI compatible consensus protocol with hybrid proofs

> “It is by will alone I set my mind in motion. It is by the juice of sapho that thoughts acquire speed, the lips acquire stains, the stains become a warning. It is by will alone I set my mind in motion.” – Frank Herbert’s Dune

## Overview

Sapho aims to be a drop in replacement for Tendermint consensus, and will have a
full ABCI API implementation so it can be used in place of TM. It replaces the
pure pBFT validator based consensus with a hybrid proof of work/proof of stake,
and provides finality with a quorum of masternodes, whose stake grows from
participating in finalization and shrinks if they do not vote.

Sapho leverages the strong security of running the block reward lottery using a
tightly constrained proof of work that requires freezing funds to win the block
reward, and then provides finality with Masternodes.

Sapho uses a hybrid proof of work/proof of stake mining protocol with difficulty
adjustment that is derived from timestamps - that is, it will reduce difficulty
to 1 eventually and progress the chain no matter what, as the longer the
interval the lower the difficulty drops. This ensures that the chain always
progresses and that its longest block interval is less than twice the target.

In order to implement Cosmos IBC protocol, to provide finality, the chain will
have masternode validator tokens and an on-chain market for novating these and
an initial token sale denominated in Saph tokens that will be acquired by miners
and prior to this they will be operated entirely by the Sapho development team.

Masternodes will receive new block proposals, vote on which proposal should
become canonical, compare votes and then create a Schnorr multisignature on the
block that is added as the sealing transaction and whose presence indicates to
nodes that no block older can be mined on, and then the nodes will accept this
block once and for all, implementing the required immediate finality for IBC.

Masternode tokens keep a ledger of blocks certified, and for each missed vote
the newest block reward is zeroed, poor quality hosting damages the network, so
it can be continuously regulated by requiring constant participation.

Masternodes must broadcast their participation so the punishment is for breaking
this promise.

Masternode tokens accumulate a small percentage reward for each day they are not
cashed out to encourage accumulation and keep rewards off the market.

It uses the account model, as this allows a more compact ledger and the basis
for network identities.

There is no scripting engine, transactions are purely created by extensions to
the API, and have monotonic version numbers and aim for backward compatibility
so nodes that are not up-to-date still help with the consensus as much as
possible.

Just like Tendermint, while one can most easily write an application in-process
using Go, it also interfaces to the pipe and protobuf interface allowing
standard ABCI applications to run on this chain. Distributed applications are
preferably coded in compiled languages for best performance, and eventually
distributed via IPFS hosting built into nodes to provide in-band updates for
nodes, automatically or upon notification.
