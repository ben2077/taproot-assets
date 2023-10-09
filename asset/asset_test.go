package asset

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/taproot-assets/fn"
	"github.com/lightninglabs/taproot-assets/internal/test"
	"github.com/lightninglabs/taproot-assets/mssmt"
	"github.com/lightningnetwork/lnd/input"
	"github.com/lightningnetwork/lnd/keychain"
	"github.com/stretchr/testify/require"
)

var (
	hashBytes1     = fn.ToArray[[32]byte](bytes.Repeat([]byte{1}, 32))
	hashBytes2     = fn.ToArray[[32]byte](bytes.Repeat([]byte{2}, 32))
	pubKeyBytes, _ = hex.DecodeString(
		"03a0afeb165f0ec36880b68e0baabd9ad9c62fd1a69aa998bc30e9a34620" +
			"2e078f",
	)
	pubKey, _   = btcec.ParsePubKey(pubKeyBytes)
	sigBytes, _ = hex.DecodeString(
		"e907831f80848d1069a5371b402410364bdf1c5f8307b0084c55f1ce2dca" +
			"821525f66a4a85ea8b71e482a74f382d2ce5ebeee8fdb2172f47" +
			"7df4900d310536c0",
	)
	sig, _     = schnorr.ParseSignature(sigBytes)
	sigWitness = wire.TxWitness{sig.Serialize()}

	generatedTestVectorName = "asset_tlv_encoding_generated.json"

	allTestVectorFiles = []string{
		generatedTestVectorName,
		"asset_tlv_encoding_error_cases.json",
	}

	splitGen = Genesis{
		FirstPrevOut: wire.OutPoint{
			Hash:  hashBytes1,
			Index: 1,
		},
		Tag:         "asset",
		MetaHash:    [MetaHashLen]byte{1, 2, 3},
		OutputIndex: 1,
		Type:        1,
	}
	testSplitAsset = &Asset{
		Version:          1,
		Genesis:          splitGen,
		Amount:           1,
		LockTime:         1337,
		RelativeLockTime: 6,
		PrevWitnesses: []Witness{{
			PrevID: &PrevID{
				OutPoint: wire.OutPoint{
					Hash:  hashBytes1,
					Index: 1,
				},
				ID:        hashBytes1,
				ScriptKey: ToSerialized(pubKey),
			},
			TxWitness:       nil,
			SplitCommitment: nil,
		}},
		SplitCommitmentRoot: nil,
		ScriptVersion:       1,
		ScriptKey:           NewScriptKey(pubKey),
		GroupKey: &GroupKey{
			GroupPubKey: *pubKey,
		},
	}
	testRootAsset = &Asset{
		Version:          1,
		Genesis:          testSplitAsset.Copy().Genesis,
		Amount:           1,
		LockTime:         1337,
		RelativeLockTime: 6,
		PrevWitnesses: []Witness{{
			PrevID: &PrevID{
				OutPoint: wire.OutPoint{
					Hash:  hashBytes2,
					Index: 2,
				},
				ID:        hashBytes2,
				ScriptKey: ToSerialized(pubKey),
			},
			TxWitness:       wire.TxWitness{{2}, {2}},
			SplitCommitment: nil,
		}},
		SplitCommitmentRoot: mssmt.NewComputedNode(hashBytes1, 1337),
		ScriptVersion:       1,
		ScriptKey:           NewScriptKey(pubKey),
		GroupKey: &GroupKey{
			GroupPubKey: *pubKey,
		},
	}
)

