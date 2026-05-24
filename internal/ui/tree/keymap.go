package tree

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
)

type keymap struct {
	up          key.Binding
	down        key.Binding
	expand      key.Binding
	expandAll   key.Binding
	collapse    key.Binding
	collapseAll key.Binding
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
	expandAll: key.NewBinding(
		key.WithKeys("ctrl+e"),
		key.WithHelp("ctrl+e", "expand all"),
	),
	collapse: key.NewBinding(
		key.WithKeys("enter", "c", "left", "h"),
		key.WithHelp("←/h", "collapse"),
	),
	collapseAll: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "collapse all"),
	),
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.up,
		k.down,
		k.expand,
		k.expandAll,
		k.collapse,
		k.collapseAll,
	}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.up, k.down},
		{k.expand, k.expandAll},
		{k.collapse, k.collapseAll},
	}
}

func KeyMap() help.KeyMap {
	return keys
}
