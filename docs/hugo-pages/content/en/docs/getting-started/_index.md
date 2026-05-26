---
title: "Getting Started"
linkTitle: "Getting Started"
weight: 10
description: "Prerequisites and first steps for working with Terraview."
---

This section covers the fastest path from cloning the repository or installing the CLI to opening your first Terraform plan in Terraview.

## Prerequisites

| Tool | Purpose | Version requirement |
| ---- | ------- | ------------------- |
| [Go](https://go.dev/doc/install) | Build Terraview from source | `1.25+` |
| [Terraform](https://developer.hashicorp.com/terraform/install) | Generate `terraform show -json` output | Any recent version that can emit JSON plans |
| Terminal emulator | Run the interactive TUI | ANSI-capable terminal recommended |

If you already have a saved plan JSON file, Terraform is not required to inspect it.

## Next steps

1. Follow the [Quick Start](quick-start/) to install or build Terraview.
2. Review the [Usage](../usage/) section to learn navigation and filtering.
3. Check the [CLI Reference](../reference/cli/) for flag details.
