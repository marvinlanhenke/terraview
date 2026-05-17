package app

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/marvinlanhenke/terraview/internal/planview"
)

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
		m.components.status.SetWidth(msg.Width - defaultMargin)

		treeWidth, treeHeight := treePaneSize(msg.Width, msg.Height)
		m.components.tree.SetSize(treeWidth, treeHeight)

		detailsWidth := detailsWidth(m.size.width, treeWidth)
		detailsHeight := treeHeight
		m.components.details.SetSize(detailsWidth, detailsHeight)

		return m, nil

	case tea.KeyPressMsg:
		switch {
		// Quit
		case key.Matches(msg, keys.Quit) && m.focus != FocusSearch:
			return m, tea.Quit

		// Search Focus
		case key.Matches(msg, keys.Search) && m.focus != FocusSearch:
			m.focus = FocusSearch
			m.components.details.Blur()
			cmds = append(cmds, m.components.search.Focus())
			// We return here, such that the `/` is not used as query input
			return m, tea.Batch(cmds...)

		// Search Enter
		case key.Matches(msg, keys.Enter) && m.focus == FocusSearch:
			m.focus = FocusTree
			m.components.search.Blur()

		// Search Escape
		case key.Matches(msg, keys.Escape) && m.focus == FocusSearch:
			m.focus = FocusTree
			m.components.search.Clear()
			m.components.search.Blur()
			m.components.tree.ApplyQuery("")
			return m, nil

		// Tree -> Details
		case key.Matches(msg, keys.RightPane) && m.focus == FocusTree:
			m.focus = FocusDetails
			m.components.details.Focus()

		case key.Matches(msg, keys.Enter) && m.focus == FocusTree:
			selected := m.components.tree.Selected()
			if selected != nil && selected.Kind == planview.NodeResource {
				m.focus = FocusDetails
				m.components.details.Focus()
			}

		// Details -> Tree
		case (key.Matches(msg, keys.Escape) || key.Matches(msg, keys.Enter) || key.Matches(msg, keys.LeftPane)) && m.focus == FocusDetails:
			m.focus = FocusTree
			m.components.details.Blur()

		// Filter Focus
		case (key.Matches(msg, keys.Filter)) && m.focus != FocusSearch:
			if m.focus == FocusFilter {
				m.focus = FocusTree
			} else {
				m.focus = FocusFilter
			}

		// Filter Exit
		case (key.Matches(msg, keys.Escape)) && m.focus == FocusFilter:
			m.focus = FocusTree
		}

	}

	switch m.focus {
	case FocusSearch:
		cmds = append(cmds, m.components.search.Update(msg))
		m.components.tree.ApplyQuery(m.components.search.Value())
		m.components.search.SetMatches(m.components.tree.GetVisible())

	case FocusTree:
		cmds = append(cmds, m.components.tree.Update(msg))
		m.components.details.SetNode(m.components.tree.Selected())
		m.components.search.SetMatches(m.components.tree.GetVisible())

	case FocusDetails:
		cmds = append(cmds, m.components.details.Update(msg))

	case FocusFilter:
		cmds = append(cmds, m.components.filter.Update(msg))
		m.components.tree.ApplyFilters(m.components.filter.GetFilters())
		m.components.details.SetNode(m.components.tree.Selected())
		m.components.search.SetMatches(m.components.tree.GetVisible())
	}

	return m, tea.Batch(cmds...)
}
