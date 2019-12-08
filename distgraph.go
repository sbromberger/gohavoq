package gohavoq

import "fmt"

type vHash = uint64

func partFunc(v vHash, k int) int {
	return int(v % vHash(k))
}

type vertexMap map[vHash][]vHash

// GraphNode represents one part of a distributed graph.
type GraphNode struct {
	nodeID, nNodes int
	vMap           vertexMap
}

func (d *GraphNode) String() string {
	return fmt.Sprintf("Distributed graph part %d of %d with %d vertices", d.nodeID, d.nNodes, len(d.vMap))
}

// IsLocal returns true if the supplied vertex is local to this node.
func (d *GraphNode) IsLocal(v vHash) bool {
	_, found := d.vMap[v]
	return found
}

// GetNodeFor returns the remote node associated with a given hash by calling partFunc().
func (d *GraphNode) GetNodeFor(v vHash) int {
	return partFunc(v, d.nNodes)
}

// LoadEdgeList loads an edgelist from an mmapped file.
// func LoadEdgeList(fn string) EdgeList {

// func LoadNode(nodeID, nNodes int, fn string) GraphNode {
// 	el := LoadEdgeList(fn)
// 	return el.ToNode(nodeID, nNodes)
// }
