package godb

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"unsafe"
)

type NodeType uint16

func (n NodeType) String() string {
	switch n {
	case NodeLeaf:
		return "NodeLeaf"
	case NodeInternal:
		return "NodeInternal"
	default:
		return "Unknown"
	}
}

const (
	NodeInternal NodeType = iota + 1
	NodeLeaf
)

// Common Node Header Layout
const (
	NodeTypeSize         = 8 // byte
	NodeTypeOffset       = 0
	IsRootSize           = 8
	IsRootOffset         = NodeTypeSize
	ParentPointerSize    = 32
	ParentPointerOffset  = IsRootOffset + IsRootSize
	CommonNodeHeaderSize = NodeTypeSize + IsRootSize + ParentPointerSize
)

// Leaf Node Header Layout
const (
	LeafNodeNumCellsSize   = 32
	LeafNodeNumCellsOffset = CommonNodeHeaderSize
	LeafNodeHeaderSize     = CommonNodeHeaderSize + LeafNodeNumCellsSize
)

const (
	LeafNodeKeySize       = 32
	LeafNodeKeyOffset     = 0
	LeafNodeValueSize     = 4 + 32 + 256
	LeafNodeValueOffset   = LeafNodeKeyOffset + LeafNodeKeySize
	LeafNodeCellSize      = LeafNodeKeySize + LeafNodeValueSize
	LeafNodeSpaceForCells = PageSize - LeafNodeHeaderSize
	LeafNodeMaxCells      = LeafNodeSpaceForCells / LeafNodeCellSize
)

const (
	maxCellsCount = 4
	// Number of children of a node is equal to the number of keys in it plus 1.
	maxChildren = maxCellsCount + 1
)

type tree struct {
	root  int
	cache map[int]*node
	file  *os.File
}

// GetNode reads node from file if there's no present in cache
func (t tree) GetNode(id int) (*node, error) {
	if n, ok := t.cache[id]; ok {
		return n, nil
	}
	_, err := t.file.Seek(int64(id*NodeSizeInt()), io.SeekStart)
	if err != nil {
		return nil, err
	}

	b := make([]byte, NodeSize())
	_, err = t.file.Read(b)
	if err != nil {
		return nil, err
	}

	n := NewNode()
	if err := UnmarshalNode(b, &n); err != nil {
		return nil, err
	}
	return &n, nil
}

type walkFnPostAction int

const (
	walkPostContinue = iota + 1
	walkPostQuit
)

// walkFn performs some action on node that is being 'walked'
type walkFn func(int, *node) (walkFnPostAction, error)

// walk is a function that walks btree in-order and perform given function
func (t tree) walk(id int, walkFn walkFn) error {
	node, err := t.GetNode(id)
	if err != nil {
		return err
	}
	// TODO: refactor that, it'll go here for all the cases after e.g. node to insert is found
	if action, err := walkFn(id, node); err != nil {
		return err
	} else if action == walkPostQuit {
		return nil
	}

	if !node.hasChildren() {
		return nil
	}
	for _, ch := range node.Children {
		if ch != 0 {
			if err := t.walk(ch, walkFn); err != nil {
				return err
			}
		}
	}
	return nil
}

// Print prints whole tree
func (t tree) Print() error {
	return t.walk(t.root, func(i int, n *node) (walkFnPostAction, error) {
		fmt.Println(n)
		return walkPostContinue, nil
	})
}

// Flush writes whole tree into file starting from root
func (t tree) Flush() error {
	return t.walk(t.root, func(id int, node *node) (walkFnPostAction, error) {
		return walkPostContinue, t.flushNode(id, node)
	})
}

func (t tree) flushNode(id int, n *node) error {
	_, err := t.file.Seek(int64(id*NodeSizeInt()), io.SeekStart)
	if err != nil {
		return err
	}

	b, err := n.Marshal()
	if err != nil {
		return err
	}
	_, err = t.file.Write(b)
	if err != nil {
		return err
	}
	return nil
}

// Add adds row(cell) to the tree
func (t tree) Add(r *Row) error {
	node, err := t.findPlaceToInsert(t.root)
	if err != nil {
		return err
	}
	if node == nil {
		// TODO: add another node here and link it to previous one
		return errors.New("there's no place to insert")
	}
	node.addRow(r)
	return nil
}

