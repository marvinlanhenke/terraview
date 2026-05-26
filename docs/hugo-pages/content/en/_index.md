---
title: "Terraview"
linkTitle: "Home"
---

{{< blocks/cover title="◎─TERRAVIEW─◉" image_anchor="top" height="med" color="primary" >}}

<p class="lead mt-3">Interactive terminal UI for exploring <code>terraform show -json</code> output without digging through raw plan files.</p>
<div class="mx-auto mt-5">
  <a class="btn btn-lg btn-primary mr-3 mb-4" href="docs/getting-started/quick-start/">
    Quick Start <i class="fas fa-arrow-alt-circle-right ml-2"></i>
  </a>
  <a class="btn btn-lg btn-secondary mr-3 mb-4" href="docs/">
    Documentation <i class="fas fa-book ml-2"></i>
  </a>
</div>
{{< /blocks/cover >}}

{{< blocks/section color="white" type="row" >}}

{{% blocks/feature icon="fas fa-project-diagram" title="Tree-first plan navigation" %}}
Browse Terraform resource changes grouped by action with counts, expansion controls, and keyboard navigation.
{{% /blocks/feature %}}

{{% blocks/feature icon="fas fa-search" title="Search and filter" %}}
Combine case-insensitive substring or regex search with action filters to narrow large plans to the resources that matter.
{{% /blocks/feature %}}

{{% blocks/feature icon="fas fa-code" title="Diff and raw plan views" %}}
Inspect flattened attribute diffs or switch to the full JSON payload for a selected resource or diagnostic.
{{% /blocks/feature %}}

{{< /blocks/section >}}

{{< blocks/section color="light" type="row" >}}

{{% blocks/feature icon="fas fa-rocket" title="Quick Start" url="docs/getting-started/quick-start/" %}}
Install or build Terraview, generate a Terraform plan JSON file, and open it in the terminal UI.
{{% /blocks/feature %}}

{{% blocks/feature icon="fas fa-keyboard" title="Usage" url="docs/usage/navigation/" %}}
Learn the pane layout, key bindings, search behavior, and filtering workflow.
{{% /blocks/feature %}}

{{% blocks/feature icon="fas fa-cogs" title="Reference" url="docs/reference/architecture/" %}}
Review CLI flags, input behavior, and the package layout behind the application.
{{% /blocks/feature %}}

{{< /blocks/section >}}
