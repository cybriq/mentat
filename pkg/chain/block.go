package chain

import (
	"math/big"

	"github.com/zeebo/blake3"
)

// Hash is your standard 32 byte fixed length array as returned by many hash
// functions and used in the blocks as identifiers of other data
type Hash [32]byte

// ToBigInt converts a hash to big.Int so it can be compared as in the case of a
// proof or other partial hash collision
func (h Hash) ToBigInt() *big.Int {
	return new(big.Int).SetBytes(h[:])
}

// ToAddress returns the address derived from a hash (which would be a public
// key, presumably)
func (h Hash) ToAddress() (a Address) {
	H := blake3.Sum256(h[:])
	hh := H[:]
	copy(a[:], hh[:20])
	return
}

// An Address is a fixed length array representing the 20 byte blake 3 256-bit
// hash truncated to 20 bytes used to identify an account in the database,
// created by sending funds to the address, generated from the public key,
// derived from the private key
type Address [20]byte

// Difficulty is a variable length long integer
type Difficulty []byte

// MinTargetBytes is the difficulty of 1 as a 32 byte slice, the largest number
// a proof can be without triggering a divide by zero error, and the base from
// which the difficulty value is subtracted to yield a value to compare to a
// block proof
var MinTargetBytes = []byte{
	0xef, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
}

// MinDifficulty represents the bytes of a proof that cannot be any smaller, it
// has the most significant bit as a zero, meaning it is one less than the
// completely impossible all bits 1 coming from a hash (same impossible as all
// zero).
var MinDifficulty = new(big.Int).SetBytes(MinTargetBytes)

// ToBigInt converts a difficulty to big.Int by subtracting it from
// MinDifficulty
func (d Difficulty) ToBigInt() *big.Int {
	if len(d) > 32 {
		panic("difficulty cannot be more than 32 bytes")
	}
	return new(big.Int).Sub(MinDifficulty, new(big.Int).SetBytes(d))
}

// BigToDifficulty converts big.Int into a difficulty value, used to compute a
// new difficulty adjustment and return a value for the block
func BigToDifficulty(bi *big.Int) Difficulty {
	b := bi.Bytes()
	if len(b) > 32 {
		panic("difficulty cannot be more than 32 bytes")
	}
	return b
}

// A Block represents a unit of progress for the network, and carries a payload
// of transactions to change the state of the accounts database for users
type Block struct {
	// Parent is the hash of the previous block
	Parent Hash
	// Time is the Unix timestamp of the block
	Time int64
	// Provenance is a value computed by the hash of the block nearest the
	// staking cooldown period, as displaced from the block's
	// timestamp, combined with the coinbase of the parent block,
	// two values that don't exist until a block has been mined
	Provenance Hash
	// Difficulty is a value that is subtracted from the lowest difficulty
	// (all but the largest bit 1 in a 256 bit hash) to create a difficulty
	// target below which proofs must be or be staked equal differing by an
	// amount covered by the stake required to issue the block
	Difficulty Difficulty
	// Coinbase is the account to which the reward for this block is paid
	Coinbase Address
	// Stake is the necessary amount of stake to normalize the difficulty
	// target, minimum of 1
	Stake uint64
	// MerkleRoot is the root hash of the transaction Merkle Tree
	MerkleRoot Hash
	// MasternodeSignature - After a block is proposed, masternodes run a
	// vote/compare/sign cycle which generates a Schnorr signature with each
	// masternode's signature overlaid on it, the masternode list being
	// visible as the issued tokens and their currently controlling account.
	//
	// This signature indicates to nodes that no previous block can now be
	// mined on, making the block immediately final, a requirement for
	// Cosmos IBC protocol
	MasternodeSignature Hash
	// transactions is the working memory storage for the individual
	// transactions that generate the merkle root above
	transactions []interface{}
}
