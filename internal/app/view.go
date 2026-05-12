package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m Model) View() tea.View {
	searchBar := m.components.search.View()
	summary := m.components.summary.View()
	tree := m.components.tree.View()

	treeWidth, treeHeight := treePaneSize(m.size.width, m.size.height)

	// TODO
	details := lipgloss.NewStyle().
		Width(max(20, m.size.width-treeWidth-defaultMargin)).
		Height(treeHeight).
		Render("Details / Diff\n\nplaceholder")

	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		tree,
		" ",
		details,
	)

	// TODO
	footer := m.theme.Styles.Footer.Render("/ search • esc back • q quit")

	view := tea.NewView(
		m.theme.Styles.App.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				searchBar,
				summary,
				" ",
				body,
				" ",
				footer,
			),
		),
	)

	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion

	return view
}

func treePaneSize(width, height int) (int, int) {
	return max(20, width/3), max(5, height-10)
}
