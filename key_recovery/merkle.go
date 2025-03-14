package keyrecovery

import (
	"encoding/hex"
	MT "github.com/txaty/go-merkletree"
)

func GenerateMerkleShares(shares []Share) (string, []MerkleShare, error) {
	root, proofs, err := generateMerkleProofs(shares)
	if err != nil {
		return "", nil, err
	}
	merkleShares := make([]MerkleShare, len(shares))
	for i, share := range shares {
		merkleShares[i] = MerkleShare{Share: share, Proof: *proofs[i]}
	}
	return root, merkleShares, nil
}

func generateMerkleProofs(shares []Share) (string, []*MT.Proof, error) {
	blocks := make([]MT.DataBlock, len(shares))
	for i, share := range shares {
		blocks[i] = &share
	}
	tree, err := MT.New(nil, blocks)
	if err != nil {
		return "", nil, err
	}
	proofs := tree.Proofs
	return hex.EncodeToString(tree.Root), proofs, nil
}

func VerifyMerkleProof(root string, proof *MT.Proof, share Share) (bool, error) {
	rootBytes, err := hex.DecodeString(root)
	if err != nil {
		return false, err
	}
	return MT.Verify(&share, proof, rootBytes, nil)
}

type MerkleShare struct {
	Share Share
	Proof MT.Proof
}