// TestGroupKeyIsEqual tests that GroupKey.IsEqual is correct.
func TestGroupKeyIsEqual(t *testing.T) {
	t.Parallel()

	testKey := &GroupKey{
		RawKey: keychain.KeyDescriptor{
			// Fill in some non-defaults.
			KeyLocator: keychain.KeyLocator{
				Family: keychain.KeyFamilyMultiSig,
				Index:  1,
			},
			PubKey: pubKey,
		},
		GroupPubKey: *pubKey,
		Witness:     sigWitness,
	}

	pubKeyCopy := *pubKey

	tests := []struct {
		a, b  *GroupKey
		equal bool
	}{
		{
			a:     nil,
			b:     nil,
			equal: true,
		},
		{
			a:     &GroupKey{},
			b:     &GroupKey{},
			equal: true,
		},
		{
			a:     nil,
			b:     &GroupKey{},
			equal: false,
		},
		{
			a: testKey,
			b: &GroupKey{
				GroupPubKey: *pubKey,
			},
			equal: false,
		},
		{
			a: testKey,
			b: &GroupKey{
				GroupPubKey: testKey.GroupPubKey,
				Witness:     testKey.Witness,
			},
			equal: false,
		},
		{
			a: testKey,
			b: &GroupKey{
				RawKey: keychain.KeyDescriptor{
					KeyLocator: testKey.RawKey.KeyLocator,
					PubKey:     nil,
				},

				GroupPubKey: testKey.GroupPubKey,
				Witness:     testKey.Witness,
			},
			equal: false,
		},
		{
			a: testKey,
			b: &GroupKey{
				RawKey: keychain.KeyDescriptor{
					PubKey: &pubKeyCopy,
				},

				GroupPubKey: testKey.GroupPubKey,
				Witness:     testKey.Witness,
			},
			equal: false,
		},
		{
			a: testKey,
			b: &GroupKey{
				RawKey: keychain.KeyDescriptor{
					KeyLocator: testKey.RawKey.KeyLocator,
					PubKey:     &pubKeyCopy,
				},

				GroupPubKey: testKey.GroupPubKey,
				Witness:     testKey.Witness,
			},
			equal: true,
		},
		{
			a: &GroupKey{
				GroupPubKey: testKey.GroupPubKey,
				Witness:     testKey.Witness,
			},
			b: &GroupKey{
				GroupPubKey: testKey.GroupPubKey,
				Witness:     testKey.Witness,
			},
			equal: true,
		},
		{
			a: &GroupKey{
				RawKey: keychain.KeyDescriptor{
					KeyLocator: testKey.RawKey.KeyLocator,
				},
				GroupPubKey: testKey.GroupPubKey,
				Witness:     testKey.Witness,
			},
			b: &GroupKey{
				RawKey: keychain.KeyDescriptor{
					KeyLocator: testKey.RawKey.KeyLocator,
				},
				GroupPubKey: testKey.GroupPubKey,
				Witness:     testKey.Witness,
			},
			equal: true,
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		require.Equal(t, testCase.equal, testCase.a.IsEqual(testCase.b))
		require.Equal(t, testCase.equal, testCase.b.IsEqual(testCase.a))
	}
}

// TestGenesisAssetClassification tests that the multiple forms of genesis asset
// are recognized correctly.
func TestGenesisAssetClassification(t *testing.T) {
	t.Parallel()

	baseGen := RandGenesis(t, Normal)
	baseScriptKey := RandScriptKey(t)
	baseAsset := RandAssetWithValues(t, baseGen, nil, baseScriptKey)
	assetValidGroup := RandAsset(t, Collectible)
	assetNeedsWitness := baseAsset.Copy()
	assetNeedsWitness.GroupKey = &GroupKey{
		GroupPubKey: *test.RandPubKey(t),
	}
	nonGenAsset := baseAsset.Copy()
	nonGenAsset.PrevWitnesses = []Witness{{
		PrevID: &PrevID{
			OutPoint: wire.OutPoint{
				Hash:  hashBytes1,
				Index: 1,
			},
			ID:        hashBytes1,
			ScriptKey: ToSerialized(pubKey),
		},
		TxWitness:       sigWitness,
		SplitCommitment: nil,
	}}
	groupMemberNonGen := nonGenAsset.Copy()
	groupMemberNonGen.GroupKey = &GroupKey{
		GroupPubKey: *test.RandPubKey(t),
	}
	splitAsset := nonGenAsset.Copy()
	splitAsset.PrevWitnesses[0].TxWitness = nil
	splitAsset.PrevWitnesses[0].SplitCommitment = &SplitCommitment{}

	tests := []struct {
		name                                string
		genAsset                            *Asset
		isGenesis, needsWitness, hasWitness bool
	}{
		{
			name:         "group anchor with witness",
			genAsset:     assetValidGroup,
			isGenesis:    false,
			needsWitness: false,
			hasWitness:   true,
		},
		{
			name:         "ungrouped genesis asset",
			genAsset:     baseAsset,
			isGenesis:    true,
			needsWitness: false,
			hasWitness:   false,
		},
		{
			name:         "group anchor without witness",
			genAsset:     assetNeedsWitness,
			isGenesis:    true,
			needsWitness: true,
			hasWitness:   false,
		},
		{
			name:         "non-genesis asset",
			genAsset:     nonGenAsset,
			isGenesis:    false,
			needsWitness: false,
			hasWitness:   false,
		},
		{
			name:         "non-genesis grouped asset",
			genAsset:     groupMemberNonGen,
			isGenesis:    false,
			needsWitness: false,
			hasWitness:   false,
		},
		{
			name:         "split asset",
			genAsset:     splitAsset,
			isGenesis:    false,
			needsWitness: false,
			hasWitness:   false,
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		a := testCase.genAsset

		hasGenWitness := a.HasGenesisWitness()
		require.Equal(t, testCase.isGenesis, hasGenWitness)
		needsGroupWitness := a.NeedsGenesisWitnessForGroup()
		require.Equal(t, testCase.needsWitness, needsGroupWitness)
		hasGroupWitness := a.HasGenesisWitnessForGroup()
		require.Equal(t, testCase.hasWitness, hasGroupWitness)
	}
}

// TestValidateAssetName tests that asset names are validated correctly.
func TestValidateAssetName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		valid bool
	}{
		{
			// A name with spaces is valid.
			name:  "a name with spaces",
			valid: true,
		},
		{
			// Capital letters are valid.
			name:  "ABC",
			valid: true,
		},
		{
			// Numbers are valid.
			name:  "1234",
			valid: true,
		},
		{
			// A mix of lower/upper, spaces, and numbers is valid.
			name:  "Name 1234",
			valid: true,
		},
		{
			// Japanese characters are valid.
			name:  "日本語",
			valid: true,
		},
		{
			// The "place of interest" character takes up multiple
			// bytes and is valid.
			name:  "⌘",
			valid: true,
		},
		{
			// Exclusively whitespace is an invalid name.
			name:  "   ",
			valid: false,
		},
		{
			// An empty name string is invalid.
			name:  "",
			valid: false,
		},
		{
			// A 65 character name is too long and therefore
			// invalid.
			name: "asdasdasdasdasdasdasdasdasdasdasdasdasdasdas" +
				"dasdasdasdadasdasdada",
			valid: false,
		},
		{
			// Invalid if tab in name.
			name:  "tab\ttab",
			valid: false,
		},
		{
			// Invalid if newline in name.
			name:  "newline\nnewline",
			valid: false,
		},
	}

	for _, testCase := range tests {
		testCase := testCase

		err := ValidateAssetName(testCase.name)
		if testCase.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

// TestAssetEncoding asserts that we can properly encode and decode assets
// through their TLV serialization.
func TestAssetEncoding(t *testing.T) {
	t.Parallel()

	testVectors := &TestVectors{}
	assertAssetEncoding := func(comment string, a *Asset) {
		t.Helper()

		require.True(t, a.DeepEqual(a.Copy()))

		var buf bytes.Buffer
		require.NoError(t, a.Encode(&buf))

		testVectors.ValidTestCases = append(
			testVectors.ValidTestCases, &ValidTestCase{
				Asset:    NewTestFromAsset(t, a),
				Expected: hex.EncodeToString(buf.Bytes()),
				Comment:  comment,
			},
		)

		var b Asset
		require.NoError(t, b.Decode(&buf))

		require.True(t, a.DeepEqual(&b))
	}
	root := testRootAsset.Copy()
	split := testSplitAsset.Copy()
	split.PrevWitnesses[0].SplitCommitment = &SplitCommitment{
		Proof:     *mssmt.RandProof(t),
		RootAsset: *root,
	}
	assertAssetEncoding("random split asset with root asset", split)

	newGen := Genesis{
		FirstPrevOut: wire.OutPoint{
			Hash:  hashBytes2,
			Index: 2,
		},
		Tag:         "asset",
		MetaHash:    [MetaHashLen]byte{1, 2, 3},
		OutputIndex: 2,
		Type:        2,
	}

	comment := "random asset with multiple previous witnesses"
	assertAssetEncoding(comment, &Asset{
		Version:          2,
		Genesis:          newGen,
		Amount:           2,
		LockTime:         1337,
		RelativeLockTime: 6,
		PrevWitnesses: []Witness{{
			PrevID:          nil,
			TxWitness:       nil,
			SplitCommitment: nil,
		}, {
			PrevID:          &PrevID{},
			TxWitness:       nil,
			SplitCommitment: nil,
		}, {
			PrevID: &PrevID{
				OutPoint: wire.OutPoint{
					Hash:  hashBytes2,
					Index: 2,
				},
				ID:        hashBytes2,
				ScriptKey: ToSerialized(pubKey),
			},
			TxWitness:       wire.TxWitness{{2}, {2}},
			SplitCommitment: nil,
		}},
		SplitCommitmentRoot: nil,
		ScriptVersion:       2,
		ScriptKey:           NewScriptKey(pubKey),
		GroupKey:            nil,
	})

	assertAssetEncoding("minimal asset", &Asset{
		Genesis: Genesis{
			MetaHash: [MetaHashLen]byte{},
		},
		ScriptKey: NewScriptKey(pubKey),
	})

	// Write test vectors to file. This is a no-op if the "gen_test_vectors"
	// build tag is not set.
	test.WriteTestVectors(t, generatedTestVectorName, testVectors)
}

// TestAssetIsBurn asserts that the IsBurn method is correct.
func TestAssetIsBurn(t *testing.T) {
	root := testRootAsset.Copy()
	split := testSplitAsset.Copy()
	split.PrevWitnesses[0].SplitCommitment = &SplitCommitment{
		Proof:     *mssmt.RandProof(t),
		RootAsset: *root,
	}

	require.False(t, root.IsBurn())
	require.False(t, split.IsBurn())

	// Update the script key to a burn script key for both of the assets.
	rootPrevID := root.PrevWitnesses[0].PrevID
	root.ScriptKey = NewScriptKey(DeriveBurnKey(*rootPrevID))
	split.ScriptKey = NewScriptKey(DeriveBurnKey(*rootPrevID))

	require.True(t, root.IsBurn())
	require.True(t, split.IsBurn())
}

// TestAssetType asserts that the number of issued assets is set according to
// the genesis type when creating a new asset.
func TestAssetType(t *testing.T) {
	t.Parallel()

	normalGen := Genesis{
		FirstPrevOut: wire.OutPoint{
			Hash:  hashBytes1,
			Index: 1,
		},
		Tag:         "normal asset",
		MetaHash:    [MetaHashLen]byte{1, 2, 3},
		OutputIndex: 1,
		Type:        Normal,
	}
	collectibleGen := Genesis{
		FirstPrevOut: wire.OutPoint{
			Hash:  hashBytes1,
			Index: 1,
		},
		Tag:         "collectible asset",
		MetaHash:    [MetaHashLen]byte{1, 2, 3},
		OutputIndex: 2,
		Type:        Collectible,
	}
	scriptKey := NewScriptKey(pubKey)

	normal, err := New(normalGen, 741, 0, 0, scriptKey, nil)
	require.NoError(t, err)
	require.EqualValues(t, 741, normal.Amount)

	_, err = New(collectibleGen, 741, 0, 0, scriptKey, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "amount must be 1 for asset")

	collectible, err := New(collectibleGen, 1, 0, 0, scriptKey, nil)
	require.NoError(t, err)
	require.EqualValues(t, 1, collectible.Amount)
}

// TestAssetID makes sure that the asset ID is derived correctly.
func TestAssetID(t *testing.T) {
	t.Parallel()

	g := Genesis{
		FirstPrevOut: wire.OutPoint{
			Hash:  hashBytes1,
			Index: 99,
		},
		Tag:         "collectible asset 1",
		MetaHash:    [MetaHashLen]byte{1, 2, 3},
		OutputIndex: 21,
		Type:        Collectible,
	}
	tagHash := sha256.Sum256([]byte(g.Tag))

	h := sha256.New()
	_ = wire.WriteOutPoint(h, 0, 0, &g.FirstPrevOut)
	_, _ = h.Write(tagHash[:])
	_, _ = h.Write(g.MetaHash[:])
	_, _ = h.Write([]byte{0, 0, 0, 21, 1})
	result := h.Sum(nil)

	id := g.ID()
	require.Equal(t, result, id[:])

	// Make sure we get a different asset ID even if everything is the same
	// except for the type.
	normalWithDifferentType := Genesis{
		FirstPrevOut: wire.OutPoint{
			Hash:  hashBytes1,
			Index: 99,
		},
		Tag:         "collectible asset 1",
		MetaHash:    [MetaHashLen]byte{1, 2, 3},
		OutputIndex: 21,
		Type:        Normal,
	}
	differentID := normalWithDifferentType.ID()
	require.NotEqual(t, id[:], differentID[:])
}

// TestAssetGroupKey tests that the asset key group is derived correctly.
func TestAssetGroupKey(t *testing.T) {
	t.Parallel()

	privKey, err := btcec.NewPrivateKey()
	groupPub := privKey.PubKey()
	require.NoError(t, err)
	privKeyCopy := btcec.PrivKeyFromScalar(&privKey.Key)
	genSigner := NewMockGenesisSigner(privKeyCopy)
	genBuilder := MockGroupTxBuilder{}
	fakeKeyDesc := test.PubToKeyDesc(groupPub)
	fakeScriptKey := NewScriptKeyBip86(fakeKeyDesc)

	g := Genesis{
		FirstPrevOut: wire.OutPoint{
			Hash:  hashBytes1,
			Index: 99,
		},
		Tag:         "normal asset 1",
		MetaHash:    [MetaHashLen]byte{1, 2, 3},
		OutputIndex: 21,
		Type:        Collectible,
	}
	groupTweak := g.ID()

	internalKey := input.TweakPrivKey(privKeyCopy, groupTweak[:])
	tweakedKey := txscript.TweakTaprootPrivKey(*internalKey, nil)

	// TweakTaprootPrivKey modifies the private key that is passed in! We
	// need to provide a copy to arrive at the same result.
	protoAsset := NewAssetNoErr(t, g, 1, 0, 0, fakeScriptKey, nil)
	keyGroup, err := DeriveGroupKey(
		genSigner, &genBuilder, fakeKeyDesc, g, protoAsset,
	)
	require.NoError(t, err)

	require.Equal(
		t, schnorr.SerializePubKey(tweakedKey.PubKey()),
		schnorr.SerializePubKey(&keyGroup.GroupPubKey),
	)

	// Group key tweaking should fail when given invalid tweaks.
	badTweak := test.RandBytes(33)
	_, err = GroupPubKey(groupPub, badTweak, badTweak)
	require.Error(t, err)

	_, err = GroupPubKey(groupPub, groupTweak[:], badTweak)
	require.Error(t, err)
}

// TestDeriveGroupKey tests that group key derivation fails for assets that are
// not eligible to be group anchors.
func TestDeriveGroupKey(t *testing.T) {
	t.Parallel()

	groupPriv := test.RandPrivKey(t)
	groupPub := groupPriv.PubKey()
	groupKeyDesc := test.PubToKeyDesc(groupPub)
	genSigner := NewMockGenesisSigner(groupPriv)
	genBuilder := MockGroupTxBuilder{}

	baseGen := RandGenesis(t, Normal)
	collectGen := RandGenesis(t, Collectible)
	baseScriptKey := RandScriptKey(t)
	protoAsset := RandAssetWithValues(t, baseGen, nil, baseScriptKey)
	nonGenProtoAsset := protoAsset.Copy()
	nonGenProtoAsset.PrevWitnesses = []Witness{{
		PrevID: &PrevID{
			OutPoint: wire.OutPoint{
				Hash:  hashBytes1,
				Index: 1,
			},
			ID:        hashBytes1,
			ScriptKey: ToSerialized(pubKey),
		},
		TxWitness:       sigWitness,
		SplitCommitment: nil,
	}}
	groupedProtoAsset := protoAsset.Copy()
	groupedProtoAsset.GroupKey = &GroupKey{
		GroupPubKey: *groupPub,
	}

	// A prototype asset is required for building the genesis virtual TX.
	_, err := DeriveGroupKey(
		genSigner, &genBuilder, groupKeyDesc, baseGen, nil,
	)
	require.Error(t, err)

	// The prototype asset must have a genesis witness.
	_, err = DeriveGroupKey(
		genSigner, &genBuilder, groupKeyDesc, baseGen, nonGenProtoAsset,
	)
	require.Error(t, err)

	// The prototype asset must not have a group key set.
	_, err = DeriveGroupKey(
		genSigner, &genBuilder, groupKeyDesc, baseGen, groupedProtoAsset,
	)
	require.Error(t, err)

	// The anchor genesis used for signing must have the same asset type
	// as the prototype asset being signed.
	_, err = DeriveGroupKey(
		genSigner, &genBuilder, groupKeyDesc, collectGen, protoAsset,
	)
	require.Error(t, err)

	groupKey, err := DeriveGroupKey(
		genSigner, &genBuilder, groupKeyDesc, baseGen, protoAsset,
	)
	require.NoError(t, err)
	require.NotNil(t, groupKey)
}

// TestAssetWitness tests that the asset group witness can be serialized and
// parsed correctly, and that signature detection works correctly.
func TestAssetWitnesses(t *testing.T) {
	t.Parallel()

	nonSigWitness := test.RandTxWitnesses(t)
	for len(nonSigWitness) == 0 {
		nonSigWitness = test.RandTxWitnesses(t)
	}

	// A witness must be unmodified after serialization and parsing.
	nonSigWitnessBytes, err := SerializeGroupWitness(nonSigWitness)
	require.NoError(t, err)

	nonSigWitnessParsed, err := ParseGroupWitness(nonSigWitnessBytes)
	require.NoError(t, err)
	require.Equal(t, nonSigWitness, nonSigWitnessParsed)

	// A witness that is a single Schnorr signature must be detected
	// correctly both before and after serialization.
	sigWitnessParsed, isSig := IsGroupSig(sigWitness)
	require.True(t, isSig)
	require.NotNil(t, sigWitnessParsed)

	sigWitnessBytes, err := SerializeGroupWitness(sigWitness)
	require.NoError(t, err)

	sigWitnessParsed, err = ParseGroupSig(sigWitnessBytes)
	require.NoError(t, err)
	require.Equal(t, sig.Serialize(), sigWitnessParsed.Serialize())

	// Adding an annex to the witness stack should not affect signature
	// parsing.
	dummyAnnex := []byte{0x50, 0xde, 0xad, 0xbe, 0xef}
	sigWithAnnex := wire.TxWitness{sigWitness[0], dummyAnnex}
	sigWitnessParsed, isSig = IsGroupSig(sigWithAnnex)
	require.True(t, isSig)
	require.NotNil(t, sigWitnessParsed)

	// Witness that are not a single Schnorr signature must also be
	// detected correctly.
	possibleSig, isSig := IsGroupSig(nonSigWitness)
	require.False(t, isSig)
	require.Nil(t, possibleSig)

	possibleSig, err = ParseGroupSig(nonSigWitnessBytes)
	require.Error(t, err)
	require.Nil(t, possibleSig)
}

// TestUnknownVersion tests that an asset of an unknown version is rejected
// before being inserted into an MS-SMT.
func TestUnknownVersion(t *testing.T) {
	t.Parallel()

	rootGen := Genesis{
		FirstPrevOut: wire.OutPoint{
			Hash:  hashBytes1,
			Index: 1,
		},
		Tag:         "asset",
		MetaHash:    [MetaHashLen]byte{1, 2, 3},
		OutputIndex: 1,
		Type:        1,
	}

	root := &Asset{
		Version:          212,
		Genesis:          rootGen,
		Amount:           1,
		LockTime:         1337,
		RelativeLockTime: 6,
		PrevWitnesses: []Witness{{
			PrevID: &PrevID{
				OutPoint: wire.OutPoint{
					Hash:  hashBytes2,
					Index: 2,
				},
				ID:        hashBytes2,
				ScriptKey: ToSerialized(pubKey),
			},
			TxWitness:       wire.TxWitness{{2}, {2}},
			SplitCommitment: nil,
		}},
		SplitCommitmentRoot: mssmt.NewComputedNode(hashBytes1, 1337),
		ScriptVersion:       1,
		ScriptKey:           NewScriptKey(pubKey),
		GroupKey: &GroupKey{
			GroupPubKey: *pubKey,
			Witness:     sigWitness,
		},
	}

	rootLeaf, err := root.Leaf()
	require.Nil(t, rootLeaf)
	require.ErrorIs(t, err, ErrUnknownVersion)

	root.Version = V0
	rootLeaf, err = root.Leaf()
	require.NotNil(t, rootLeaf)
	require.Nil(t, err)
}

func FuzzAssetDecode(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		r := bytes.NewReader(data)
		a := &Asset{}
		if err := a.Decode(r); err != nil {
			return
		}
	})
}

