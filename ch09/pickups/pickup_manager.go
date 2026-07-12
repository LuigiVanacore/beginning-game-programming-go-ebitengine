package pickups

import (
	. "book/code/ch09/internal/core"
	"math/rand"
)

// PickupManager owns active orb and potion lists and their per-frame collection queues.
type PickupManager struct {
	orbs             []*OrbEnt
	potions          []*PotionEnt
	toCollectOrbs    []*OrbEnt
	toCollectPotions []*PotionEnt
}

// NewPickupManager creates an empty PickupManager.
func NewPickupManager() *PickupManager {
	return &PickupManager{
		orbs:             make([]*OrbEnt, 0),
		potions:          make([]*PotionEnt, 0),
		toCollectOrbs:    make([]*OrbEnt, 0),
		toCollectPotions: make([]*PotionEnt, 0),
	}
}

// AddOrb appends a pre-created orb to the active list.
func (m *PickupManager) AddOrb(o *OrbEnt) { m.orbs = append(m.orbs, o) }

// MaybeDropPotion randomly spawns a potion at (x, y).
func (m *PickupManager) MaybeDropPotion(engine *Engine, x, y float64) {
	tex, _ := engine.ResourceManager().GetTexture(PotionTexture)
	if tex == nil {
		return
	}
	if rand.Float64() >= potionDropChance {
		return
	}
	m.potions = append(m.potions, CreatePotion(engine, x, y, tex))
}

// Orbs returns the active orb slice (used by collision callbacks).
func (m *PickupManager) Orbs() []*OrbEnt { return m.orbs }

// Potions returns the active potion slice (used by collision callbacks).
func (m *PickupManager) Potions() []*PotionEnt { return m.potions }

// QueueOrbCollect records an orb touched by the player this frame.
func (m *PickupManager) QueueOrbCollect(o *OrbEnt) {
	m.toCollectOrbs = append(m.toCollectOrbs, o)
}

// QueuePotionCollect records a potion touched by the player this frame.
func (m *PickupManager) QueuePotionCollect(p *PotionEnt) {
	m.toCollectPotions = append(m.toCollectPotions, p)
}

// CollectOrbs applies XP for every queued orb, removes it from the world, and clears the queue.
func (m *PickupManager) CollectOrbs(world *World, cm *CollisionManager, onCollect func(int)) {
	for _, picked := range m.toCollectOrbs {
		onCollect(picked.Value)
		world.RemoveNode(picked)
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

// CollectPotions applies healing for every queued potion, removes it from the world, and clears the queue.
func (m *PickupManager) CollectPotions(world *World, cm *CollisionManager, onCollect func()) {
	for _, p := range m.toCollectPotions {
		onCollect()
		world.RemoveNode(p)
		cm.RemoveCollider(p.Col)
		for i, active := range m.potions {
			if active == p {
				m.potions = append(m.potions[:i], m.potions[i+1:]...)
				break
			}
		}
	}
	m.toCollectPotions = m.toCollectPotions[:0]
}
