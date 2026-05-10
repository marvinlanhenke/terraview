package search

import (
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

var (
	searchBar = lipgloss.NewStyle().Background(theme.BackgroundBlur)

	searchBarFocused = searchBar.Background(theme.BackgroundFocus)

	searchNugget = lipgloss.NewStyle().
			Foreground(theme.TextBlur).
			Background(theme.AccentPrimary).
			Padding(0, 1).
			Bold(true)

	searchInput = lipgloss.NewStyle().
			Foreground(theme.TextBlur).
			Background(theme.BackgroundBlur).
			Padding(0, 1)

	searchInputFocused = searchInput.Foreground(theme.TextFocus)

	searchStatus = lipgloss.NewStyle().
			Foreground(theme.TextBlur).
			Background(theme.AccentPrimary).
			Padding(0, 1)
)
