package chain

import (
	"math/big"

	"github.com/niubaoshu/gotiny"
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

// An Address is a fixed length array representing the blake 3 256-bit hash of
// the public key truncated to 20 bytes used to identify an account in the
// database, created by sending funds to the address, generated from the public
// key, derived from the private key
type Address [20]byte

// Difficulty is a variable length long integer of maximum 32 byte length
// representing the difference between minimum difficulty and the block proof
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

// BigToDifficulty converts big. Int from a proof or target into a difficulty
// value, used to compute a new difficulty adjustment and return a value for the
// block
func BigToDifficulty(bi *big.Int) Difficulty {
	bb := new(big.Int).Sub(MinDifficulty, bi)
	b := bb.Bytes()
	if len(b) > 32 {
		panic("difficulty cannot be more than 32 bytes")
	}
	return b
}

// A BlockNode is a node in a dual linked list that allows walking the chain
// history for various purposes
type BlockNode struct {
	// parent is the pointer to the stored block structure,
	// if nil it needs to be read in from the database
	parent *Block
	// Block is the block itself
	*Block
	// children are the descendents of the block, there can be multiple
	// during the process of finalization by masternode validators.
	//
	// The one that gets picked and ends up as the only node after this
	// block will tend to be the one that is first on the most nodes.
	children []*Block
}

// A Block represents a unit of progress for the network, and carries a payload
// of transactions to change the state of the accounts database for users
type Block struct {
	Header
	// Transactions is the list of transactions ordered in the merkle
	// tree structure with the necessary empty nodes in the list as
	// required. The MerkleRoot is computed from this list which must be
	// structured correctly for this.
	Transactions []interface{}
	// MasternodeSignature - After a block is proposed, masternodes run a
	// vote/compare/sign cycle which generates a Schnorr signature with each
	// masternode's signature overlaid on it, the masternode list being
	// visible as the issued tokens and their currently controlling account.
	//
	// This signature indicates to nodes that no previous block can now be
	// mined on, making the block immediately final, a requirement for
	// Cosmos IBC protocol.
	//
	// This signature is the product of the peer validation process for
	// signatures on the block, the first signature chain to acquire the 2/3
	// quorum of masternodes triggers the acceptance and finality of the
	// new block.
	//
	// The signature is applied to the Hash of the Header, which a verifier
	// must generate from the serialized block header.
	// Without this signature a block may not be mined on,
	// and the masternodes only accept one and finalize it immediately
	// before a new block is allowed to happen (15 second minimum)
	MasternodeSignature []byte
}

// Header is the metadata for the block and the entropy for mining
type Header struct {
	// Parent is the hash of the previous block
	Parent Hash
	// Time is the Unix timestamp of the block
	Time int64
	// Provenance is a value computed by the app state hash of the block
	// nearest the staking cooldown period, as displaced from the block's
	// timestamp, combined with the coinbase of the parent block, two values
	// that don't exist until a block has been mined.
	//
	// Before the cooldown period the genesis block is always the app
	// state hash used to generate this
	Provenance Hash
	// Difficulty is a value that is subtracted from the lowest difficulty
	// (all but the largest bit 1 in a 256 bit hash) to create a difficulty
	// target which defines the largest proof that can be accepted.
	//
	// If the following is true:
	//
	//   stake * proof < target
	//
	// then the block has valid proof and could be accepted.
	Difficulty Difficulty
	// Coinbase is the account to which the reward for this block is paid
	Coinbase Address
	// Stake is the necessary amount of stake to normalize the difficulty
	// target, minimum of 1. This functions as a multiplier on the proof
	// to produce the value that is compared to the difficulty target
	Stake uint64
	// MerkleRoot is the root hash of the transaction merkle tree
	MerkleRoot Hash
}

// Serialize the block
func (h *Header) Serialize() (serialized []byte) {
	return gotiny.Marshal(h)
}

// Deserialize the block
func Deserialize(serialized []byte) (h *Header) {
	gotiny.Unmarshal(serialized, h)
	return
}

func (h *Header) Hash() (ha Hash) {
	return blake3.Sum256(h.Serialize())
}

// GetStake returns the multiplier that makes the Hash of the Header equal to or
// below the value generated by the Difficulty, which is locked up for the
// cool down period (SaphoConsensus.StakeCoolDown).
//
// If the proof is divided by the stake it will produce a number smaller than
// the value generated from Difficulty, and thus pass
func (h *Header) GetStake() (stake uint64) {
	target := h.Difficulty.ToBigInt()
	proof := h.Hash().ToBigInt()
	if target.Cmp(proof) < 0 {
		multiplier, mod := new(big.Int).DivMod(target, proof, new(big.Int))
		stake = multiplier.Uint64()
		// round the stake up if there is a remainder in the division
		//
		if mod.Cmp(big.NewInt(0)) > 0 {
			stake++
		}
		return
	}
	return 0
}

// ValidateStake divides the proof by the stake and if the result is smaller
// than the target then it returns true, otherwise false.
//
// Miners should use GetStake to determine the stake that will be used, so they
// stake the minimal required for the solution, and reject solutions that
// require more stake than they want to place on one block.
//
// More stake can be used but it makes no sense to lock up more than required.
// Staking's main benefit is getting ahead of the crowd for mining a block as it
// effectively lowers the target.
func (h *Header) ValidateStake() bool {
	target := h.Difficulty.ToBigInt()
	// the effective proof is the hash of the Block Header divided by the
	// Stake
	proof := new(big.Int).Div(h.Hash().ToBigInt(), big.NewInt(int64(h.Stake)))
	if target.Cmp(proof) < 0 {
		return true
	}
	return false
}
