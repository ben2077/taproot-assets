package taropsbt

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/lightninglabs/taro/address"
	"github.com/lightninglabs/taro/asset"
	"github.com/lightninglabs/taro/commitment"
	"github.com/lightninglabs/taro/internal/test"
	"github.com/lightningnetwork/lnd/keychain"
	"github.com/stretchr/testify/require"
)

var (
	testParams = &address.MainNetTaro
)

func randomPacket(t testing.TB) *VPacket {
	testPubKey := test.RandPubKey(t)
	op := test.RandOp(t)
	keyDesc := keychain.KeyDescriptor{
		PubKey: testPubKey,
		KeyLocator: keychain.KeyLocator{
			Family: 123,
			Index:  456,
		},
	}
	inputScriptKey := asset.NewScriptKeyBip86(keyDesc)
	inputScriptKey.Tweak = []byte("merkle root")

	bip32Derivation, trBip32Derivation := Bip32DerivationFromKeyDesc(
		keyDesc, testParams.HDCoinType,
	)
	bip32Derivations := []*psbt.Bip32Derivation{bip32Derivation}
	trBip32Derivations := []*psbt.TaprootBip32Derivation{trBip32Derivation}
	testAsset := asset.RandAsset(t, asset.Normal)
	testAsset.ScriptKey = inputScriptKey

	testOutputAsset := asset.RandAsset(t, asset.Normal)
	testOutputAsset.ScriptKey = asset.NewScriptKeyBip86(keyDesc)

	// The raw key won't be serialized within the asset, so let's blank it
	// out here to get a fully, byte-by-byte comparable PSBT.
	testAsset.GroupKey.RawKey = keychain.KeyDescriptor{}
	testOutputAsset.GroupKey.RawKey = keychain.KeyDescriptor{}
	testOutputAsset.ScriptKey.TweakedScriptKey = nil
	leaf1 := txscript.TapLeaf{
		LeafVersion: txscript.BaseLeafVersion,
		Script:      []byte("not a valid script"),
	}
	testPreimage1 := commitment.NewPreimageFromLeaf(leaf1)
	testPreimage2 := commitment.NewPreimageFromBranch(
		txscript.NewTapBranch(leaf1, leaf1),
	)

	vPacket := &VPacket{
		Inputs: []*VInput{{
			PrevID: asset.PrevID{
				OutPoint:  op,
				ID:        asset.RandID(t),
				ScriptKey: asset.RandSerializedKey(t),
			},
			Anchor: Anchor{
				Value:             777,
				PkScript:          []byte("anchor pkscript"),
				SigHashType:       txscript.SigHashSingle,
				InternalKey:       testPubKey,
				MerkleRoot:        []byte("merkle root"),
				TapscriptSibling:  []byte("sibling"),
				Bip32Derivation:   bip32Derivations,
				TrBip32Derivation: trBip32Derivations,
			},
		}, {
			// Empty input.
		}},
		Outputs: []*VOutput{{
			Amount:                             123,
			IsSplitRoot:                        true,
			Interactive:                        true,
			AnchorOutputIndex:                  0,
			AnchorOutputInternalKey:            testPubKey,
			AnchorOutputBip32Derivation:        bip32Derivations,
			AnchorOutputTaprootBip32Derivation: trBip32Derivations,
			Asset:                              testOutputAsset,
			ScriptKey:                          testOutputAsset.ScriptKey,
			SplitAsset:                         testOutputAsset,
			AnchorOutputTapscriptPreimage:      testPreimage1,
		}, {
			Amount:                             345,
			IsSplitRoot:                        false,
			Interactive:                        false,
			AnchorOutputIndex:                  1,
			AnchorOutputInternalKey:            testPubKey,
			AnchorOutputBip32Derivation:        bip32Derivations,
			AnchorOutputTaprootBip32Derivation: trBip32Derivations,
			Asset:                              testOutputAsset,
			ScriptKey:                          testOutputAsset.ScriptKey,
			AnchorOutputTapscriptPreimage:      testPreimage2,
		}},
		ChainParams: testParams,
	}
	vPacket.SetInputAsset(0, testAsset, []byte("this is a proof"))

	return vPacket
}

// TestEncodeAsPsbt tests the encoding of a virtual packet as a PSBT.
func TestEncodeAsPsbt(t *testing.T) {
	t.Parallel()

	pkg := randomPacket(t)
	packet, err := pkg.EncodeAsPsbt()
	require.NoError(t, err)

	b64, err := packet.B64Encode()
	require.NoError(t, err)

	require.NotEmpty(t, b64)
}