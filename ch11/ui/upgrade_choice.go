package ui

// UpgradeChoice holds display fields for an upgrade panel (no game callback; use parallel []func() from game).
type UpgradeChoice struct {
	WeaponName  string
	UpgradeDesc string
	IconKey     string
}
