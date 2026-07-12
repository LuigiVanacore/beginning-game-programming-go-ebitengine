package pickups

import (
	. "book/code/ch08/internal/core"
)

// PickupManager owns the active orb list and the per-frame collection queue.
type PickupManager struct {
	orbs          []*Orb
	toCollectOrbs []*Orb
}

// NewPickupManager creates an empty PickupManager.
func NewPickupManager() *PickupManager {
	return &PickupManager{
		orbs:          make([]*Orb, 0),
		toCollectOrbs: make([]*Orb, 0),
	}
}

// AddOrb appends a pre-created orb to the active list.
func (m *PickupManager) AddOrb(o *Orb) { m.orbs = append(m.orbs, o) }

// Orbs returns the active orb slice (used by collision callbacks).
func (m *PickupManager) Orbs() []*Orb { return m.orbs }

// QueueOrbCollect records an orb touched by the player this frame.
func (m *PickupManager) QueueOrbCollect(o *Orb) {
	m.toCollectOrbs = append(m.toCollectOrbs, o)
}

// CollectOrbs applies XP for every queued orb, removes it from the world, and clears the queue.
func (m *PickupManager) CollectOrbs(world *World, cm *CollisionManager, onCollect func(int)) {
	for _, picked := range m.toCollectOrbs {
		onCollect(picked.Value)
		world.RemoveNode(picked.Col)
		cm.RemoveCollider(picked.Col)
		for i, active := range m.orbs {
			if active == picked {
				m.orbs = append(m.orbs[:i], m.orbs[i+1:]...)
				break
			}
		}
	}
	m.toCollectOrbs = m.toCollectOrbs[:0]
}
