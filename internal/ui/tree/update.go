package tree

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func (t *Tree) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		// Up
		case key.Matches(msg, keys.up):
			t.cursor--

		// Down
		case key.Matches(msg, keys.down):
			t.cursor++

		// Expand
		case key.Matches(msg, keys.expand):
			selected, ok := t.selectedRow()
			if ok && selected.expandable {
				t.setExpanded(selected.node.Id, !selected.expanded)
				t.rebuildRows()
			}

		// Expand All
		case key.Matches(msg, keys.expandAll):
			if len(t.rows) > 0 {
				for _, r := range t.rows {
					if r.expandable {
						t.setExpanded(r.node.Id, true)
					}
				}
				t.rebuildRows()
			}

		// Collapse
		case key.Matches(msg, keys.collapse):
			selected, ok := t.selectedRow()
			if !ok {
				break
			}

			if selected.expandable && selected.expanded {
				t.setExpanded(selected.node.Id, false)
				t.rebuildRows()
			} else if selected.parent >= 0 {
				t.cursor = selected.parent
			}

		// Collapse All
		case key.Matches(msg, keys.collapseAll):
			if len(t.rows) > 0 {
				for _, r := range t.rows {
					if r.expandable {
						t.setExpanded(r.node.Id, false)
					}
				}
				t.rebuildRows()
			}

		}
	}

	t.clampCursor()
	t.syncViewport()

	return nil
}
