package tree

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
)

type keymap struct {
	up       key.Binding
	down     key.Binding
	expand   key.Binding
	collapse key.Binding
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
	expand: key.NewBinding(
		key.WithKeys("enter", "e", "right", "l"),
		key.WithHelp("→/l", "expand"),
	),
	collapse: key.NewBinding(
		key.WithKeys("enter", "c", "left", "h"),
		key.WithHelp("←/h", "collapse"),
	),
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.up,
		k.down,
		k.expand,
		k.collapse,
	}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.up, k.down},
		{k.expand, k.collapse},
	}
}

func KeyMap() help.KeyMap {
	return keys
}