// TestBIPTestVectors tests that the BIP test vectors are passing.
func TestBIPTestVectors(t *testing.T) {
	t.Parallel()

	for idx := range allTestVectorFiles {
		var (
			fileName    = allTestVectorFiles[idx]
			testVectors = &TestVectors{}
		)
		test.ParseTestVectors(t, fileName, &testVectors)
		t.Run(fileName, func(tt *testing.T) {
			tt.Parallel()

			runBIPTestVector(tt, testVectors)
		})
	}
}

// runBIPTestVector runs the tests in a single BIP test vector file.
func runBIPTestVector(t *testing.T, testVectors *TestVectors) {
	for _, validCase := range testVectors.ValidTestCases {
		validCase := validCase

		t.Run(validCase.Comment, func(tt *testing.T) {
			tt.Parallel()

			a := validCase.Asset.ToAsset(tt)

			var buf bytes.Buffer
			err := a.Encode(&buf)
			require.NoError(tt, err)

			areEqual := validCase.Expected == hex.EncodeToString(
				buf.Bytes(),
			)

			// Create nice diff if things don't match.
			if !areEqual {
				expectedBytes, err := hex.DecodeString(
					validCase.Expected,
				)
				require.NoError(tt, err)

				expectedAsset := &Asset{}
				err = expectedAsset.Decode(bytes.NewReader(
					expectedBytes,
				))
				require.NoError(tt, err)

				require.Equal(tt, a, expectedAsset)

				// Make sure we still fail the test.
				require.Equal(
					tt, validCase.Expected,
					hex.EncodeToString(buf.Bytes()),
				)
			}

			// We also want to make sure that the asset is decoded
			// correctly from the encoded TLV stream.
			decoded := &Asset{}
			err = decoded.Decode(&buf)
			require.NoError(tt, err)

			require.Equal(tt, a, decoded)
		})
	}

	for _, invalidCase := range testVectors.ErrorTestCases {
		invalidCase := invalidCase

		t.Run(invalidCase.Comment, func(tt *testing.T) {
			tt.Parallel()

			require.PanicsWithValue(t, invalidCase.Error, func() {
				invalidCase.Asset.ToAsset(tt)
			})
		})
	}
}

