package filter

import (
	"github.com/marvinlanhenke/terraview/internal/ui"
	"github.com/marvinlanhenke/terraview/internal/ui/theme"
)

type Intent struct {
	Action    ui.Action
	HasToggle bool
	Reset     bool
}

type Option struct {
	Action ui.Action
	Label  string
	Count  string
}

type Modal struct {
	options []Option

	cursor int
	styles styles
}

func New(t theme.Theme) Modal {
	return Modal{
		styles: newStyles(t),
	}
}

func (f *Modal) SetOptions(options []Option) {
	f.options = options
	f.clampCursor()
}

func (f *Modal) selected() *Option {
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
