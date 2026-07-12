package core

type Node2D struct {
	Node
	localTransform Transform
	worldTransform Transform
	isDirty        bool
}

func NewNode2D(name string) *Node2D {
	return &Node2D{
		Node:          *NewNode(name),
		localTransform: NewTransform(ZeroVector2D(), ZeroVector2D(), 0),
		worldTransform: NewTransform(ZeroVector2D(), ZeroVector2D(), 0),
		isDirty:        true,
	}
}

func (n *Node2D) GetTransform() Transform { return n.localTransform }
func (n *Node2D) SetTransform(t Transform) { n.localTransform = t; n.MarkDirty() }
func (n *Node2D) SetPosition(x, y float64) { n.localTransform.SetPosition(NewVector2D(x, y)); n.MarkDirty() }
func (n *Node2D) GetPosition() Vector2D   { return n.localTransform.GetPosition() }
func (n *Node2D) SetRotation(r float64)   { n.localTransform.SetRotation(r); n.MarkDirty() }
func (n *Node2D) GetRotation() float64    { return n.localTransform.GetRotation() }
func (n *Node2D) SetScale(x, y float64)   { n.localTransform.SetScale(x, y); n.MarkDirty() }
func (n *Node2D) GetScale() Vector2D     { return n.localTransform.GetScale() }
func (n *Node2D) GetPivot() Vector2D     { return n.localTransform.GetPivot() }
func (n *Node2D) SetPivot(x, y float64)  { n.localTransform.SetPivot(x, y); n.MarkDirty() }

// AddChildren overrides Node.AddChildren so the stored parent is the Node2D (Transformable),
// not the embedded Node. Delegates to Node.addChildWithParent to avoid code duplication.
func (n *Node2D) AddChildren(child SceneNode) {
	n.Node.addChildWithParent(n, child)
}

func (n *Node2D) MarkDirty() {
	if n.isDirty {
		return
	}
	n.isDirty = true
	n.Node.MarkDirty()
}

func (n *Node2D) GetWorldTransform() Transform {
	if !n.isDirty {
		return n.worldTransform
	}
	world := NewTransform(ZeroVector2D(), ZeroVector2D(), 0)
	if parent := n.GetParent(); parent != nil {
		if pt, ok := parent.(Transformable); ok {
			world = pt.GetWorldTransform()
		}
	}
	world.Concat(n.localTransform)
	n.worldTransform = world
	n.isDirty = false
	return n.worldTransform
}

// GetWorldPosition returns the world-space position (convenience for collision etc.).
func (n *Node2D) GetWorldPosition() Vector2D {
	wt := n.GetWorldTransform()
	return wt.GetPosition()
}