// findPlaceToInsert searches btree in order to find first place to insert new Row
//
// it's simple version that searches tree in-order and returns first node that will fit Row
func (t tree) findPlaceToInsert(id int) (*node, error) {
	node, err := t.GetNode(id)
	if err != nil {
		return nil, err
	}

	if node.hasPlaceForRow() {
		return node, nil
	}

	if !node.hasChildren() {
		return nil, nil
	}
	for _, ch := range node.Children {
		if ch != 0 {
			if n, err := t.findPlaceToInsert(ch); err != nil {
				return nil, err
			} else if n != nil {
				return n, nil
			}
		}
	}
	return nil, nil
}

type node struct {
	// Children contains 'pointers' to specific node in bytesystem
	Children [maxChildren]int
	IsRoot   bool
	NodeType NodeType
	// Cells contains Rows that node stores
	Cells map[uint8]*Row
}

func NewNode() node {
	n := node{
		Children: [maxChildren]int{},
		IsRoot:   false,
		NodeType: 0,
		Cells:    make(map[uint8]*Row, maxCellsCount),
	}

	for i := uint8(0); i < maxCellsCount; i++ {
		n.Cells[i] = nil
	}
	return n
}

// NodeSize returns byte size of node
func NodeSize() uint32 {
	return maxChildren*8 + 2 + 1 + maxCellsCount*RowSize()
}

// NodeSizeInt returns byte size of node
func NodeSizeInt() int {
	return int(maxChildren*8 + 2 + 1 + maxCellsCount*RowSize())
}

// Marshal takes node object and parses it to slice of bytes
func (n *node) Marshal() ([]byte, error) {
	b := make([]byte, NodeSize())

	// marshal children ids
	for i := uint32(0); i < maxChildren; i++ {
		binary.LittleEndian.PutUint64(b[uint64Size*(i):uint64Size*(i+1)], uint64(n.Children[i]))
	}
	currentOffset := uint64Size * 5 // size of children
	binary.LittleEndian.PutUint16(b[currentOffset:], uint16(n.NodeType))
	currentOffset += 2
	if n.IsRoot {
		b[currentOffset] = byte(1)
	} else {
		b[currentOffset] = byte(0)
	}
	currentOffset++
	for i := uint32(0); i < maxCellsCount; i++ {
		v := n.Cells[uint8(i)]
		if v == nil {
			copy(b[currentOffset+(i*RowSize()):], make([]byte, RowSize()))
		} else {
			marshaledRow, err := v.Marshal()
			if err != nil {
				return nil, err
			}
			copy(b[currentOffset+(i*RowSize()):], marshaledRow)
		}
	}
	return b, nil
}

var uint64Size = uint32(unsafe.Sizeof(uint64(0)))

// UnmarshalNode takes slice of bytes and parses it to node object
func UnmarshalNode(data []byte, n *node) error {
	chSize := uintptr(8 * maxChildren)
	ntSize := unsafe.Sizeof(n.NodeType)

	ch := data[:chSize] // 0:40
	nt := data[chSize : chSize+ntSize]
	ir := data[chSize+ntSize : chSize+ntSize+1]

	rowsOffset := uint32(chSize + ntSize + 1)

	for i := uint32(1); i <= maxCellsCount; i++ {
		var r Row
		if err := UnmarshalRow(data[rowsOffset+RowSize()*(i-1):rowsOffset+RowSize()*i], &r); err != nil {
			return err
		}
		n.Cells[uint8(i-1)] = &r
	}
	for i := uint32(0); i < maxChildren-1; i++ {
		n.Children[i] = int(binary.LittleEndian.Uint64(ch[uint64Size*i : uint64Size*(i+1)]))
	}

	switch binary.LittleEndian.Uint16(nt) {
	case uint16(NodeInternal):
		n.NodeType = NodeInternal
	case uint16(NodeLeaf):
		n.NodeType = NodeLeaf
	default:
		return errors.New("invalid node type")
	}
	n.IsRoot = ir[0] == byte(1)
	return nil
}

func (n node) String() string {
	return fmt.Sprintf("Children: %v, IsRoot: %v, NodeType: %s, Cells: %v",
		n.Children, n.IsRoot, n.NodeType, n.Cells)
}

func (n node) hasChildren() bool {
	for _, ch := range n.Children {
		if ch != 0 && ch != -1 {
			return true
		}
	}
	return false
}

// addRow adds row to node rows
func (n *node) addRow(r *Row) {
	for k, v := range n.Cells {
		if v == nil {
			n.Cells[k] = r
			return
		}
	}
}

// hasPlaceForRow checks if node has place for another row
func (n node) hasPlaceForRow() bool {
	for _, c := range n.Cells {
		if c == nil {
			return true
		}
	}
	return false
}
