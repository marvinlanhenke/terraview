package app

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type keymap struct {
	Search key.Binding
	Enter  key.Binding
	Escape key.Binding
	Quit   key.Binding
}

var keys = keymap{
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "apply search"),
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
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size.width = msg.Width
		m.size.height = msg.Height

		m.components.search.SetWidth(msg.Width - defaultMargin)
		m.components.summary.SetWidth(msg.Width - defaultMargin)

		treeWidth, treeHeight := treePaneSize(msg.Width, msg.Height)
		m.components.tree.SetSize(treeWidth, treeHeight)

		return m, nil

	case tea.KeyPressMsg:
		switch {
		// Quit
		case key.Matches(msg, keys.Quit) && m.focus != FocusSearch:
			return m, tea.Quit

		// Search Focus
		case key.Matches(msg, keys.Search) && m.focus != FocusSearch:
			m.focus = FocusSearch
			cmds = append(cmds, m.components.search.Focus())
			// We return here, such that the `/` is not used as query input
			return m, tea.Batch(cmds...)

		// Search Enter
		case key.Matches(msg, keys.Enter, keys.Enter) && m.focus == FocusSearch:
			m.focus = FocusTree
			m.components.search.Blur()

		case key.Matches(msg, keys.Escape) && m.focus == FocusSearch:
			m.focus = FocusTree
			m.components.search.Clear()
			m.components.search.Blur()
			m.components.tree.ApplyFilter("")
			return m, nil
		}
	}

	switch m.focus {
	case FocusSearch:
		cmds = append(cmds, m.components.search.Update(msg))
		m.components.tree.ApplyFilter(m.components.search.Value())
		m.components.search.SetMatches(m.components.tree.GetVisible())

	case FocusTree:
		cmds = append(cmds, m.components.tree.Update(msg))

	case FocusDetails:
		// TODO
	}

	return m, tea.Batch(cmds...)
}
