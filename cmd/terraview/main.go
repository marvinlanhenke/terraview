package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/marvinlanhenke/terraview/internal/app"
	"github.com/marvinlanhenke/terraview/internal/planview"
	"github.com/marvinlanhenke/terraview/internal/terraform"
)

const filepath = "/home/mlanhenke/dev/projects/terraview/testdata/plans"

func main() {
	data, err := os.ReadFile(filepath + "/mixed.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read terraform plan: %v\n", err)
		os.Exit(1)
	}

	plan, err := terraform.Parse(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse terraform plan: %v\n", err)
		os.Exit(1)
	}

	root, err := planview.FromTerraform(plan)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build tree from plan: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(app.New(root))

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run terraview: %v\n", err)
		os.Exit(1)
	}
}
