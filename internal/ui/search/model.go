package search

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/marvinlanhenke/terraview/internal/ui"
)

const placeholder = "search resources..."

type Search struct {
	input textinput.Model

	width   int
	matches int
	styles  styles
}

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

func (s *Search) SetWidth(width int) {
	s.width = max(0, width)
	s.syncInputWidth()
}

func (s *Search) SetMatches(matches int) {
	s.matches = max(0, matches)
	s.syncInputWidth()
}

func (s *Search) Focus() tea.Cmd {
	s.input.Placeholder = ""
	s.input.Focus()
	s.syncInputWidth()

	return textinput.Blink
}

func (s *Search) Focused() bool {
	return s.input.Focused()
}

func (s *Search) Blur() {
	s.input.Placeholder = placeholder
	s.input.Blur()
	s.syncInputWidth()
}

func (s *Search) Value() string {
	return s.input.Value()
}

func (s *Search) Clear() {
	s.matches = 0
	s.input.SetValue("")
	s.input.Placeholder = placeholder
}

func (s *Search) syncInputWidth() {
	s.input.SetWidth(s.layout().inputViewportWidth)
}
