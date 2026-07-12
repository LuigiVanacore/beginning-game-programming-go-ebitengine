package core

import "github.com/hajimehoshi/ebiten/v2"

// Layers manages draw order: lower layer index = drawn first (background).
// Within a layer, nodes are drawn in LIFO order (last pushed = drawn first).
type Layers struct {
	layers []*Stack[func()]
}

func NewLayers() *Layers {
	l := &Layers{layers: make([]*Stack[func()], 8)}
	for i := range l.layers {
		l.layers[i] = NewStack[func()]()
	}
	return l
}

func (l *Layers) ensureLayer(idx int) {
	for idx >= len(l.layers) {
		l.layers = append(l.layers, NewStack[func()]())
	}
}

func (l *Layers) AddDraw(layerIndex int, f func()) {
	if layerIndex < 0 {
		return
	}
	l.ensureLayer(layerIndex)
	l.layers[layerIndex].Push(f)
}

func (l *Layers) AddNode(layerIndex int, node Drawable, target *ebiten.Image, op ebiten.DrawImageOptions) {
	if layerIndex < 0 {
		return
	}
	opCopy := op
	l.AddDraw(layerIndex, func() {
		node.Draw(target, &opCopy)
	})
}

func (l *Layers) DrawAll() {
	for _, stk := range l.layers {
		for !stk.IsEmpty() {
			if f, ok := stk.Pop(); ok {
				f()
			}
		}
	}
}
