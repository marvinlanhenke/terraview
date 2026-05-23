package filter

import (
	"charm.land/lipgloss/v2"
)

const (
	iconWidth  = 4
	labelWidth = 12
	countWidth = 8
	modalWidth = 28
)

func (f *Modal) View(active map[Action]bool) string {
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

		label := option.Label
		count := option.Count

		iconCol := row.Width(iconWidth).Render(icon)
		labelCol := row.Width(labelWidth).Render(label)
		countCol := row.Width(countWidth).Align(lipgloss.Right).Render(count)

		rows[i] = lipgloss.JoinHorizontal(lipgloss.Top, iconCol, labelCol, countCol)
	}

	list := lipgloss.JoinVertical(lipgloss.Left, rows...)

	header := f.styles.header.Width(modalWidth).Render("⚲ Filter:")

	content := lipgloss.JoinVertical(lipgloss.Left, header, "", list)

	return f.styles.modal.Width(modalWidth).Render(content)
}
