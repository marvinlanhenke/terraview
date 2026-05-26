---
title: "Documentation"
linkTitle: "Documentation"
weight: 1
description: "Documentation for the Terraview repository."
---

Terraview is a terminal UI for exploring Terraform plan output in a structured way. It turns `terraform show -json` into an interactive tree with diff and raw JSON inspection, live search, and action filtering.

## Overview

| Topic | Summary |
| ----- | ------- |
| Input | Terraform plan JSON from a file or stdin |
| Interface | Tree pane, details pane, search bar, status bar, and filter modal |
| Output | Interactive terminal session for browsing resource changes |
| Build | `make build` or `go build ./cmd/terraview` |

## Sections

| Section | Description |
| ------- | ----------- |
| [Getting Started](getting-started/) | Install or build Terraview and run it against a plan file |
| [Usage](usage/) | Navigate the interface, search results, and filter large plans |
| [Reference](reference/) | CLI flags, input behavior, and code structure |

## Repository structure

```text
cmd/
`-- terraview/
    `-- main.go                # CLI entrypoint and input handling
internal/
|-- app/                      # Bubble Tea application composition
|-- planview/                 # Terraform change grouping and diff shaping
|-- terraform/                # Terraform plan JSON parsing
`-- ui/                       # Tree, details, filter, and status components
docs/
`-- hugo-pages/               # Hugo documentation site
testdata/                     # Sample Terraform plans used by tests
Makefile                      # Common fmt, lint, test, and build targets
README.md                     # Project overview and CLI usage
```

## Core capabilities

- Group resources into Create, Update, Delete, Replace, No-op, and Error sections
- Flatten nested attribute diffs into stable dot-path rows
- Toggle between diff view and full JSON plan view
- Search resource labels, actions, addresses, and payloads with substring or regex matching
- Filter visible resources by action type during navigation
