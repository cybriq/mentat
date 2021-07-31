package chain

import "math/big"

type Hash [32]byte

func (h Hash) ToBigInt() *big.Int {
	return new(big.Int).SetBytes(h[:])
}

type Address [20]byte

type Difficulty []byte

var MaxDifficulty = new(big.Int).Lsh(big.NewInt(1), 256)

func (d Difficulty) ToBigInt() *big.Int {
	return new(big.Int).Sub(MaxDifficulty, new(big.Int).SetBytes(d))
}

type Block struct {
	Parent       Hash
	Time         int64
	Nonce        uint64
	Difficulty   Difficulty
	Coinbase     Address
	Stake        uint64
	Transactions []interface{}
}
