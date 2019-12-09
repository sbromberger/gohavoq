package gohavoq

//GoHavoq implements HAVOQ-GT in Go.

import "fmt"

type vertexMap map[uint64][]uint64

// GraphNode represents one part of a distributed graph.
type GraphNode struct {
	nodeID, nNodes int
	vMap           vertexMap
}

func (d *GraphNode) String() string {
	return fmt.Sprintf("Distributed graph part %d of %d with %d vertices", d.nodeID, d.nNodes, len(d.vMap))
}

// IsLocal returns true if the supplied vertex is local to this node.
func (d *GraphNode) IsLocal(v uint64) bool {
	_, found := d.vMap[v]
	return found
}

func partFunc(v uint64, k int) int {
	return int(v % uint64(k))
}

// GetNodeFor returns the remote node associated with a given hash by calling partFunc().
func (d *GraphNode) GetNodeFor(v uint64) int {
	return partFunc(v, d.nNodes)
}

func makePartFn(fn string, n int) string {
	return fmt.Sprintf("%s-%d", fn, n)
}

// LoadNode creates a node from a nodeID, total number of nodes, and a file prefix that contains partitions.
func LoadNode(nodeID, nNodes int, fn string) (GraphNode, error) {
	partFn := makePartFn(fn, nodeID)
	el, err := Load(partFn)
	if err != nil {
		return GraphNode{}, err
	}

	return el.ToNode(nodeID, nNodes), nil
}
