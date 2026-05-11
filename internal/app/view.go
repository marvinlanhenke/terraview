package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

func (m Model) View() tea.View {
	searchBar := m.search.View(max(0, m.width-defaultMargin), m.matchCount)

	summary := m.summary.View()

	tree := m.tree.View()
	treeWidth, treeHeight := treePaneSize(m.width, m.height)

	details := lipgloss.NewStyle().
		Width(max(20, m.width-treeWidth-defaultMargin)).
		Height(treeHeight).
		Render("Details / Diff\n\nplaceholder")

	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		tree,
		" ",
		details,
	)

	footer := theme.Footer.Render("/ search • esc back • q quit")

	view := tea.NewView(
		theme.App.Render(
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