// TestAssetEncodingNoWitness tests that we can properly encode and decode an
// asset using the v1 version where the witness is not included.
func TestAssetEncodingNoWitness(t *testing.T) {
	t.Parallel()

	// First, start by copying the root asset re-used across tests.
	root := testRootAsset.Copy()

	// We'll make another copy that we'll use to modify the witness field.
	root2 := root.Copy()

	// We'll now modify the witness field of the second root.
	root2.PrevWitnesses[0].TxWitness[0][0] ^= 1

	// If we encode both of these assets then, then final encoding should
	// be identical as we use the EncodeNoWitness method.
	var b1, b2 bytes.Buffer
	require.NoError(t, root.EncodeNoWitness(&b1))
	require.NoError(t, root2.EncodeNoWitness(&b2))

	require.Equal(t, b1.Bytes(), b2.Bytes())

	// The leaf encoding for these two should also be identical.
	root1Leaf, err := root.Leaf()
	require.NoError(t, err)
	root2Leaf, err := root2.Leaf()
	require.NoError(t, err)

	require.Equal(t, root1Leaf.NodeHash(), root2Leaf.NodeHash())
}

// TestNewAssetWithCustomVersion tests that a custom version can be set for
// newly created assets.
func TestNewAssetWithCustomVersion(t *testing.T) {
	t.Parallel()

	// We'll use the root asset as a template, to re-use some of its static
	// data.
	rootAsset := testRootAsset.Copy()

	const newVersion = 10

	assetCustomVersion, err := New(
		rootAsset.Genesis, rootAsset.Amount, 0, 0, rootAsset.ScriptKey, nil,
		WithAssetVersion(newVersion),
	)
	require.NoError(t, err)

	require.Equal(t, int(assetCustomVersion.Version), newVersion)
}
