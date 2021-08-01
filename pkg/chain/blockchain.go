package chain

import "github.com/tendermint/tendermint/abci/types"

type SaphoConsensus struct {
	// MinBlockInterval is the minimum difference between a new block and
	// the parent block
	MinBlockInterval int
	// BlockIntervalTarget is the set point for the difficulty adjustment,
	// important for regulating the token emission rate
	BlockIntervalTarget int
	// AdjustmentInverseBase is the reciprocal of the scaling value
	// compounded with the difference between the block timestamp and the
	// time target computed by the difficulty for a given parent block and
	// the timestamp of the proposed block.
	//
	// The closer the timestamp is to the computed target timestamp the less
	// that difficulty is adjusted.
	AdjustmentInverseBase int
	// EMAWindow is the number of previous blocks used to compute the
	// average used to derive the difficulty
	EMAWindow int
	// EMANumer is the numerator of the EMA smoothing factor (multiplies)
	EMANumer int
	// EMADenom is the denominator of the EMA smoothing factor (divides)
	EMADenom int
	// StakeCoolDown is the number of seconds that stake is subtracted
	// from available balance when querying to perform another
	// transaction. During this period the stake cannot be spent or staked.
	StakeCoolDown int
}

// GetDefaultSaphoConsensus returns a default configured Sapho consensus
// structure
func GetDefaultSaphoConsensus() (sc *SaphoConsensus) {
	return &SaphoConsensus{
		MinBlockInterval:      15,
		BlockIntervalTarget:   30,
		AdjustmentInverseBase: 10,
		EMAWindow:             120,
		EMANumer:              2,
		EMADenom:              6,
	}
}

type Blockchain struct {
	abciParams *types.ConsensusParams
	genesis    *Block
	params     *SaphoConsensus
}

func NewBlockchain(abciParams *types.ConsensusParams, genesis *Block) *Blockchain {
	return &Blockchain{
		abciParams: abciParams, genesis: genesis,
	}
}
