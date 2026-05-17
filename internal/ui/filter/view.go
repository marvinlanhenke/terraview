package filter

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func (f FilterModal) View() string {
	rows := make([]string, len(f.options))

	for i, option := range f.options {
		row := f.styles.row
		isCurrent := i == f.cursor

		icon := "[ ]"

		if f.filters[option.action] {
			icon = "[x]"
		}

		if isCurrent {
			row = f.styles.rowAlt
		}

		label := strings.ToUpper(string(option.action[0])) + string(option.action[1:])
		count := option.count

		iconCol := row.Width(4).Render(icon)
		labelCol := row.Width(12).Render(label)
		countCol := row.Width(8).Align(lipgloss.Right).Render(count)

		rows[i] = lipgloss.JoinHorizontal(lipgloss.Top, iconCol, labelCol, countCol)
	}

	list := lipgloss.JoinVertical(lipgloss.Left, rows...)

	header := f.styles.header.Width(28).Render("⚲ Filter:")

	content := lipgloss.JoinVertical(lipgloss.Left, header, "", list)

	return f.styles.modal.Width(28).Render(content)
}
