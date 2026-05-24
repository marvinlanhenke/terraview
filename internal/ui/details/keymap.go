package details

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
)

type keymap struct {
	toggle key.Binding
}

var keys = keymap{
	toggle: key.NewBinding(
		key.WithKeys("p", "t"),
		key.WithHelp("p/t", "switch diff/plan"),
	),
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.toggle}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.toggle},
	}
}

func KeyMap() help.KeyMap {
	return keys
}
