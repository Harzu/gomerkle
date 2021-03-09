package gomerkle

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMerkleTree_Proof(t *testing.T) {
	req := require.New(t)

	cases := map[string]struct{
		data  []string
		check string
		want  bool
	}{
		"success": {
			data: encodeDataSliceToSha3([]string{"a", "b", "c", "d"}),
			check: sha3Encode("c"),
			want:  true,
		},
		"not even": {
			data:  encodeDataSliceToSha3([]string{"a", "b", "c"}),
			check: sha3Encode("c"),
			want:  true,
		},
		"one element": {
			data:  encodeDataSliceToSha3([]string{"a"}),
			check: sha3Encode("a"),
			want:  true,
		},
		"proof false": {
			data:  encodeDataSliceToSha3([]string{"a", "b", "c", "d"}),
			check: sha3Encode("w"),
			want:  false,
		},
		"empty": {
			data:  []string{},
			check: sha3Encode("w"),
			want:  false,
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			tree := BuildTree(cs.data)
			proof := tree.Proof(cs.check, 0)
			req.Equal(cs.want, proof)
		})
	}
}

func TestMerkleTree_ProofLayers(t *testing.T) {
	req := require.New(t)

	cases := map[string]struct{
		data              []string
		layer             int64
		layerNodePosition int
		want              bool
		isFail            bool
		failNodeHash      string
	}{
		"success": {
			data:              encodeDataSliceToSha3([]string{"a", "b", "c", "d", "e"}),
			layer:             2,
			layerNodePosition: 1,
			want:              true,
		},
		"fail": {
			data:              encodeDataSliceToSha3([]string{"a", "b", "c", "d", "e"}),
			layer:             3,
			layerNodePosition: 0,
			want:              false,
			isFail:            true,
			failNodeHash:      sha3Encode("test"),
		},
		"0 layer": {
			data:              encodeDataSliceToSha3([]string{"a", "b", "c", "d", "e"}),
			layer:             0,
			layerNodePosition: 1,
			want:              true,
		},
		"root": {
			data:              encodeDataSliceToSha3([]string{"a", "b", "c", "d", "e"}),
			layer:             3,
			layerNodePosition: 0,
			want:              true,
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			tree := BuildTree(cs.data)

			var proof bool
			if cs.isFail {
				proof = tree.Proof(cs.failNodeHash, cs.layer)
			} else {
				nodes := tree.GetLayer(cs.layer)
				proof = tree.Proof(nodes[cs.layerNodePosition].hash, cs.layer)
			}

			req.Equal(cs.want, proof)
		})
	}
}

func encodeDataSliceToSha3(data []string) []string {
	for index := range data {
		data[index] = sha3Encode(data[index])
	}
	return data
}