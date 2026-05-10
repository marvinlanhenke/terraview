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
	treeCreate  = treeBase.Foreground(theme.AccentTertiary)
	treeUpdate  = treeBase.Foreground(theme.AccentTertiary)
	treeDelete  = treeBase.Foreground(theme.AccentTertiary)
	treeReplace = treeBase.Foreground(theme.AccentTertiary)
	treeNoOp    = treeBase.Foreground(theme.AccentTertiary)
	treeError   = treeBase.Foreground(theme.AccentTertiary)
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
