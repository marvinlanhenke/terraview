package theme

import "charm.land/lipgloss/v2"

// TODO: Should we split this into multiple files theme/search.go, etc. and keep here only reusable color defs etc?
var (
	App = lipgloss.NewStyle().
		Padding(1, 2)

	Summary = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(0, 1)

	Pane = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(0, 1)

	Footer = lipgloss.NewStyle().
		Faint(true)

	// SearchBar
	SearchBarBG = lipgloss.Color("#353533")

	SearchBarFocusBG = lipgloss.Color("#242424")

	SearchBar = lipgloss.NewStyle().
			Background(SearchBarBG)

	SearchNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	SearchInput = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C1C6B2")).
			Background(SearchBarBG).
			Padding(0, 1)

	SearchStatus = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#6124DF")).
			Padding(0, 1)

	SearchInputFocused = SearchInput.
				Foreground(lipgloss.Color("#FFFDF5"))

	SearchBarFocused = SearchBar.
				Background(SearchBarFocusBG)
)
