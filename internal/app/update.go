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

// Update satisfies tea.Model by routing app-level messages before focused input.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.applyLayout(msg.Width, msg.Height)
		m.logger.Debug("window resized", "width", msg.Width, "height", msg.Height)
		return m, nil

	case tea.KeyPressMsg:
		cmd, consumed := m.routeKeyPress(msg)
		cmds = append(cmds, cmd)

		if consumed {
			m.logger.Debug("tea cmd consumed", "key", msg.String(), "consumed", consumed)
			return m, tea.Batch(cmds...)
		}
	}

	cmds = append(cmds, m.updateFocused(msg))

	return m, tea.Batch(cmds...)
}

// routeKeyPress applies app-level key bindings before component-specific updates.
// A false consumed value lets the focused component also process the keypress.
func (m *Model) routeKeyPress(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	switch {
	case key.Matches(msg, keys.Quit) && m.focus != focusSearch:
		m.logger.Debug("key match", "key", keys.Quit)
		return tea.Quit, true

	case key.Matches(msg, keys.Search) && m.focus != focusSearch:
		m.logger.Debug("key match", "key", keys.Search)
		return m.focusSearch(), true

	case key.Matches(msg, keys.Enter) && m.focus == focusSearch:
		m.logger.Debug("key match", "key", keys.Enter, "focus", "search", "action", "apply search")
		m.focusTree()
		return nil, true

	case key.Matches(msg, keys.Escape) && m.focus == focusSearch:
		m.logger.Debug("key match", "key", keys.Escape, "focus", "search", "action", "clear and exit search")
		m.components.search.Clear()
		m.controls.query = ""
		m.refreshTreeFromControls()
		m.focusTree()
		return nil, true

	case m.focus == focusTree && (key.Matches(msg, keys.RightPane) || key.Matches(msg, keys.Enter)):
		m.logger.Debug("key match", "key", msg.String(), "focus", "tree", "action", "inspect selected")
		m.focusDetailsIfInspectable()
		return nil, false

	case m.focus == focusDetails && (key.Matches(msg, keys.LeftPane) || key.Matches(msg, keys.Enter) || key.Matches(msg, keys.Escape)):
		m.logger.Debug("key match", "key", msg.String(), "focus", "details", "action", "return to tree")
		m.focusTree()
		return nil, false

	case key.Matches(msg, keys.Filter) && m.focus != focusSearch:
		m.logger.Debug("key match", "key", keys.Filter, "action", "toggle filter")
		m.toggleFilter()
		return nil, false

	case key.Matches(msg, keys.Escape) && m.focus == focusFilter:
		m.logger.Debug("key match", "key", keys.Escape, "focus", "filter", "action", "close filter")
		m.focusTree()
		return nil, false
	}

	return nil, false
}

// focusSearch moves focus to the search field and starts cursor blinking.
func (m *Model) focusSearch() tea.Cmd {
	m.logger.Debug("focus changed", "from", m.focus, "to", focusSearch)
	m.focus = focusSearch
	m.components.details.Blur()
	return m.components.search.Focus()
}

// focusTree moves focus back to the tree and clears child component focus styles.
func (m *Model) focusTree() {
	m.logger.Debug("focus changed", "from", m.focus, "to", focusTree)
	m.focus = focusTree
	m.components.search.Blur()
	m.components.details.Blur()
}

// focusDetailsIfInspectable moves focus to details when a resource row is selected.
func (m *Model) focusDetailsIfInspectable() {
	selected := m.components.tree.Selected()
	if selected == nil || !selected.IsInspectable() {
		m.logger.Debug("focus details skipped", "reason", "selected node not inspectable")
		return
	}

	m.logger.Debug("focus changed", "from", m.focus, "to", focusDetails, "node", selected.Id)
	m.focus = focusDetails
	m.components.search.Blur()
	m.components.details.Blur()
	m.components.details.Focus()
}

// toggleFilter opens the filter modal, or closes it when already active.
func (m *Model) toggleFilter() {
	if m.focus == focusFilter {
		m.logger.Debug("filter modal closed")
		m.focusTree()
		return
	}

	m.logger.Debug("filter modal opened", "from", m.focus)
	m.focus = focusFilter
	m.components.search.Blur()
	m.components.details.Blur()
}

// updateFocused forwards a message to the currently focused child component.
func (m *Model) updateFocused(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch m.focus {
	case focusSearch:
		cmds = append(cmds, m.components.search.Update(msg))

		if m.applySearchQuery() {
			m.logger.Debug("search query applied", "query", m.controls.query)
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
			m.logger.Debug("filters applied", "active_count", len(m.controls.filters))
			m.refreshTreeFromControls()
		}

		m.components.status.SetActiveFilterCount(len(m.controls.filters))
	}

	return tea.Batch(cmds...)
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
	if intent.Reset() {
		if len(m.controls.filters) == 0 {
			return false
		}

		m.logger.Debug("filter reset", "previously_active", len(m.controls.filters))
		clear(m.controls.filters)

		return true
	}

	if action, ok := intent.Toggle(); ok {
		if m.controls.filters[action] {
			m.logger.Debug("filter removed", "action", action)
			delete(m.controls.filters, action)
		} else {
			m.logger.Debug("filter added", "action", action)
			m.controls.filters[action] = true
		}

		return true
	}

	return false
}
