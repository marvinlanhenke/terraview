package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/marvinlanhenke/terraview/internal/app"
	"github.com/marvinlanhenke/terraview/internal/planview"
	"github.com/marvinlanhenke/terraview/internal/terraform"
)

func main() {
	planPath := flag.String("file", "", "path to terraform plan JSON file, or - for stdin")
	flag.Parse()

	data, err := readPlanInput(*planPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
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

func readPlanInput(planPath string) ([]byte, error) {
	switch {
	case planPath == "-":
		return io.ReadAll(os.Stdin)

	case planPath != "":
		return os.ReadFile(planPath)

	default:
		info, err := os.Stdin.Stat()
		if err != nil {
			return nil, fmt.Errorf("failed to inspect stdin: %w", err)
		}

		if info.Mode()&os.ModeCharDevice == 0 {
			return io.ReadAll(os.Stdin)
		}

		return nil, fmt.Errorf("usage: terraview -file <plan.json> or terraview -file -")
	}
}
