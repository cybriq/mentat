package chain

import (
	. "github.com/tendermint/tendermint/abci/types"
	"github.com/zeebo/blake3"
)

type ABCI struct {
	*Blockchain
}

func (a *ABCI) Info(reqInfo RequestInfo) (ri ResponseInfo) {
	return
}

func (a *ABCI) SetOption(reqSetOpt RequestSetOption) (rso ResponseSetOption) {
	return
}

func (a *ABCI) Query(reqQuery RequestQuery) (rq ResponseQuery) {
	return
}

func (a *ABCI) CheckTx(reqChkTx RequestCheckTx) (rct ResponseCheckTx) {
	return
}

func (a *ABCI) InitChain(reqInitChain RequestInitChain) (
	ric ResponseInitChain) {
	// const maxBytes = 1 << 22 // 4Mb
	//
	// rr := RequestInitChain{
	// 	Time:    time.Time{},
	// 	ChainId: "",
	// 	ConsensusParams: &ConsensusParams{
	// 		Header: &BlockParams{
	// 			MaxBytes: 0,
	// 			MaxGas:   0,
	// 		},
	// 		Evidence: &types.EvidenceParams{
	// 			MaxAgeNumBlocks: 0,
	// 			MaxAgeDuration:  0,
	// 			MaxBytes:        0,
	// 		},
	// 		Validator: &types.ValidatorParams{
	// 			PubKeyTypes: []string{typs.ABCIPubKeyTypeEd25519},
	// 		},
	// 		Version: &types.VersionParams{
	// 			AppVersion: 0,
	// 		},
	// 	},
	// 	Validators:    nil,
	// 	AppStateBytes: nil,
	// 	InitialHeight: 0,
	// }

	a.Blockchain.abciParams = reqInitChain.ConsensusParams

	appHash := blake3.Sum256(reqInitChain.AppStateBytes)
	ric = ResponseInitChain{
		ConsensusParams: reqInitChain.ConsensusParams,
		// Validators:      reqInitChain.Validators,
		AppHash: appHash[:],
	}
	return
}

func (a *ABCI) BeginBlock(reqBegBlk RequestBeginBlock) (
	rbb ResponseBeginBlock) {
	return
}

func (a *ABCI) DeliverTx(reqDelTx RequestDeliverTx) (rdt ResponseDeliverTx) {
	return
}

func (a *ABCI) EndBlock(reqEndBlk RequestEndBlock) (reb ResponseEndBlock) {
	return
}

func (a *ABCI) Commit() (rc ResponseCommit) {
	return
}

func (a *ABCI) ListSnapshots(reqListSnap RequestListSnapshots) (
	rls ResponseListSnapshots) {
	return
}

func (a *ABCI) OfferSnapshot(
	reqOfferSnap RequestOfferSnapshot,
) (ros ResponseOfferSnapshot) {
	return
}

func (a *ABCI) LoadSnapshotChunk(
	reqLoadSnapChunk RequestLoadSnapshotChunk,
) (rlsc ResponseLoadSnapshotChunk) {
	return
}

func (a *ABCI) ApplySnapshotChunk(
	reqApplySnapChunk RequestApplySnapshotChunk,
) (rasc ResponseApplySnapshotChunk) {
	return
}
