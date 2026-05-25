// Package search renders and updates the resource search bar.
package search

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/marvinlanhenke/terraview/internal/ui"
)

const placeholder = "search resources..."

// Search renders and updates the resource search bar.
type Search struct {
	input textinput.Model

	width   int
	matches int
	styles  styles
}

// New returns an initialized Search using t for styling.
func New(t ui.Theme) Search {
	s := newStyles(t)
	input := textinput.New()

	styles := input.Styles()

	styles.Focused.Placeholder = s.background
	styles.Focused.Text = s.inputAlt
	styles.Focused.Prompt = s.inputAlt

	styles.Blurred.Placeholder = s.backgroundMuted
	styles.Blurred.Text = s.input
	styles.Blurred.Prompt = s.input

	input.SetStyles(styles)

	input.Placeholder = placeholder
	input.CharLimit = 256
	input.Blur()

	return Search{
		input:  input,
		styles: s,
	}
}

// SetWidth sets the rendered search bar width.
func (s *Search) SetWidth(width int) {
	s.width = max(0, width)
	s.syncInputWidth()
}

// SetMatches records how many resources match the current search criteria.
func (s *Search) SetMatches(matches int) {
	s.matches = max(0, matches)
	s.syncInputWidth()
}

// Focus moves keyboard focus to the search input and starts cursor blinking.
func (s *Search) Focus() tea.Cmd {
	s.input.Placeholder = ""
	s.input.Focus()
	s.syncInputWidth()

	return textinput.Blink
}

// Focused reports whether the search input has keyboard focus.
func (s *Search) Focused() bool {
	return s.input.Focused()
}

// Blur removes keyboard focus from the search input.
func (s *Search) Blur() {
	s.input.Placeholder = placeholder
	s.input.Blur()
	s.syncInputWidth()
}

// Value returns the current search query.
func (s *Search) Value() string {
	return s.input.Value()
}

// Clear resets the search input, match count, and placeholder.
func (s *Search) Clear() {
	s.matches = 0
	s.input.SetValue("")
	s.input.Placeholder = placeholder
	s.syncInputWidth()
}

// syncInputWidth updates the underlying text input viewport width.
func (s *Search) syncInputWidth() {
	s.input.SetWidth(s.layout().inputViewportWidth)
}
