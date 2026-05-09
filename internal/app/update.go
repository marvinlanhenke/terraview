package app

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type keymap struct {
	Search key.Binding
	Escape key.Binding
	Quit   key.Binding
}

var keys = keymap{
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (m Model) Init() tea.Cmd {
	return m.search.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.search.SetWidth(max(0, msg.Width-8))
		return m, nil

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.Quit) && m.focus != FocusSearch:
			return m, tea.Quit

		case key.Matches(msg, keys.Search) && m.focus != FocusSearch:
			m.focus = FocusSearch
			cmds = append(cmds, m.search.Focus())
			return m, tea.Batch(cmds...)

		case key.Matches(msg, keys.Escape) && m.focus == FocusSearch:
			m.focus = FocusTree
			m.search.Clear()
			m.search.Blur()
			return m, nil
		}
	}

	switch m.focus {
	case FocusSearch:
		cmds = append(cmds, m.search.Update(msg))

	case FocusTree:
		// TODO

	case FocusDetails:
		// TODO
	}

	return m, tea.Batch(cmds...)
}
