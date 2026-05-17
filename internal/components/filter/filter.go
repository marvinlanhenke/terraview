package filter

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/plan"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

type option struct {
	action plan.Action
	label  string
	count  string
}

type FilterModal struct {
	filters map[plan.Action]bool
	options []option

	cursor int
	width  int
	height int
	styles styles
}

func New(t theme.Theme) FilterModal {
	s := newStyles(t)
	f := make(map[plan.Action]bool)

	return FilterModal{
		filters: f,
		styles:  s,
	}
}

func (f *FilterModal) SetOptions(nodes []*plan.Node) {
	f.options = f.options[:0]

	seen := make(map[plan.Action]struct{})

	for _, n := range nodes {
		if _, exists := seen[n.Action]; !exists {
			option := option{
				action: n.Action,
				label:  n.Label,
				count:  n.LabelCount,
			}
			f.options = append(f.options, option)
		}
		seen[n.Action] = struct{}{}
	}
}

func (f *FilterModal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		// Up
		case key.Matches(msg, keys.up):
			f.cursor--

		// Down
		case key.Matches(msg, keys.down):
			f.cursor++

		// Toggle Filter
		case key.Matches(msg, keys.toggle):
			selected := f.Selected()
			if selected != nil {
				f.ToggleSingleFilter(selected.action)
			}

		// Reset Filters
		case key.Matches(msg, keys.reset):
			f.ResetFilters()
		}
	}

	f.clampCursor()

	return nil
}

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

func (f *FilterModal) Selected() *option {
	if len(f.options) == 0 {
		return nil
	}

	return &f.options[f.cursor]
}

func (f *FilterModal) ToggleFilters(actions []plan.Action) {
	for _, action := range actions {
		f.ToggleSingleFilter(action)
	}
}

func (f *FilterModal) ToggleSingleFilter(action plan.Action) {
	before, exists := f.filters[action]

	if !exists {
		f.filters[action] = true
		return
	}

	f.filters[action] = !before
}

func (f *FilterModal) ResetFilters() {
	f.filters = nil
	f.filters = make(map[plan.Action]bool)
}

func (f FilterModal) GetFilters() map[plan.Action]bool {
	return f.filters
}

func (f *FilterModal) clampCursor() {
	if len(f.options) == 0 {
		f.cursor = 0
		return
	}

	if f.cursor < 0 {
		f.cursor = 0
	}

	if f.cursor >= len(f.options) {
		f.cursor = len(f.options) - 1
	}
}
