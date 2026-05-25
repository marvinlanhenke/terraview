package tree

import (
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/ui"
)

type actionStyle struct {
	marker string
	style  lipgloss.Style
}

type styles struct {
	palette       ui.Palette
	base          lipgloss.Style
	empty         lipgloss.Style
	background    lipgloss.Style
	backgroundAlt lipgloss.Style
	selected      lipgloss.Style
	create        lipgloss.Style
	update        lipgloss.Style
	delete        lipgloss.Style
	replace       lipgloss.Style
	noOp          lipgloss.Style
	error         lipgloss.Style
	label         lipgloss.Style
	labelAlt      lipgloss.Style
}

func newStyles(t ui.Theme) styles {
	p := t.Palette
	s := t.Styles

	base := lipgloss.NewStyle().Padding(0, 1).Background(p.Surface)

	return styles{
		palette:       p,
		base:          base,
		empty:         base.Faint(true),
		background:    s.Background,
		backgroundAlt: s.BackgroundAlt,
		selected:      base.Foreground(p.Text).Background(p.SurfaceAlt),

		create:  lipgloss.NewStyle().Foreground(p.Success),
		update:  lipgloss.NewStyle().Foreground(p.Warning),
		delete:  lipgloss.NewStyle().Foreground(p.Danger),
		replace: lipgloss.NewStyle().Foreground(p.Secondary),
		noOp:    lipgloss.NewStyle().Foreground(p.Text),
		error:   lipgloss.NewStyle().Foreground(p.Danger),

		label:    lipgloss.NewStyle().Foreground(p.Text).Background(p.Surface),
		labelAlt: lipgloss.NewStyle().Foreground(p.Text).Background(p.SurfaceAlt),
	}
}

func (s styles) actionMarker(a ui.Action) actionStyle {
	switch a {
	case ui.ActionCreate:
		return actionStyle{"+", s.create}
	case ui.ActionUpdate:
		return actionStyle{"~", s.update}
	case ui.ActionDelete:
		return actionStyle{"-", s.delete}
	case ui.ActionReplace:
		return actionStyle{"*", s.replace}
	case ui.ActionNoOp:
		return actionStyle{"=", s.noOp}
	case ui.ActionError:
		return actionStyle{"!", s.error}
	default:
		return actionStyle{" ", lipgloss.NewStyle()}
	}
}
