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
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	planPath := flag.String("file", "", "path to terraform plan JSON file, or - for stdin")
	debug := flag.Bool("debug", false, "enable debug logging")
	logFile := flag.String("log-file", "debug.log", "path to log file")
	flag.Parse()

	logger, f, err := setupLogger(debug, logFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	if f != nil {
		defer f.Close()
	}

	logger.Debug("starting terraview", "file", *planPath, "debug", *debug, "log_file", *logFile)

	data, err := readPlanInput(*planPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	logger.Debug("plan input read", "bytes", len(data))

	plan, err := terraform.Parse(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse terraform plan\n: %v", err)
		os.Exit(1)
	}

	logger.Debug("terraform plan parsed", "terraform_version", plan.TerraformVersion, "resource_changes", len(plan.ResourceChanges), "diagnostics", len(plan.Diagnostics))

	root, err := planview.FromTerraform(plan, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build tree from plan: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(app.New(root, logger))

	logger.Debug("starting bubbletea program")

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

	rollingFile := &lumberjack.Logger{
		Filename:   *logFile,
		MaxSize:    10, // megabytes
		MaxBackups: 5,  // how many files to keep
		MaxAge:     30, // days
		Compress:   true,
	}

	return slog.New(slog.NewTextHandler(rollingFile, &slog.HandlerOptions{Level: slog.LevelDebug})), rollingFile, nil
}
