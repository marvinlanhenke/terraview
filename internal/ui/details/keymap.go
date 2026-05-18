package details

import (
	"charm.land/bubbles/v2/key"
)

type keymap struct {
	toggle key.Binding
}

var keys = keymap{
	toggle: key.NewBinding(
		key.WithKeys("p", "t"),
		key.WithHelp("p/t", "toggle diff/plan"),
	),
}
