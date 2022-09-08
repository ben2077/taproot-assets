// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package sqlite

import (
	"database/sql"
	"time"
)

type Addr struct {
	ID           int32
	Version      int16
	AssetID      []byte
	FamKey       []byte
	ScriptKeyID  int32
	TaprootKeyID int32
	Amount       int64
	AssetType    int16
	CreationTime time.Time
}

type Asset struct {
	AssetID                  int32
	Version                  int32
	ScriptKeyID              int32
	AssetFamilySigID         sql.NullInt32
	ScriptVersion            int32
	Amount                   int64
	LockTime                 sql.NullInt32
	RelativeLockTime         sql.NullInt32
	SplitCommitmentRootHash  []byte
	SplitCommitmentRootValue sql.NullInt64
	AnchorUtxoID             sql.NullInt32
}

type AssetFamily struct {
	FamilyID       int32
	TweakedFamKey  []byte
	InternalKeyID  int32
	GenesisPointID int32
}

type AssetFamilySig struct {
	SigID      int32
	GenesisSig []byte
	GenAssetID int32
	KeyFamID   int32
}

type AssetMintingBatch struct {
	BatchID            int32
	BatchState         int16
	MintingTxPsbt      []byte
	MintingOutputIndex sql.NullInt16
	GenesisID          sql.NullInt32
	CreationTimeUnix   time.Time
}

type AssetProof struct {
	ProofID   int32
	AssetID   int32
	ProofFile []byte
}

type AssetSeedling struct {
	SeedlingID      int32
	AssetName       string
	AssetType       int16
	AssetSupply     int64
	AssetMeta       []byte
	EmissionEnabled bool
	AssetID         sql.NullInt32
	BatchID         int32
}

type AssetWitness struct {
	WitnessID            int32
	AssetID              int32
	PrevOutPoint         []byte
	PrevAssetID          []byte
	PrevScriptKey        []byte
	WitnessStack         []byte
	SplitCommitmentProof []byte
}

type ChainTxn struct {
	TxnID       int32
	Txid        []byte
	RawTx       []byte
	BlockHeight sql.NullInt32
	BlockHash   []byte
	TxIndex     sql.NullInt32
}

type GenesisAsset struct {
	GenAssetID     int32
	AssetID        []byte
	AssetTag       string
	MetaData       []byte
	OutputIndex    int32
	AssetType      int16
	GenesisPointID int32
}

type GenesisPoint struct {
	GenesisID  int32
	PrevOut    []byte
	AnchorTxID sql.NullInt32
}

type InternalKey struct {
	KeyID     int32
	RawKey    []byte
	KeyFamily int32
	KeyIndex  int32
}

type Macaroon struct {
	ID      []byte
	RootKey []byte
}

type ManagedUtxo struct {
	UtxoID           int32
	Outpoint         []byte
	AmtSats          int64
	InternalKeyID    int32
	TapscriptSibling []byte
	TaroRoot         []byte
	TxnID            int32
}

type MssmtNode struct {
	HashKey   []byte
	LHashKey  []byte
	RHashKey  []byte
	Key       []byte
	Value     []byte
	Sum       int64
	Namespace string
}

type MssmtRoot struct {
	Namespace string
	RootHash  []byte
}
