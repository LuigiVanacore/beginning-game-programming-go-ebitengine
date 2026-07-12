package pickups

import (
	. "book/code/ch13/internal/core"
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

// AddPotion appends a pre-created potion to the active list.
func (m *PickupManager) AddPotion(p *PotionEnt) { m.potions = append(m.potions, p) }

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
	for _, orb := range m.toCollectOrbs {
		onCollect(orb.Value)
		world.RemoveNode(orb)
		cm.RemoveCollider(orb.Col)
		for i, o := range m.orbs {
			if o == orb {
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
		for i, existing := range m.potions {
			if existing == p {
				m.potions = append(m.potions[:i], m.potions[i+1:]...)
				break
			}
		}
	}
	m.toCollectPotions = m.toCollectPotions[:0]
}

// MaybeDropPotion randomly spawns a potion at (x, y).
func (m *PickupManager) MaybeDropPotion(engine *Engine, rm *ResourceManager, x, y float64) {
	if rand.Float64() >= potionDropChance {
		return
	}
	potionTex, _ := rm.GetTexture("potion")
	if potionTex == nil {
		return
	}
	m.potions = append(m.potions, CreatePotion(engine, x, y, potionTex))
}
