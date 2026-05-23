package theme

import (
	"charm.land/lipgloss/v2"
	"image/color"
)

type Palette struct {
	Surface         color.Color
	SurfaceAlt      color.Color
	SurfaceMuted    color.Color
	SurfaceEmbedded color.Color

	Text      color.Color
	TextMuted color.Color

	Primary   color.Color
	Secondary color.Color
	Info      color.Color
	Success   color.Color
	Warning   color.Color
	Danger    color.Color
}

type CommonStyles struct {
	App           lipgloss.Style
	Footer        lipgloss.Style
	Background    lipgloss.Style
	BackgroundAlt lipgloss.Style
}

type Theme struct {
	Palette Palette
	Styles  CommonStyles
}

func Default() Theme {
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
			Footer:        lipgloss.NewStyle().Faint(true),
			Background:    lipgloss.NewStyle().Background(p.Surface),
			BackgroundAlt: lipgloss.NewStyle().Background(p.SurfaceAlt),
		},
	}
}
