package filter

import "github.com/marvinlanhenke/terraview/internal/theme"

type styles struct {
	palette theme.Palette
}

func newStyles(t theme.Theme) styles {
	p := t.Palette

	return styles{
		palette: p,
	}
}
