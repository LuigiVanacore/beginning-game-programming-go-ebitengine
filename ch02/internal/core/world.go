package core

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// World manages the scene graph and draws it to the screen.
type World struct {
	rootScene SceneNode
}

func NewWorld() *World {
	return &World{
		rootScene: NewNode("root"),
	}
}

func (w *World) AddNodeToDefaultLayer(node SceneNode) {
	w.rootScene.AddChildren(node)
}

func (w *World) RemoveNode(node SceneNode) bool {
	parent := node.GetParent()
	if parent == nil {
		return false
	}
	if !parent.DetachChild(node) {
		return false
	}
	node.AttachParent(nil)
	return true
}

func (w *World) Update() {
	w.updateNode(w.rootScene)
}

func (w *World) updateNode(node SceneNode) {
	if node == nil {
		return
	}
	// Updatable support can be added later
	for _, child := range node.GetChildren() {
		w.updateNode(child)
	}
}

func (w *World) Draw(target *ebiten.Image) {
	w.drawNode(w.rootScene, ebiten.GeoM{}, target)
}

func (w *World) drawNode(node SceneNode, parentGeoM ebiten.GeoM, target *ebiten.Image) {
	if node == nil {
		return
	}

	// Draw this node first if it's a Drawable (parents before children = parents behind)
	if drawable, ok := node.(Drawable); ok {
		op := &ebiten.DrawImageOptions{}
		op.GeoM = buildGeoMFromTransform(drawable.GetWorldTransform())
		drawable.Draw(target, op)
	}

	// Recurse into children
	for _, child := range node.GetChildren() {
		w.drawNode(child, ebiten.GeoM{}, target)
	}
}

func buildGeoMFromTransform(t Transform) ebiten.GeoM {
	g := ebiten.GeoM{}
	pivot := t.GetPivot()
	pos := t.GetPosition()
	scale := t.GetScale()
	rot := t.GetRotation()

	// Move pivot to origin -> Scale -> Rotate -> Place at world position
	g.Translate(-pivot.X(), -pivot.Y())
	g.Scale(scale.X(), scale.Y())
	g.Rotate(rot)
	g.Translate(pos.X(), pos.Y())
	return g
}
