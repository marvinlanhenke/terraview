package app

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

const defaultMargin = 4

func (m Model) View() tea.View {
	searchBar := theme.SearchBar.
		Width(max(0, m.width-defaultMargin)).
		Render(m.search.View())

	summary := theme.Summary.
		Width(max(0, m.width-defaultMargin)).
		Render(fmt.Sprintf("Summary: %q", m.summary))

	tree := theme.Pane.
		Width(max(20, m.width/3)).
		Height(max(5, m.height-10)).
		Render("Resource Tree\n\nplaceholder")

	details := theme.Pane.
		Width(max(20, m.width-m.width/3-defaultMargin)).
		Height(max(5, m.height-10)).
		Render("Details / Diff\n\nplaceholder")

	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		tree,
		details,
	)

	footer := theme.Footer.Render("/ search • esc back • q quit")

	view := tea.NewView(
		theme.App.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				searchBar,
				summary,
				body,
				footer,
			),
		),
	)

	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion

	return view
}
