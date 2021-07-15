package godb

import (
	"fmt"
	"os"
	"testing"
)

func Test_node_Marshal(t *testing.T) {
	file, err := os.OpenFile("test", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		t.Fatal(err)
	}

	n := NewNode()
	n.Children[0] = 123
	n.Children[3] = 1231254
	n.IsRoot = true
	n.NodeType = NodeInternal
	n.Cells[0] = prepRowWithValues(0, "username0", "email0")
	n.Cells[1] = prepRowWithValues(1, "username1", "email1")
	fmt.Println(n)

	b, err := n.Marshal()
	file.Write(b)

	var n1 = NewNode()
	if err := UnmarshalNode(b, &n1); err != nil {
		t.Fatal(err)
	}

	fmt.Println(n1)

	var n2 = NewNode()
	readBytes := make([]byte, NodeSize())
	_, err = file.ReadAt(readBytes, int64(NodeSize()))
	if err != nil {
		t.Fatal(err)
	}

	if err := UnmarshalNode(b, &n2); err != nil {
		t.Fatal(err)
	}
	fmt.Println(n2)
}
