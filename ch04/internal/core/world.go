package core

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

const DefaultLayerIndex = 0

// World manages the scene graph with layers.
type World struct {
	rootScene  SceneNode
	layerRoots []SceneNode
	drawLayers *Layers
	drawTarget *ebiten.Image // set each Draw call, used during queueNode
}

func NewWorld(w, h uint) *World {
	return &World{
		rootScene:  NewNode("root"),
		layerRoots: make([]SceneNode, 0),
		drawLayers: NewLayers(),
	}
}

func (w *World) AddNodeToLayer(node SceneNode, layerIndex int) {
	if layerIndex < 0 {
		return
	}
	for layerIndex >= len(w.layerRoots) {
		root := NewNode("layer_root")
		w.layerRoots = append(w.layerRoots, root)
		w.rootScene.AddChildren(root)
	}
	w.layerRoots[layerIndex].AddChildren(node)
}

func (w *World) AddNodeToDefaultLayer(node SceneNode) {
	w.AddNodeToLayer(node, DefaultLayerIndex)
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
	for _, child := range node.GetChildren() {
		w.updateNode(child)
	}
}

func (w *World) queueNode(node SceneNode, layerIndex int) {
	if node == nil {
		return
	}

	children := node.GetChildren()
	sort.Slice(children, func(i, j int) bool {
		li, lj := 0, 0
		if d, ok := children[i].(Drawable); ok {
			li = d.GetLayer()
		}
		if d, ok := children[j].(Drawable); ok {
			lj = d.GetLayer()
		}
		return li > lj // higher layer first = pushed first = drawn first (behind)
	})

	for _, child := range children {
		w.queueNode(child, layerIndex)
	}

	if drawable, ok := node.(Drawable); ok {
		op := ebiten.DrawImageOptions{}
		if tr, ok := node.(Transformable); ok {
			op.GeoM = buildGeoMFromTransform(tr.GetWorldTransform())
		}
		w.drawLayers.AddNode(layerIndex, drawable, w.drawTarget, op)
	}
}

func (w *World) Draw(target *ebiten.Image) {
	w.drawTarget = target
	for i := range w.layerRoots {
		w.queueNode(w.layerRoots[i], i)
	}
	w.drawLayers.DrawAll()
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
