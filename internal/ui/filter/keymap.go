package filter

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
)

type keymap struct {
	up     key.Binding
	down   key.Binding
	reset  key.Binding
	toggle key.Binding
}

var keys = keymap{
	up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	reset: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset all filters"),
	),
	toggle: key.NewBinding(
		key.WithKeys("enter", "space"),
		key.WithHelp("enter/space", "toggle filter"),
	),
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.up,
		k.down,
		k.reset,
		k.toggle,
	}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.up, k.down},
		{k.reset, k.toggle},
	}
}

// KeyMap returns the key bindings used by the filter modal.
func KeyMap() help.KeyMap {
	return keys
}
