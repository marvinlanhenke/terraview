package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m Model) View() tea.View {
	searchBar := m.components.search.View()
	status := m.components.status.View()
	tree := m.components.tree.View()
	details := m.components.details.View()

	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		tree,
		" ",
		details,
	)

	footer := m.theme.Styles.Footer.Render("/ search • esc back • q quit")

	appContent := lipgloss.JoinVertical(
		lipgloss.Left,
		searchBar,
		status,
		" ",
		body,
		" ",
		footer,
	)

	content := appContent

	if m.focus == FocusFilter {
		modal := m.components.filter.View(m.controls.filterView())

		x := max(0, (m.size.width-lipgloss.Width(modal))/2)
		y := max(0, (m.size.height-lipgloss.Height(modal))/2)

		base := lipgloss.NewLayer(appContent).Z(1)
		popup := lipgloss.NewLayer(modal).X(x).Y(y).Z(2)

		content = lipgloss.NewCompositor(base, popup).Render()
	}

	view := tea.NewView(m.theme.Styles.App.Render(content))
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion

	return view
}

func treePaneSize(width, height int) (int, int) {
	return max(20, width/3), max(5, height-10)
}

func detailsWidth(appWidth, treeWidth int) int {
	return max(20, appWidth-treeWidth-defaultMargin-1)
}
