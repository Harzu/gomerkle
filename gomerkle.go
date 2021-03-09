package gomerkle

import "fmt"

type MerkleTree interface {
	Root() string
	Proof(hash string, layerNum int64) bool
	GetLayer(layerNum int64) []node
}

type merkleTree struct {
	layers map[int64][]node
}

type node struct {
	hash string
	left, right *node
}

func BuildTree(hashes []string) MerkleTree {
	tree := &merkleTree{layers: make(map[int64][]node)}
	if len(hashes) == 0 {
		return tree
	}

	var layerNum int64 = 0
	tree.layers[layerNum] = tree.buildZeroLayer(hashes)

	for len(tree.layers[layerNum]) > 1 {
		nodes := tree.buildLayer(tree.layers[layerNum])
		layerNum++
		tree.layers[layerNum] = nodes
	}

	return tree
}

func (t *merkleTree) buildZeroLayer(hashes []string) []node {
	result := make([]node, 0, len(hashes))
	for _, hash := range hashes {
		result = append(result, node{hash: hash})
	}
	if len(result) % 2 == 1 {
		result = append(result, node{hash: hashes[len(hashes)-1]})
	}
	return result
}

func (t *merkleTree) buildLayer(leavers []node) []node {
	result := make([]node, 0, len(leavers)/2)

	const chunkSize = 2
	nodesFullSize := len(leavers)

	startIndex := 0
	for startIndex < nodesFullSize {
		endIndex := startIndex + chunkSize
		if endIndex >= nodesFullSize {
			endIndex = nodesFullSize
		}

		leafPair := leavers[startIndex:endIndex]
		if len(leafPair) == 1 {
			leafPair = append(leafPair, leafPair[0])
		}

		result = append(result, node{
			hash:  sha3Encode(fmt.Sprintf("%s%s", leafPair[0].hash, leafPair[1].hash)),
			left:  &leafPair[0],
			right: &leafPair[1],
		})

		startIndex = endIndex
	}

	return result
}

func (t *merkleTree) Proof(hash string, layerNum int64) bool {
	nodes, layerExist := t.layers[layerNum]
	if !layerExist {
		return false
	}

	var hashExist bool
	for _, node := range nodes {
		if node.hash == hash {
			hashExist = true
			break
		}
	}

	if hashExist {
		var layerNum int64 = 0
		tree := &merkleTree{layers: map[int64][]node{layerNum: nodes}}

		for len(tree.layers[layerNum]) > 1 {
			nodes := tree.buildLayer(tree.layers[layerNum])
			layerNum++
			tree.layers[layerNum] = nodes
		}

		return t.Root() == tree.Root()
	}

	return false
}

func (t *merkleTree) Root() string {
	rootLayer, exists := t.layers[int64(len(t.layers))-1]
	if !exists || len(rootLayer) != 1 {
		return ""
	}

	return rootLayer[0].hash
}

func (t *merkleTree) GetLayer(layerNum int64) []node {
	return t.layers[layerNum]
}