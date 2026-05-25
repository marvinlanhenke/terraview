package ui

import (
	"charm.land/lipgloss/v2"
	"image/color"
)

// Palette defines the colors shared by UI components.
type Palette struct {
	// Surface is the default component background color.
	Surface color.Color
	// SurfaceAlt is the highlighted component background color.
	SurfaceAlt color.Color
	// SurfaceMuted is the subdued component background color.
	SurfaceMuted color.Color
	// SurfaceEmbedded is the nested content background color.
	SurfaceEmbedded color.Color

	// Text is the default foreground color.
	Text color.Color
	// TextMuted is the subdued foreground color.
	TextMuted color.Color

	// Primary is the primary accent color.
	Primary color.Color
	// Secondary is the secondary accent color.
	Secondary color.Color
	// Info is the informational accent color.
	Info color.Color
	// Success is the success accent color.
	Success color.Color
	// Warning is the warning accent color.
	Warning color.Color
	// Danger is the destructive or error accent color.
	Danger color.Color
}

// CommonStyles contains base styles shared by UI components.
type CommonStyles struct {
	// App styles the top-level application container.
	App lipgloss.Style
	// Footer styles the app footer.
	Footer lipgloss.Style
	// Background styles default pane backgrounds.
	Background lipgloss.Style
	// BackgroundAlt styles highlighted pane backgrounds.
	BackgroundAlt lipgloss.Style
}

// Theme groups shared colors and styles for the UI.
type Theme struct {
	// Palette contains the theme colors.
	Palette Palette
	// Styles contains reusable common styles.
	Styles CommonStyles
}

// DefaultTheme returns the default Terraview UI theme.
func DefaultTheme() Theme {
	p := Palette{
		Surface:         lipgloss.Color("#303446"),
		SurfaceAlt:      lipgloss.Color("#3d4258"),
		SurfaceMuted:    lipgloss.Color("#292c3c"),
		SurfaceEmbedded: lipgloss.Color("#373b4f"),

		Text:      lipgloss.Color("#bcc0cc"),
		TextMuted: lipgloss.Color("#51576d"),

		Primary:   lipgloss.Color("#dd7878"),
		Secondary: lipgloss.Color("#dc8a78"),
		Info:      lipgloss.Color("#7287fd"),
		Success:   lipgloss.Color("#40a02b"),
		Warning:   lipgloss.Color("#df8e1d"),
		Danger:    lipgloss.Color("#d20f39"),
	}

	return Theme{
		Palette: p,
		Styles: CommonStyles{
			App:           lipgloss.NewStyle().Padding(1, 2),
			Footer:        lipgloss.NewStyle(),
			Background:    lipgloss.NewStyle().Background(p.Surface),
			BackgroundAlt: lipgloss.NewStyle().Background(p.SurfaceAlt),
		},
	}
}
