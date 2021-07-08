package godb

type NodeType uint8

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

type node struct {
	IsRoot   bool
	NodeType NodeType
	Cells    map[uint32]*Row
}
