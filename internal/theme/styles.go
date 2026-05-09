package theme

import "charm.land/lipgloss/v2"

var (
	App = lipgloss.NewStyle().
		Padding(1, 2)

	SearchBar = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(0, 1)

	Summary = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(0, 1)

	Pane = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(0, 1)

	Footer = lipgloss.NewStyle().
		Faint(true)
)
