package app

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/marvinlanhenke/terraview/internal/ui/filter"
)

// Init satisfies tea.Model and returns no initial command.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update satisfies tea.Model and routes messages to the focused component.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.applyLayout(msg.Width, msg.Height)
		return m, nil

	case tea.KeyPressMsg:
		switch {
		// Quit
		case key.Matches(msg, keys.Quit) && m.focus != focusSearch:
			return m, tea.Quit

		// Search Focus
		case key.Matches(msg, keys.Search) && m.focus != focusSearch:
			m.focus = focusSearch
			m.components.details.Blur()
			cmds = append(cmds, m.components.search.Focus())
			// We return here, such that the `/` is not used as query input
			return m, tea.Batch(cmds...)

		// Search Enter
		case key.Matches(msg, keys.Enter) && m.focus == focusSearch:
			m.focus = focusTree
			m.components.search.Blur()
			return m, nil

		// Search Escape
		case key.Matches(msg, keys.Escape) && m.focus == focusSearch:
			m.focus = focusTree
			m.components.search.Clear()
			m.components.search.Blur()
			m.controls.query = ""
			m.refreshTreeFromControls()
			return m, nil

		// Tree -> Details
		case (key.Matches(msg, keys.RightPane) || key.Matches(msg, keys.Enter)) && m.focus == focusTree:
			selected := m.components.tree.Selected()
			if selected != nil && selected.IsResource() {
				m.focus = focusDetails
				m.components.details.Focus()
			}

		// Details -> Tree
		case (key.Matches(msg, keys.Escape) || key.Matches(msg, keys.Enter) || key.Matches(msg, keys.LeftPane)) && m.focus == focusDetails:
			m.focus = focusTree
			m.components.details.Blur()

		// Filter Focus
		case key.Matches(msg, keys.Filter) && m.focus != focusSearch:
			if m.focus == focusFilter {
				m.focus = focusTree
				m.components.details.Blur()
			} else {
				m.focus = focusFilter
			}

		// Filter Exit
		case key.Matches(msg, keys.Escape) && m.focus == focusFilter:
			m.focus = focusTree
		}

	}

	switch m.focus {
	case focusSearch:
		cmds = append(cmds, m.components.search.Update(msg))
		if m.applySearchQuery() {
			m.refreshTreeFromControls()
		}

	case focusTree:
		cmds = append(cmds, m.components.tree.Update(msg))
		m.syncTreeOutputs()

	case focusDetails:
		cmds = append(cmds, m.components.details.Update(msg))

	case focusFilter:
		intent, cmd := m.components.filter.Update(msg)

		cmds = append(cmds, cmd)

		if m.applyFilterIntent(intent) {
			m.refreshTreeFromControls()
		}

		m.components.status.SetActiveFilterCount(len(m.controls.filters))
	}

	return m, tea.Batch(cmds...)
}

// applySearchQuery stores the current search query and reports whether it changed.
func (m *Model) applySearchQuery() bool {
	nextQuery := m.components.search.Value()

	if nextQuery != m.controls.query {
		m.controls.query = nextQuery
		return true
	}

	return false
}

// applyFilterIntent applies filter modal state changes and reports whether filters changed.
func (m *Model) applyFilterIntent(intent filter.Intent) bool {
	if intent.Reset {
		if len(m.controls.filters) == 0 {
			return false
		}

		clear(m.controls.filters)

		return true
	}

	if intent.HasToggle {
		if m.controls.filters[intent.Action] {
			delete(m.controls.filters, intent.Action)
		} else {
			m.controls.filters[intent.Action] = true
		}

		return true
	}

	return false
}
