package filter

import "github.com/marvinlanhenke/terraview/internal/ui"

type intentKind uint8

const (
	intentNone intentKind = iota
	intentToggle
	intentReset
)

type Intent struct {
	kind   intentKind
	action ui.Action
}

func (i Intent) Toggle() (ui.Action, bool) {
	if i.kind != intentToggle {
		return "", false
	}

	return i.action, true
}

func (i Intent) Reset() bool {
	return i.kind == intentReset
}

func ToggleIntent(action ui.Action) Intent {
	return Intent{kind: intentToggle, action: action}
}

func ResetIntent() Intent {
	return Intent{kind: intentReset}
}
