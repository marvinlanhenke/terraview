---
title: "Quick Start"
linkTitle: "Quick Start"
weight: 10
description: "Build or install Terraview and inspect your first Terraform plan."
---

## 1. Clone the repository

```bash
git clone https://github.com/marvinlanhenke/terraview.git
cd terraview
```

## 2. Build or install Terraview

Build from source with the project `Makefile`:

```bash
make build
./bin/terraview -file plan.json
```

Or install the CLI directly into your Go bin directory:

```bash
go install github.com/marvinlanhenke/terraview/cmd/terraview@latest
```

## 3. Generate Terraform plan JSON

```bash
terraform plan -out=tfplan
terraform show -json tfplan > plan.json
```

If you already have a saved JSON plan file, you can skip this step.

## 4. Open the plan in Terraview

From a file:

```bash
./bin/terraview -file plan.json
```

Or pipe JSON directly from Terraform:

```bash
terraform show -json tfplan | terraview
```

Or force stdin explicitly:

```bash
./bin/terraview -file - < plan.json
```

## 5. Navigate the interface

- Use `j` / `k` or arrow keys to move through the tree
- Use `enter`, `l`, or `e` to expand a node
- Use `/` to search resources
- Use `f` to toggle action filters
- Use `ctrl+l` to focus the details pane and `p` to switch to raw plan view

## 6. Enable debug logging when needed

```bash
./bin/terraview -file plan.json -debug -log-file /tmp/terraview.log
```

Terraview writes debug logs to the configured file without polluting stdout or stderr during normal TUI use.
