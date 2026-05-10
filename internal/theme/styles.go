package theme

import "charm.land/lipgloss/v2"

// TODO: Should we split this into multiple files theme/search.go, etc. and keep here only reusable color defs etc?
var (
	App = lipgloss.NewStyle().Padding(1, 2)

	Summary = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(0, 1)

	Pane = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(0, 1)

	Footer = lipgloss.NewStyle().Faint(true)

	// Colors
	BackgroundBlur  = lipgloss.Color("#303446")
	BackgroundFocus = lipgloss.Color("#5c5f77")

	TextBlur  = lipgloss.Color("#51576d")
	TextFocus = lipgloss.Color("#bcc0cc")

	AccentPrimary   = lipgloss.Color("#dd7878")
	AccentSecondary = lipgloss.Color("#dc8a78")
	AccentTertiary  = lipgloss.Color("#7287fd")

	// SearchBar
	SearchBar = lipgloss.NewStyle().Background(BackgroundBlur)

	SearchBarFocused = SearchBar.Background(BackgroundFocus)

	SearchNugget = lipgloss.NewStyle().
			Foreground(TextBlur).
			Background(AccentPrimary).
			Padding(0, 1).
			Bold(true)

	SearchInput = lipgloss.NewStyle().
			Foreground(TextBlur).
			Background(BackgroundBlur).
			Padding(0, 1)

	SearchInputFocused = SearchInput.Foreground(TextFocus)

	SearchStatus = lipgloss.NewStyle().
			Foreground(TextBlur).
			Background(AccentPrimary).
			Padding(0, 1)
)
