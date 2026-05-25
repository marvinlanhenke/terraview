# ◎─TERRAVIEW─◉

> Interactively explore Terraform plan outputs in the terminal.

Terraview is a terminal UI for navigating `terraform show -json` output. Instead of scrolling through walls of JSON, it gives you a structured, filterable, searchable tree of planned resource changes with a syntax-highlighted diff pane.

---

## Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Interface](#interface)
- [Key Bindings](#key-bindings)
- [Search](#search)
- [Filtering](#filtering)
- [Debug Logging](#debug-logging)

---

## Features

- **Action-grouped resource tree** — resources are organised into six fixed groups (Create, Update, Delete, Replace, No-op, Error) with per-group and total counts
- **Diff view** — inspect the before/after attribute changes for any selected resource, flattened to dot-path rows and syntax-highlighted with the Catppuccin Macchiato theme
- **Raw plan view** — toggle to see the full Terraform `ResourceChange` or `Diagnostic` JSON payload, pretty-printed and highlighted
- **Substring and regex search** — filter visible resources by plain text or a `/regex/` pattern; match counter updates live
- **Action filters** — a floating modal lets you show only the action types you care about; filter count is shown in the status bar
- **Expand / collapse** — navigate the tree with keyboard or collapse all at once; search automatically expands matching ancestors
- **Diagnostic errors** — plan errors are surfaced as first-class nodes in the Error group alongside resource changes
- **Stdin and file input** — pipe output directly from Terraform or point at a saved JSON file
- **Rolling debug log** — optional structured debug log with automatic rotation

---

## Installation

```bash
go install github.com/marvinlanhenke/terraview/cmd/terraview@latest
```

Or build from source:

```bash
git clone https://github.com/marvinlanhenke/terraview
cd terraview
go build ./cmd/terraview
```

---

## Usage

### From a saved plan file

```bash
terraform plan -out=tfplan
terraform show -json tfplan > plan.json
terraview -file plan.json
```

### Piped directly from Terraform

```bash
terraform show -json tfplan | terraview
```

### Explicit stdin

```bash
terraview -file - < plan.json
```

### All flags

| Flag               | Default     | Description                                                                                                                                 |
| ------------------ | ----------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| `-file <path>`     | —           | Path to a Terraform plan JSON file. Use `-` for stdin. If omitted, stdin is read automatically when it is a pipe.                           |
| `-debug`           | `false`     | Enable structured debug logging.                                                                                                            |
| `-log-file <path>` | `debug.log` | File to write debug logs to. Required when `-debug` is set. Logs rotate automatically (10 MB max, 5 backups, 30-day retention, compressed). |

---

## Interface

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ [S]  search resources...                    3 matches  ◎─TERRAVIEW─◉       │
│ ⚑ Plan: [+7] [~2] [-1] [*0] [=2] [!0]                    ⚲ Filter: 1       │
├──────────────────────────────┬──────────────────────────────────────────────┤
│ ⌘ Resources                 │ ▤ Details · Diff                            │
│                              │                                              │
│  ○ Create            (7/12)  │  attribute: tags.owner                       │
│  ● aws_ecs_cluster.main      │  (−) before:                                 │
│  ● aws_iam_role.app          │   null                                       │
│  ○ Update            (2/12)  │  (+) after:                                  │
│  ○ Delete            (1/12)  │   "platform"                                 │
│  ○ No-op             (2/12)  │                                              │
│                              │                                              │
├──────────────────────────────┴──────────────────────────────────────────────┤
│ q quit  esc back  / search  f filter  ctrl+h left pane  ctrl+l right pane   │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Search bar** — shows the active query, live match count, and the Terraview banner.

**Status bar** — plan summary with per-action counts (`+` create, `~` update, `-` delete, `*` replace, `=` no-op, `!` errors) and the active filter count.

**Tree pane** — left third of the terminal. Groups are collapsible; the selected node is highlighted. Groups with no resources are hidden.

**Details pane** — right two-thirds of the terminal. Shows the diff or raw plan for the selected resource. Scrollable with keyboard and mouse.

**Footer** — context-sensitive key binding hints that update as focus moves between panes.

---

## Key Bindings

### Global

| Key           | Action                                                                             |
| ------------- | ---------------------------------------------------------------------------------- |
| `q`, `ctrl+c` | Quit                                                                               |
| `/`           | Open search bar                                                                    |
| `enter`       | Apply search (when search is focused)                                              |
| `esc`         | Clear search and return to tree / close filter modal / return to tree from details |
| `f`           | Open / close filter modal                                                          |
| `ctrl+l`      | Move focus to details pane                                                         |
| `ctrl+h`      | Move focus to tree pane                                                            |

### Tree pane

| Key                    | Action                                                          |
| ---------------------- | --------------------------------------------------------------- |
| `j`, `↓`               | Move cursor down                                                |
| `k`, `↑`               | Move cursor up                                                  |
| `l`, `→`, `e`, `enter` | Expand selected node                                            |
| `h`, `←`, `c`, `enter` | Collapse selected node (or jump to parent if already collapsed) |
| `ctrl+e`               | Expand all nodes                                                |
| `ctrl+r`               | Collapse all nodes                                              |

### Details pane

| Key      | Action                                                    |
| -------- | --------------------------------------------------------- |
| `p`, `t` | Toggle between **Diff** view and **Plan** (raw JSON) view |

### Filter modal

| Key              | Action                             |
| ---------------- | ---------------------------------- |
| `j`, `↓`         | Move cursor down                   |
| `k`, `↑`         | Move cursor up                     |
| `enter`, `space` | Toggle highlighted filter on / off |
| `r`              | Reset all filters                  |
| `esc`            | Close modal                        |

---

## Search

Press `/` to open the search bar. The tree filters live as you type. Press `enter` to confirm and return focus to the tree, or `esc` to clear and cancel.

**Substring search** (default) — matches any resource whose address, label, action type, or JSON payload contains the query (case-insensitive):

```
s3
```

Matches `aws_s3_bucket.artifacts`, `aws_s3_bucket_policy.artifacts`, etc.

**Regex search** — wrap the pattern in forward slashes:

```
/aws_(s3|iam).*/
```

Matching is case-insensitive. If the pattern fails to compile it falls back silently to plain substring matching.

The match counter in the search bar (`N matches`) reflects the number of visible resource nodes that satisfy the current query and any active action filters.

When a search is active, all group nodes that contain matches are automatically expanded so results are always visible.

---

## Filtering

Press `f` to open the filter modal. Each row represents an action type with its resource count.

```
┌────────────────────────────┐
│ ⚲ Filter:                  │
│  [x] Create        (7/12)  │
│  [ ] Update        (2/12)  │
│  [ ] Delete        (1/12)  │
│  [ ] Replace       (0/12)  │
│  [ ] No-op         (2/12)  │
│  [ ] Error         (0/12)  │
└────────────────────────────┘
```

- `enter` or `space` — toggle a filter on or off
- `r` — reset all filters
- `f` or `esc` — close the modal

When filters are active the tree shows only the matching action groups. The status bar filter indicator changes from `⊙ Filter: 0` to `⚲ Filter: N` to reflect the number of active filters. Filters and search can be combined.

---

## Debug Logging

Enable structured debug logging to trace internal events (key presses, focus changes, search queries, filter toggles, tree updates):

```bash
terraview -file plan.json -debug -log-file /tmp/terraview.log
```

Logs are written in text format and rotated automatically. They are never written to stdout or stderr so the TUI is unaffected.
