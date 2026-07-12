package core

import "github.com/hajimehoshi/ebiten/v2"

// ActionID identifies an action in the ActionMap (e.g. MoveUp, MoveDown).
type ActionID int

// ActionMap maps ActionID to Action for the state-based input API.
// Used with RegisterAction and IsActionActive for an ActionID-based workflow.
type ActionMap struct {
	actions map[ActionID]Action
}

// NewActionMap creates an empty ActionMap.
func NewActionMap() *ActionMap {
	return &ActionMap{
		actions: make(map[ActionID]Action),
	}
}

// AddAction binds an Action to the given ActionID.
func (am *ActionMap) AddAction(id ActionID, action Action) {
	am.actions[id] = action
}

// RemoveAction removes the action for the given ID.
func (am *ActionMap) RemoveAction(id ActionID) {
	delete(am.actions, id)
}

// GetAction returns the action for the given ID, and whether it exists.
func (am *ActionMap) GetAction(id ActionID) (Action, bool) {
	action, exists := am.actions[id]
	return action, exists
}

// ClearActions removes all actions.
func (am *ActionMap) ClearActions() {
	am.actions = make(map[ActionID]Action)
}

// SetKeyBinding rebinds the action for the given ID to a new key.
// Only KeyAction types are modified; others are replaced with a new KeyAction.
func (am *ActionMap) SetKeyBinding(id ActionID, key ebiten.Key) error {
	action, exists := am.actions[id]
	if !exists {
		return ErrActionNotFound
	}
	am.actions[id] = NewKeyAction(key, action.GetMode())
	return nil
}

// ErrActionNotFound is returned when an ActionID is not found.
var ErrActionNotFound = &ActionNotFoundError{}

// ActionNotFoundError indicates the action ID was not found.
type ActionNotFoundError struct{}

func (e *ActionNotFoundError) Error() string {
	return "action not found"
}
