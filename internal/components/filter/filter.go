package filter

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/components/tree"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

type option struct {
	action tree.Action
	label  string
	count  string
}

type FilterModal struct {
	filters map[tree.Action]bool
	options []option

	cursor int
	width  int
	height int
	styles styles
}

func New(t theme.Theme) FilterModal {
	s := newStyles(t)
	f := make(map[tree.Action]bool)

	return FilterModal{
		filters: f,
		styles:  s,
	}
}

func (f *FilterModal) SetOptions(nodes []*tree.Node) {
	f.options = f.options[:0]

	seen := make(map[tree.Action]struct{})

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
	return nil
}

func (f FilterModal) View() string {
	rows := make([]string, len(f.options))

	for i, option := range f.options {
		row := lipgloss.NewStyle().
			Foreground(f.styles.palette.Text).
			Background(f.styles.palette.Surface)

		icon := "[ ]"

		if f.filters[option.action] {
			icon = "[x]"
		}

		label := strings.ToUpper(string(option.action[0])) + string(option.action[1:])
		count := option.count

		iconCol := row.Width(4).Render(icon)
		labelCol := row.Width(12).Render(label)
		countCol := row.Width(8).Align(lipgloss.Right).Render(count)

		rows[i] = lipgloss.JoinHorizontal(lipgloss.Top, iconCol, labelCol, countCol)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)

	return lipgloss.NewStyle().
		Width(28).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderBackground(f.styles.palette.Surface).
		Background(f.styles.palette.Surface).
		Render(content)
}

func (f *FilterModal) ToggleFilters(actions []tree.Action) {
	for _, action := range actions {
		f.ToggleSingleFilter(action)
	}
}

func (f *FilterModal) ToggleSingleFilter(action tree.Action) {
	before, exists := f.filters[action]

	if !exists {
		f.filters[action] = true
		return
	}

	f.filters[action] = !before
}

func (f *FilterModal) ResetFilters() {
	f.filters = nil
	f.filters = make(map[tree.Action]bool)
}

func (f FilterModal) GetFilters() map[tree.Action]bool {
	return f.filters
}
