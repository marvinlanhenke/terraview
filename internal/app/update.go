package app

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/marvinlanhenke/terraview/internal/ui/filter"
	"github.com/marvinlanhenke/terraview/internal/ui/tree"
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
			m.controls.query = ""
			m.refreshTreeFromControls()
			return m, nil

		// Tree -> Details
		case key.Matches(msg, keys.RightPane) && m.focus == FocusTree:
			m.focus = FocusDetails
			m.components.details.Focus()

		case key.Matches(msg, keys.Enter) && m.focus == FocusTree:
			selected := m.components.tree.Selected()
			if selected != nil && selected.IsResource() {
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
		if m.applySearchQuery() {
			m.refreshTreeFromControls()
		}

	case FocusTree:
		cmds = append(cmds, m.components.tree.Update(msg))
		m.syncTreeOutputs()

	case FocusDetails:
		cmds = append(cmds, m.components.details.Update(msg))

	case FocusFilter:
		intent, cmd := m.components.filter.Update(msg)

		cmds = append(cmds, cmd)

		if m.applyFilterIntent(intent) {
			m.refreshTreeFromControls()
		}

		m.components.status.SetActiveFilterCount(len(m.controls.filters))
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) applySearchQuery() bool {
	nextQuery := m.components.search.Value()

	if nextQuery != m.controls.query {
		m.controls.query = nextQuery
		return true
	}

	return false
}

func (m *Model) applyFilterIntent(intent filter.Intent) bool {
	if intent.Reset {
		if len(m.controls.filters) == 0 {
			return false
		}

		clear(m.controls.filters)

		return true
	}

	if intent.HasToggle {
		if m.controls.filters[tree.Action(intent.Action)] {
			delete(m.controls.filters, tree.Action(intent.Action))
		} else {
			m.controls.filters[tree.Action(intent.Action)] = true
		}

		return true
	}

	return false
}
