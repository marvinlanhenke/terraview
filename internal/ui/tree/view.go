package tree

import (
	"charm.land/lipgloss/v2"
)

func (t *Tree) View() string {
	if len(t.rows) == 0 {
		empty := t.styles.empty.
			Width(t.width).
			MaxWidth(t.width).
			Height(t.height - lipgloss.Height(t.header)).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render("Nothing to show...")

		return lipgloss.JoinVertical(lipgloss.Left, t.header, empty)
	}

	return lipgloss.JoinVertical(lipgloss.Left, t.header, t.viewport.View())
}
