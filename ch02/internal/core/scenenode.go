package core

// SceneNode defines the interface for an element in the scene graph.
// Every node in the tree implements this interface.
type SceneNode interface {
	GetID() uint64
	GetName() string
	AddChildren(child SceneNode)
	GetChildren() []SceneNode
	AttachParent(node SceneNode)
	GetParent() SceneNode
	DetachChild(node SceneNode) bool
	MarkDirty()
}
