---
title: "Navigation and Interface"
linkTitle: "Navigation"
weight: 10
description: "Pane layout, key bindings, search, and filtering behavior in Terraview."
---

## Interface layout

Terraview is split into a few focused areas:

- Search bar: active query, match count, and project banner
- Status bar: counts for Create, Update, Delete, Replace, No-op, and Error actions
- Tree pane: grouped resource changes with expandable nodes
- Details pane: flattened diff rows or raw JSON plan view for the selected resource
- Filter modal: action-specific visibility toggles

## Global key bindings

| Key | Action |
| --- | ------ |
| `q`, `ctrl+c` | Quit |
| `/` | Open search bar |
| `enter` | Apply search when search is focused |
| `esc` | Clear search, close filter modal, or return to tree from details |
| `f` | Open or close the filter modal |
| `ctrl+l` | Move focus to details pane |
| `ctrl+h` | Move focus to tree pane |

## Tree pane bindings

| Key | Action |
| --- | ------ |
| `j`, `down` | Move cursor down |
| `k`, `up` | Move cursor up |
| `l`, `right`, `e`, `enter` | Expand selected node |
| `h`, `left`, `c`, `enter` | Collapse selected node or jump to parent |
| `ctrl+e` | Expand all nodes |
| `ctrl+r` | Collapse all nodes |

## Details pane bindings

| Key | Action |
| --- | ------ |
| `p`, `t` | Toggle between diff view and raw plan view |

## Search behavior

Search updates the visible tree as you type.

- Plain text uses case-insensitive substring matching
- Patterns wrapped in `/.../` are treated as regex
- Matching covers resource addresses, labels, action names, and searchable payload content
- Matching groups are expanded automatically so results remain visible

Example substring search:

```text
s3
```

Example regex search:

```text
/aws_(s3|iam).*/
```

## Filtering behavior

Press `f` to open the filter modal.

- `enter` or `space` toggles the highlighted action
- `r` resets all filters
- `esc` closes the modal

Filters combine with search. If both are active, the tree only shows resources that match the query and the enabled action types.
