package core

import (
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

const DefaultLayerIndex = 0

// World manages the scene graph with layers and a camera.
type World struct {
	rootScene  SceneNode
	layerRoots []SceneNode
	drawLayers *Layers
	camera     *Camera
	drawTarget *ebiten.Image
}

func NewWorld(w, h uint) *World {
	return &World{
		rootScene:  NewNode(NameRoot),
		layerRoots: make([]SceneNode, 0),
		drawLayers: NewLayers(),
		camera:     NewCamera(w, h),
	}
}

func (w *World) Camera() *Camera {
	return w.camera
}

func (w *World) AddNodeToLayer(node SceneNode, layerIndex int) {
	if layerIndex < 0 {
		return
	}
	w.ensureLayerRoots(layerIndex)
	w.layerRoots[layerIndex].AddChildren(node)
}

// ensureLayerRoots grows layerRoots until layerIndex is valid, appending one scene root per new layer.
func (w *World) ensureLayerRoots(layerIndex int) {
	for layerIndex >= len(w.layerRoots) {
		root := NewNode(NameLayerRoot)
		w.layerRoots = append(w.layerRoots, root)
		w.rootScene.AddChildren(root)
	}
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
	sortSceneChildrenByDrawableLayerDesc(children)

	for _, child := range children {
		w.queueNode(child, layerIndex)
	}

	if drawable, ok := node.(Drawable); ok {
		op := drawImageOptionsForSceneNode(node)
		w.camera.ApplyOffset(&op)
		w.drawLayers.AddNode(layerIndex, drawable, w.drawTarget, op)
	}
}

// sortSceneChildrenByDrawableLayerDesc sorts the node's child slice in place (higher Drawable layer first).
func sortSceneChildrenByDrawableLayerDesc(children []SceneNode) {
	sort.Slice(children, func(i, j int) bool {
		return drawableSortLayer(children[i]) > drawableSortLayer(children[j])
	})
}

func drawableSortLayer(n SceneNode) int {
	if d, ok := n.(Drawable); ok {
		return d.GetLayer()
	}
	return 0
}

func drawImageOptionsForSceneNode(node SceneNode) ebiten.DrawImageOptions {
	op := ebiten.DrawImageOptions{}
	if tr, ok := node.(Transformable); ok {
		op.GeoM = buildGeoMFromTransform(tr.GetWorldTransform())
	}
	return op
}

func (w *World) Draw(target *ebiten.Image) {
	w.camera.Update()
	w.prepareDrawSurface()

	for i := range w.layerRoots {
		w.queueNode(w.layerRoots[i], i)
	}
	w.drawLayers.DrawAll()
	w.camera.DrawToScreen(target)
}

func (w *World) prepareDrawSurface() {
	w.drawTarget = w.camera.GetSurface()
	w.drawTarget.Fill(color.RGBA{40, 40, 50, 255})
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
