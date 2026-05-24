package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/marvinlanhenke/terraview/internal/app"
	"github.com/marvinlanhenke/terraview/internal/planview"
	"github.com/marvinlanhenke/terraview/internal/terraform"
)

func main() {
	planPath := flag.String("file", "", "path to terraform plan JSON file")
	flag.Parse()

	if *planPath == "" {
		fmt.Fprintln(os.Stderr, "usage: terraview -file <path-to-plan.json>")
		os.Exit(2)
	}

	data, err := os.ReadFile(*planPath)
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
