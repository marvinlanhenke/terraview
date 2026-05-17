package filter

import (
	"github.com/marvinlanhenke/terraview/internal/planview"
	"github.com/marvinlanhenke/terraview/internal/ui/theme"
)

type Intent struct {
	ToggleAction planview.Action
	HasToggle    bool
	Reset        bool
}

type option struct {
	action planview.Action
	label  string
	count  string
}

type Modal struct {
	options []option

	cursor int
	styles styles
}

func New(t theme.Theme) Modal {
	return Modal{
		styles: newStyles(t),
	}
}

func (f *Modal) SetOptions(nodes []*planview.Node) {
	f.options = f.options[:0]

	seen := make(map[planview.Action]struct{})

	for _, n := range nodes {
		if n == nil {
			continue
		}

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

	f.clampCursor()
}

func (f *Modal) selected() *option {
	if len(f.options) == 0 {
		return nil
	}

	return &f.options[f.cursor]
}

func (f *Modal) clampCursor() {
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
