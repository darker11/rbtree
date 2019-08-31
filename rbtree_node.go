package rbtree

type color uint8

const (
	red color = iota
	black
)

// Node defines red-black-tree node structure.
type Node struct {
	key    uint32
	left   *Node
	right  *Node
	parent *Node
	color  color
}

// NewNode used to generate a node.
func NewNode(key uint32, c color) *Node {
	return &Node{
		key:    key,
		left:   nil,
		right:  nil,
		parent: nil,
		color:  c,
	}
}

func (n *Node) setColor(color color) {
	n.color = color
}

func (n *Node) judgeColor(color color) bool {
	return n.color == color
}

func mix(node, sentinel *Node) *Node {
	for node.left != sentinel {
		node = node.left
	}
	return node
}

