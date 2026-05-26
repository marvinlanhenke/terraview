---
title: "Architecture"
linkTitle: "Architecture"
weight: 20
description: "High-level package structure and data flow inside Terraview."
---

## Package layout

| Package | Responsibility |
| ------- | -------------- |
| `cmd/terraview` | CLI entrypoint, flag parsing, input reading, and logger setup |
| `internal/terraform` | Parse Terraform `show -json` output into typed Go structures |
| `internal/planview` | Compare before and after values, group resources by action, and build the tree model |
| `internal/app` | Compose the Bubble Tea program and adapt plan data into UI models |
| `internal/ui` | Shared UI types, theming, and pane-specific components |

## Runtime flow

1. `cmd/terraview` parses flags and reads plan JSON from a file or stdin.
2. `internal/terraform` parses the JSON payload into structured Terraform plan data.
3. `internal/planview` converts that plan into grouped nodes and flattened change sets.
4. `internal/app` adapts the plan nodes into UI-facing tree, details, filter, and status models.
5. `internal/ui` renders the Bubble Tea interface and handles focus, scrolling, search, and filtering.

## UI composition

The interface is intentionally split into small components:

- Tree: resource groups and resource selection
- Details: diff view and raw plan view
- Filter: action toggles with per-action counts
- Status: plan summary and filter state
- Search: live query input and visible match counting

## Design choices

- Nested map diffs are flattened into dot-path rows so related changes are easy to scan.
- Search operates across both labels and payload content to make large plans navigable.
- Raw plan inspection remains available so no detail is lost when the flattened diff omits surrounding context.
- Diagnostic errors are modeled as first-class tree nodes instead of being treated as a separate output mode.
