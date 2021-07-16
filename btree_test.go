package godb

import (
	"fmt"
	"os"
	"testing"
)

func generateNode(s string) *node {
	n := NewNode()
	n.IsRoot = true
	n.NodeType = NodeInternal
	n.Cells[0] = prepRowWithValues(0, s, s)
	return &n
}

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

func Test_tree_Print(t *testing.T) {
	file, err := os.CreateTemp("", "test")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	root := generateNode("root") // 1 - root
	root.Children[0] = 1
	root.Children[1] = 2
	n1 := generateNode("n1") // 2
	n2 := generateNode("n2") // 3

	tr := tree{
		cache: make(map[int]*node),
		file:  file,
		root:  0,
	}
	tr.cache[0] = root
	tr.cache[1] = n1
	tr.cache[2] = n2

	tr.Print()

	if err := tr.Flush(); err != nil {
		t.Fatal(err)
	}
	fmt.Println("======================")

	tr2 := tree{file: file}
	if err := tr2.Print(); err != nil {
		t.Fatal(err)
	}
}

func Test_tree_PrintOneNode(t *testing.T) {
	file, err := os.CreateTemp("", "testFlushOneNode")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	n1 := generateNode("n1") // 1 - root
	n1.IsRoot = false
	n1.NodeType = NodeLeaf

	tr := tree{
		cache: make(map[int]*node),
		file:  file,
		root:  0,
	}
	tr.cache[0] = n1

	tr.Print()

	if err := tr.Flush(); err != nil {
		t.Fatal(err)
	}

	tr2 := tree{file: file}
	if err := tr2.Print(); err != nil {
		t.Fatal(err)
	}
}

func Test_tree_PrintTwoNodes(t *testing.T) {
	file, err := os.CreateTemp("", "testFush2")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	root := generateNode("root") // 1 - root
	root.IsRoot = false
	root.NodeType = NodeLeaf
	root.Children[0] = 1

	n1 := generateNode("n1")

	tr := tree{
		cache: make(map[int]*node),
		file:  file,
		root:  0,
	}
	tr.cache[0] = root
	tr.cache[1] = n1
	tr.Print()

	if err := tr.Flush(); err != nil {
		t.Fatal(err)
	}
	fmt.Println("======================")

	tr2 := tree{file: file}
	if err := tr2.Print(); err != nil {
		t.Fatal(err)
	}
}

func Test_tree_Add(t *testing.T) {
	root := generateNode("root") // 1 - root
	root.IsRoot = true
	root.NodeType = NodeLeaf
	root.Cells[0] = prepRowWithValues(0, "username0", "email0")
	root.Cells[1] = prepRowWithValues(1, "username1", "email1")
	root.Cells[2] = prepRowWithValues(1, "username1", "email1")
	tr := tree{
		cache: map[int]*node{
			0: root,
		},
		root: 0,
	}

	err := tr.Add(prepRowWithValues(4, "qweqwe", "poqwkeqpowke"))
	if err != nil {
		t.Fatal(err)
	}

	err = tr.Print()
	if err != nil {
		t.Fatal(err)
	}
	err = tr.Add(prepRowWithValues(4, "username4", "email4"))
	if err == nil {
		t.Fatal("this error should be nil")
	}
}
