package core

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
