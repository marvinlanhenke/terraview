package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

const defaultMargin = 4

func (m Model) View() tea.View {
	searchBar := m.search.View(max(0, m.width-defaultMargin), m.matchCount)

	summary := m.summary.View(max(0, m.width-defaultMargin))

	tree := m.tree.View(max(20, m.width/3), max(5, m.height-10))

	details := lipgloss.NewStyle().
		Width(max(20, m.width-m.width/3-defaultMargin)).
		Height(max(5, m.height-10)).
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
