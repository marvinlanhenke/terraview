package details

import (
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/planview"
)

func (d Details) View() string {
	if d.node == nil || d.node.Kind == planview.NodeGroup {
		empty := d.styles.empty.
			Width(d.width).
			MaxWidth(d.width).
			Height(d.height - lipgloss.Height(d.header)).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render("Nothing to show yet...")

		return lipgloss.JoinVertical(lipgloss.Left, d.header, empty)
	}

	return lipgloss.JoinVertical(lipgloss.Left, d.header, d.viewport.View())
}
