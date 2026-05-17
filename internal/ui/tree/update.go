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
			selected := t.Selected()
			if selected != nil && selected.HasChildren() {
				selected.Expanded = !selected.Expanded
				t.rebuildVisible()
			}

		// Collapse
		case key.Matches(msg, keys.collapse):
			selected := t.Selected()
			if selected == nil {
				break
			}
			if selected.HasChildren() && selected.Expanded {
				selected.Expanded = false
				t.rebuildVisible()
			} else if selected.Parent != nil {
				for i, n := range t.visible {
					if n == selected.Parent {
						t.cursor = i
						break
					}
				}
			}
		}
	}

	t.clampCursor()
	t.syncViewport()

	return nil
}
