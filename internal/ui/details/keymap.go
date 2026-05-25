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

// ShortHelp returns the compact help bindings for the details keymap.
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.toggle}
}

// FullHelp returns the expanded help bindings for the details keymap.
func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.toggle},
	}
}

// KeyMap returns the key bindings used by the details pane.
func KeyMap() help.KeyMap {
	return keys
}
