package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// View satisfies tea.Model by rendering the app layout and active overlays.
func (m Model) View() tea.View {
	content := m.renderAppContent()

	if m.focus == focusFilter {
		content = m.renderFilterOverlay(content)
	}

	return m.newView(content)
}

// renderAppContent renders the base app layout without modal overlays.
func (m Model) renderAppContent() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.components.search.View(),
		m.components.status.View(),
		" ",
		m.renderBody(),
		" ",
		m.renderFooter(),
	)
}

// renderBody renders the main tree/details pane split.
func (m Model) renderBody() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.components.tree.View(),
		" ",
		m.components.details.View(),
	)
}

// renderFooter renders general and focus-specific key help.
func (m Model) renderFooter() string {
	general := m.help.ShortHelpView(m.generalFooterBindings())
	specific := m.help.ShortHelpView(m.specificFooterBindings())

	bindings := []string{general}

	if specific != "" {
		bindings = append(bindings, specific)
	}

	return m.theme.Styles.Footer.
		Render(lipgloss.JoinVertical(lipgloss.Left, bindings...))
}

// renderFilterOverlay layers the filter modal above the base app content.
func (m Model) renderFilterOverlay(baseContent string) string {
	modal := m.components.filter.View(m.controls.filters)

	x := max(0, (m.size.width-lipgloss.Width(modal))/2)
	y := max(0, (m.size.height-lipgloss.Width(modal))/2)

	base := lipgloss.NewLayer(baseContent).Z(1)
	popup := lipgloss.NewLayer(modal).X(x).Y(y).Z(2)

	return lipgloss.NewCompositor(base, popup).Render()
}

// newView wraps rendered content in the Bubble Tea view settings used by the app.
func (m Model) newView(content string) tea.View {
	view := tea.NewView(m.theme.Styles.App.Render(content))
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion

	return view
}

// treePaneSize returns the tree pane dimensions for an app size.
func treePaneSize(width, height int) (int, int) {
	return max(20, width/3), max(5, height-13)
}

// detailsWidth returns the details pane width beside the tree pane.
func detailsWidth(appWidth, treeWidth int) int {
	return max(20, appWidth-treeWidth-defaultMargin-1)
}
