package filter

import "github.com/marvinlanhenke/terraview/internal/ui"

type intentKind uint8

const (
	intentNone intentKind = iota
	intentToggle
	intentReset
)

// Intent describes a filter action requested by the modal.
type Intent struct {
	kind   intentKind
	action ui.Action
}

// Toggle reports the action to toggle when the intent is a toggle request.
func (i Intent) Toggle() (ui.Action, bool) {
	if i.kind != intentToggle {
		return "", false
	}

	return i.action, true
}

// Reset reports whether the intent is a reset request.
func (i Intent) Reset() bool {
	return i.kind == intentReset
}

// ToggleIntent returns an intent that toggles action.
func ToggleIntent(action ui.Action) Intent {
	return Intent{kind: intentToggle, action: action}
}

// ResetIntent returns an intent that resets all filters.
func ResetIntent() Intent {
	return Intent{kind: intentReset}
}
