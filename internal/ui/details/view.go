package details

import (
	"charm.land/lipgloss/v2"
)

func (d Details) View() string {
	if d.showEmptyState() {
		empty := d.styles.empty.
			Width(d.width).
			MaxWidth(d.width).
			Height(max(0, d.height-lipgloss.Height(d.header))).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(d.emptyStateMessage())

		return lipgloss.JoinVertical(lipgloss.Left, d.header, empty)

	}

	return lipgloss.JoinVertical(lipgloss.Left, d.header, d.viewport.View())
}
