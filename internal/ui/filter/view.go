package filter

import (
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/ui"
)

const (
	iconWidth  = 4
	labelWidth = 12
	countWidth = 8
	modalWidth = 28
)

// View renders the modal using active to mark enabled filters.
func (f *Modal) View(active map[ui.Action]bool) string {
	rows := make([]string, len(f.options))

	for i, option := range f.options {
		row := f.styles.row
		isCurrent := i == f.cursor

		icon := "[ ]"

		if active[option.Action] {
			icon = "[x]"
		}

		if isCurrent {
			row = f.styles.rowAlt
		}

		iconCol := row.Width(iconWidth).Render(icon)
		labelCol := row.Width(labelWidth).Render(option.Label)
		countCol := row.Width(countWidth).Align(lipgloss.Right).Render(option.Count)

		rows[i] = lipgloss.JoinHorizontal(lipgloss.Top, iconCol, labelCol, countCol)
	}

	list := lipgloss.JoinVertical(lipgloss.Left, rows...)

	header := f.styles.header.Width(modalWidth).Render("⚲ Filter:")

	content := lipgloss.JoinVertical(lipgloss.Left, header, "", list)

	return f.styles.modal.Width(modalWidth).Render(content)
}
