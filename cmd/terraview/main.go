package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/marvinlanhenke/terraview/internal/app"
	"github.com/marvinlanhenke/terraview/internal/planview"
	"github.com/marvinlanhenke/terraview/internal/terraform"
)

func main() {
	planPath := flag.String("file", "", "path to terraform plan JSON file, or - for stdin")
	debug := flag.Bool("debug", false, "enable debug logging")
	logFile := flag.String("log-file", "debug.log", "path to log file")
	flag.Parse()

	data, err := readPlanInput(*planPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	plan, err := terraform.Parse(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse terraform plan\n: %v", err)
		os.Exit(1)
	}

	root, err := planview.FromTerraform(plan)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build tree from plan: %v\n", err)
		os.Exit(1)
	}

	logger, f, err := setupLogger(debug, logFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	if f != nil {
		defer f.Close()
	}

	p := tea.NewProgram(app.New(root, logger))

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

func setupLogger(debug *bool, logFile *string) (*slog.Logger, io.Closer, error) {
	if *debug == false {
		return slog.New(slog.NewTextHandler(io.Discard, nil)), nil, nil
	}

	if *logFile == "" {
		return nil, nil, fmt.Errorf("-log-file is required when -debug is enabled")
	}

	f, err := os.OpenFile(*logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file %q: %w", *logFile, err)
	}

	return slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug})), f, nil
}
