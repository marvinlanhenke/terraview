package tree

import (
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

var (
	treeBase = lipgloss.NewStyle().
			Padding(0, 1).
			Background(theme.BackgroundBlur)

	treeEmpty = treeBase.Faint(true)

	treeBackground = lipgloss.NewStyle().Background(theme.BackgroundBlur)

	treeSelected = treeBase.
			Foreground(theme.TextFocus).
			Background(theme.BackgroundFocus)

	// TODO: Color Coding
	treeCreate          = lipgloss.NewStyle().Foreground(theme.AccentTertiary)
	treeUpdate          = lipgloss.NewStyle().Foreground(theme.AccentTertiary)
	treeDelete          = lipgloss.NewStyle().Foreground(theme.AccentTertiary)
	treeReplace         = lipgloss.NewStyle().Foreground(theme.AccentTertiary)
	treeNoOp            = lipgloss.NewStyle().Foreground(theme.TextBlur)
	treeError           = lipgloss.NewStyle().Foreground(theme.AccentTertiary)
	treeLabel           = lipgloss.NewStyle().Foreground(theme.TextBlur).Background(theme.BackgroundBlur)
	treeLabelSelected   = lipgloss.NewStyle().Foreground(theme.TextFocus).Background(theme.BackgroundFocus)
	treeBackgroundFocus = theme.BackgroundFocus
)

func treeActionMarkerWithStyle(action Action) (string, lipgloss.Style) {
	switch action {
	case ActionCreate:
		return "+", treeCreate
	case ActionUpdate:
		return "~", treeUpdate
	case ActionDelete:
		return "-", treeDelete
	case ActionReplace:
		return "-/+", treeReplace
	case ActionNoOp:
		return "=", treeNoOp
	case ActionError:
		return "!", treeError
	default:
		return " ", lipgloss.NewStyle()
	}
}
