---
title: "CLI Reference"
linkTitle: "CLI"
weight: 10
description: "Flags, input modes, and common project commands for Terraview."
---

## Flags

| Flag | Default | Description |
| ---- | ------- | ----------- |
| `-file <path>` | unset | Path to a Terraform plan JSON file. Use `-` to force stdin. If omitted, Terraview reads stdin automatically when it is piped. |
| `-debug` | `false` | Enable structured debug logging. |
| `-log-file <path>` | `debug.log` | Path to the debug log file. Required when `-debug` is enabled. |

## Input modes

Terraview accepts plan data in three ways:

1. `-file plan.json` to read a saved JSON file
2. `-file -` to force stdin
3. No `-file` flag when stdin is already piped

If neither a file nor piped stdin is provided, Terraview exits with a usage error.

## Common commands

Build the binary:

```bash
make build
```

Run formatting, vetting, and tests:

```bash
make check
```

Build without the `Makefile`:

```bash
go build ./cmd/terraview
```

Run the CLI against a plan:

```bash
./bin/terraview -file plan.json
```

## Logging behavior

When `-debug` is disabled, Terraview discards log output. When it is enabled, logs are written through a rotating file writer with:

- 10 MB maximum file size
- 5 backup files
- 30 day retention
- compression enabled for rotated logs
