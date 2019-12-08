package gohavoq

import (
	"os"
	"unsafe"

	mmap "github.com/edsrzf/mmap-go"
)

type Edge struct {
	Src, Dst uint64
}

type EdgeList []Edge

func (el EdgeList) toVMap() vertexMap {
	vMap := make(vertexMap)
	for _, e := range el {
		s := uint64(e.Src)
		d := uint64(e.Dst)
		vMap[s] = append(vMap[s], d)
	}
	return vMap
}

func (el EdgeList) ToNode(nodeID, nNodes int) GraphNode {
	vMap := el.toVMap()
	return GraphNode{nodeID: nodeID, nNodes: nNodes, vMap: vMap}
}

// raw defines a struct for binary EdgeList format.
type raw struct {
	file *os.File
	data mmap.MMap

	nEdges uint64
	eList  EdgeList // interleaved uint64 src/dst
}

func (raw *raw) close() error {
	var err1, err2 error
	if raw.data != nil {
		err1 = raw.data.Unmap()
	}
	if raw.file != nil {
		err2 = raw.file.Close()
	}
	if err1 != nil {
		return err1
	}
	return err2
}

func loadRaw(filename string) (*raw, error) {
	var err error
	raw := &raw{}

	raw.file, err = os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		raw.close()
		return nil, err
	}

	raw.data, err = mmap.Map(raw.file, mmap.RDONLY, 0)
	if err != nil {
		raw.close()
		return nil, err
	}

	x := 0
	copy((*[8]byte)(unsafe.Pointer(&raw.nEdges))[:], raw.data[x:x+8])
	x += 8

	raw.eList = ((*[1 << 40]Edge)(unsafe.Pointer(&raw.data[x])))[0:int(raw.nEdges)]
	return raw, nil
}

// Save saves an EdgeList to a file in raw (binary) format.
func (el EdgeList) Save(filename string) error {
	output, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer output.Close()

	elLen := int64(len(el))

	elBytes := 16 * len(el)

	err = output.Truncate(int64(8 + elBytes))
	if err != nil {
		return err
	}

	data, err := mmap.Map(output, mmap.RDWR, 0)
	if err != nil {
		return err
	}
	defer data.Unmap()

	x := 0

	copy(data[x:x+8], ((*[8]byte)(unsafe.Pointer(&elLen))[:]))
	x += 8

	if len(el) > 0 {
		copy(data[x:x+elBytes],
			((*[1 << 40]byte)(unsafe.Pointer(&el[0]))[:elBytes]))
	}

	return nil
}

// Load loads a raw (binary) SimpleGraph file and returns am EdgeList.
func Load(fn string) (EdgeList, error) {
	raw, err := loadRaw(fn)
	if err != nil {
		return EdgeList{}, err
	}

	return raw.eList, nil

}
