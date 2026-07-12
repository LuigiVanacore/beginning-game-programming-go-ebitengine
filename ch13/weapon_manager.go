package game

import (
	. "book/code/ch13/internal/core"
)

// weaponColliders is implemented by weapons that own a persistent collider (the
// contact-damage weapons: Holy Shield and Sacred Book). Their colliders must be
// registered with the CollisionManager only while the weapon is mounted, otherwise
// an unmounted weapon would leave a live damaging collider at the world origin.
type weaponColliders interface {
	weaponColliders() []*Collider
}

// WeaponManager owns weapon instances, upgrade multipliers, unlock flags, and the active playerWeapon list.
// It replaces the per-player weapon loop from ch09 so Game can control weapon ticks centrally.
type WeaponManager struct {
	weapons []playerWeapon
	cm      *CollisionManager // registers/unregisters persistent weapon colliders on mount/unmount

	// Shared multipliers (1.0 = base). Per-weapon multipliers now live on the weapon structs.
	HolyShieldRadiusMult float64
	WeaponDamageMult     float64
	GlobalCooldownMult   float64

	KnifeUnlocked, FlyingAxeUnlocked       bool
	SacredBookUnlocked, HolyShieldUnlocked bool

	KnifeSpeedLvl, KnifeCooldownLvl                               int
	SacredBookSpeedLvl                                            int
	HolyShieldRadiusLvl                                           int
	FlyingAxeSpeedLvl, FlyingAxeCooldownLvl, FlyingAxeRotationLvl int

	Knife      *BloodyKnifeWeapon
	Axe        *FlyingAxeWeapon
	SacredBook *SacredBook
	HolyShield *HolyShieldWeapon
}

// NewWeaponLoadout builds all weapon instances and default multipliers for a new game session.
func NewWeaponLoadout(engine *Engine, player *Player) *WeaponManager {
	rm := engine.ResourceManager()
	wm := &WeaponManager{
		weapons:              make([]playerWeapon, 0),
		cm:                   engine.CollisionManager(),
		HolyShieldRadiusMult: 1.0,
		WeaponDamageMult:     1.0,
		GlobalCooldownMult:   1.0,
	}

	knifeTex, _ := rm.GetTexture(BloodyKnifeTexture)
	knifePool := NewProjectilePool(engine, knifeTex, knifeRadius, knifeProjectileSpeedPxPerFrame, DrawLayerPlayer, NameBloodyKnife, projectilePoolSize)
	wm.Knife = NewBloodyKnifeWeapon(knifePool, &wm.GlobalCooldownMult, &wm.WeaponDamageMult)

	axeTex, _ := rm.GetTexture(FlyingAxeTexture)
	axePool := NewProjectilePool(engine, axeTex, flyingAxeRadius, flyingAxeProjectileSpeedPxPerFrame, DrawLayerPlayer, NameFlyingAxe, projectilePoolSize)
	wm.Axe = NewFlyingAxeWeapon(axePool, &wm.GlobalCooldownMult, &wm.WeaponDamageMult)

	if sacredBookTex, ok := rm.GetTexture(SacredBookTexture); ok {
		wm.SacredBook = NewSacredBook(engine, player, OrbitWeaponDistance, sacredBookTex)
	}
	wm.HolyShield = NewHolyShieldWeapon(engine, player, HolyShieldRadius)

	return wm
}

// Mount adds a weapon to the scene graph under parent, registers it for updates,
// and activates its persistent collider (if any) in the CollisionManager. A weapon
// only deals contact damage while mounted.
func (m *WeaponManager) Mount(w playerWeapon, parent SceneNode) {
	if m == nil || w == nil || parent == nil {
		return
	}
	parent.AddChildren(w)
	m.weapons = append(m.weapons, w)
	if wc, ok := w.(weaponColliders); ok && m.cm != nil {
		for _, c := range wc.weaponColliders() {
			m.cm.AddCollider(c)
		}
	}
}

// Unregister removes a weapon from the update list and deactivates its persistent
// collider (e.g. before replacing Holy Shield on a radius upgrade).
func (m *WeaponManager) Unregister(w playerWeapon) {
	if m == nil || w == nil {
		return
	}
	if wc, ok := w.(weaponColliders); ok && m.cm != nil {
		for _, c := range wc.weaponColliders() {
			m.cm.RemoveCollider(c)
		}
	}
	for i, x := range m.weapons {
		if x == w {
			m.weapons = append(m.weapons[:i], m.weapons[i+1:]...)
			return
		}
	}
}

// Update calls updateWeapon on every mounted weapon.
func (m *WeaponManager) Update(engine *Engine, player *Player, cursor *Cursor) {
	if m == nil {
		return
	}
	for _, w := range m.weapons {
		w.updateWeapon(engine, player, cursor)
	}
}

// TryReleaseProjectileByCollider returns a pooled projectile to its weapon when hit,
// replacing the ch09 Player.TryReleaseProjectileByCollider call in processRemovals.
func (m *WeaponManager) TryReleaseProjectileByCollider(c *Collider) bool {
	if m == nil {
		return false
	}
	for _, w := range m.weapons {
		if pc, ok := w.(projectileCarrier); ok && pc.tryReleaseProjectile(c) {
			return true
		}
	}
	return false
}
