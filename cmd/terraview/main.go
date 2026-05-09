package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/marvinlanhenke/terraview/internal/app"
)

func main() {
	p := tea.NewProgram(app.New())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run terraview: %v\n", err)
		os.Exit(1)
	}
}
