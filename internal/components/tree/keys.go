package tree

import (
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
		key.WithKeys("enter", "right", "l"),
		key.WithHelp("enter", "expand"),
	),
	collapse: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "collapse"),
	),
}
