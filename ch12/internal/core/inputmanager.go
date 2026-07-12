package core

// ActionBinding maps action names to RawInputButtons.
// Multiple buttons can trigger the same action (e.g. arrows and WASD for movement).
type ActionBinding []RawInputButton

// InputManager handles input with two APIs:
// 1) String-based: AddAction(name, RawInputButton), IsActionPressed(name)
// 2) ActionID-based: RegisterAction(id, Action), IsActionActive(id) via ActionMap + StateBuffer
type InputManager struct {
	actions   map[string]ActionBinding
	actionMap *ActionMap
	stateBuf  *StateBuffer
}

// NewInputManager creates an InputManager with empty bindings.
func NewInputManager() *InputManager {
	return &InputManager{
		actions:   make(map[string]ActionBinding),
		actionMap: NewActionMap(),
		stateBuf:  NewStateBuffer(),
	}
}

// Update polls hardware and updates the StateBuffer.
// Call this once per frame (typically from Engine.Update).
func (i *InputManager) Update() {
	i.stateBuf.Update()
}

// AddAction binds a RawInputButton to an action name.
// Multiple buttons can be bound to the same action (e.g. W and ArrowUp for move_up).
func (i *InputManager) AddAction(actionName string, button RawInputButton) {
	if _, exists := i.actions[actionName]; !exists {
		i.actions[actionName] = ActionBinding{}
	}
	i.actions[actionName] = append(i.actions[actionName], button)
}

// RemoveAction removes the action mapping entirely.
func (i *InputManager) RemoveAction(actionName string) {
	delete(i.actions, actionName)
}

// IsActionPressed returns true if any button bound to the action is currently held.
func (i *InputManager) IsActionPressed(actionName string) bool {
	bindings, exists := i.actions[actionName]
	if !exists {
		return false
	}
	for _, btn := range bindings {
		if btn.IsPressed() {
			return true
		}
	}
	return false
}

// IsActionJustPressed returns true only on the first frame any bound button is pressed.
func (i *InputManager) IsActionJustPressed(actionName string) bool {
	bindings, exists := i.actions[actionName]
	if !exists {
		return false
	}
	for _, btn := range bindings {
		if btn.IsJustPressed() {
			return true
		}
	}
	return false
}

// IsActionJustReleased returns true only on the frame any bound button is released.
func (i *InputManager) IsActionJustReleased(actionName string) bool {
	bindings, exists := i.actions[actionName]
	if !exists {
		return false
	}
	for _, btn := range bindings {
		if btn.IsJustReleased() {
			return true
		}
	}
	return false
}

// IsActionReleased returns true only if all buttons bound to the action are unpressed.
func (i *InputManager) IsActionReleased(actionName string) bool {
	bindings, exists := i.actions[actionName]
	if !exists {
		return true
	}
	for _, btn := range bindings {
		if btn.IsPressed() {
			return false
		}
	}
	return true
}

// RegisterAction binds an Action to an ActionID for the state-based API.
func (i *InputManager) RegisterAction(id ActionID, action Action) {
	i.actionMap.AddAction(id, action)
}

// IsActionActive returns true if the action for the given ID is active (per its mode).
func (i *InputManager) IsActionActive(id ActionID) bool {
	action, ok := i.actionMap.GetAction(id)
	if !ok {
		return false
	}
	return i.stateBuf.IsActionActive(action)
}

// ActionMap returns the ActionMap for direct access (e.g. SetKeyBinding).
func (i *InputManager) ActionMap() *ActionMap {
	return i.actionMap
}
