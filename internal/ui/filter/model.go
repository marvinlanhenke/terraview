// Package filter renders and updates the action filter modal.
package filter

import (
	"github.com/marvinlanhenke/terraview/internal/ui"
)

// Option describes one action filter row shown in the modal.
type Option struct {
	Action ui.Action
	Label  string
	Count  string
}

// Modal renders and updates the action filter modal.
type Modal struct {
	options []Option

	cursor int
	styles styles
}

// New returns an initialized Modal using t for styling.
func New(t ui.Theme) Modal {
	return Modal{
		styles: newStyles(t),
	}
}

// SetOptions replaces the filter options shown by the modal.
func (f *Modal) SetOptions(options []Option) {
	f.options = append([]Option(nil), options...)
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
