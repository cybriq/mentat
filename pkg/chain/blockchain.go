package chain

import "github.com/tendermint/tendermint/abci/types"

type Blockchain struct {
	abciParams *types.ConsensusParams
	genesis    *Block
}

func NewBlockchain(abciParams *types.ConsensusParams, genesis *Block) *Blockchain {
	return &Blockchain{
		abciParams: abciParams, genesis: genesis,
	}
}
