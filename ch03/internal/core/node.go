package core

import "sync/atomic"

var globalNodeID uint64

type Node struct {
	id       uint64
	name     string
	children []SceneNode
	parent   SceneNode
}

func NewNode(name string) *Node {
	id := atomic.AddUint64(&globalNodeID, 1)
	return &Node{id: id, name: name, children: make([]SceneNode, 0)}
}

func (n *Node) GetID() uint64        { return n.id }
func (n *Node) GetName() string      { return n.name }
func (n *Node) SetName(name string)  { n.name = name }
func (n *Node) GetParent() SceneNode { return n.parent }

// addChildWithParent attaches child to logicalParent and appends to n.children.
// Used by Node2D so that the stored parent is the Node2D (Transformable), not the embedded Node.
func (n *Node) addChildWithParent(logicalParent SceneNode, child SceneNode) {
	child.AttachParent(logicalParent)
	n.children = append(n.children, child)
}

func (n *Node) AddChildren(child SceneNode) {
	n.addChildWithParent(n, child)
}

func (n *Node) DetachChild(node SceneNode) bool {
	for i, c := range n.children {
		if c == node {
			n.children[i] = n.children[len(n.children)-1]
			n.children = n.children[:len(n.children)-1]
			node.AttachParent(nil)
			return true
		}
	}
	return false
}

func (n *Node) AttachParent(node SceneNode) { n.parent = node }
func (n *Node) GetChildren() []SceneNode   { return n.children }

func (n *Node) MarkDirty() {
	for _, c := range n.children {
		c.MarkDirty()
	}
}
