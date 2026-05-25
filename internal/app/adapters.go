package app

import (
	"fmt"

	"github.com/marvinlanhenke/terraview/internal/planview"
	"github.com/marvinlanhenke/terraview/internal/ui"
	"github.com/marvinlanhenke/terraview/internal/ui/details"
	"github.com/marvinlanhenke/terraview/internal/ui/filter"
	"github.com/marvinlanhenke/terraview/internal/ui/status"
	"github.com/marvinlanhenke/terraview/internal/ui/tree"
)

// buildTreeNode adapts a planview node tree into the tree component model.
func buildTreeNode(n *planview.Node) *tree.Node {
	if n == nil {
		return nil
	}

	out := &tree.Node{
		Id:         n.Id,
		Label:      n.Label,
		LabelCount: n.LabelCount,
		Kind:       convertPlanNodeKind(n.Kind),
		Action:     convertPlanAction(n.Action),
		Payload:    n.Payload,
		Changes:    ui.ChangeSet{Before: n.ChangeSetBefore(), After: n.ChangeSetAfter()},
	}

	if len(n.Children) > 0 {
		out.Children = make([]*tree.Node, len(n.Children))
		for i, child := range n.Children {
			out.Children[i] = buildTreeNode(child)
		}
	}

	return out
}

// buildStats adapts the plan tree into status action counters.
func buildStats(n *planview.Node) status.Stats {
	if n == nil {
		return status.Stats{}
	}

	stats := status.Stats{}
	collectStats(n, &stats)

	return stats
}

// collectStats recursively counts resource nodes by action.
func collectStats(n *planview.Node, stats *status.Stats) {
	if n == nil {
		return
	}

	if n.Kind == planview.NodeResource {
		switch n.Action {
		case planview.ActionCreate:
			stats.Create++
		case planview.ActionUpdate:
			stats.Update++
		case planview.ActionDelete:
			stats.Delete++
		case planview.ActionReplace:
			stats.Replace++
		case planview.ActionNoOp:
			stats.NoOp++
		case planview.ActionError:
			stats.Errors++
		}
	}

	for _, child := range n.Children {
		collectStats(child, stats)
	}
}

// buildDetailsContent adapts the selected tree node into details pane content.
func buildDetailsContent(n *tree.Node) details.Content {
	if n == nil {
		return details.Content{Kind: details.KindNone}
	}

	content := details.Content{
		Key:   n.Id,
		Label: n.Label,
	}

	switch n.Kind {
	case tree.NodeGroup:
		content.Kind = details.KindGroup
	case tree.NodeResource:
		content.Kind = details.KindResource

		content.Changes = ui.ChangeSet{
			Before: n.Changes.Before,
			After:  n.Changes.After,
		}

		content.Payload = n.Payload

		content.IsError = n.IsError()
	default:
		content.Kind = details.KindNone
	}

	return content
}

// buildFilterOptions adapts top-level action groups into filter modal options.
func buildFilterOptions(groups []*planview.Node) []filter.Option {
	options := make([]filter.Option, 0, len(groups))

	for _, group := range groups {
		if group == nil {
			continue
		}

		options = append(options, filter.Option{
			Action: convertPlanAction(group.Action),
			Label:  group.Label,
			Count:  group.LabelCount,
		})
	}

	return options
}

// convertPlanAction maps a planview action to the matching UI action.
func convertPlanAction(a planview.Action) ui.Action {
	switch a {
	case planview.ActionCreate:
		return ui.ActionCreate
	case planview.ActionUpdate:
		return ui.ActionUpdate
	case planview.ActionDelete:
		return ui.ActionDelete
	case planview.ActionReplace:
		return ui.ActionReplace
	case planview.ActionNoOp:
		return ui.ActionNoOp
	case planview.ActionError:
		return ui.ActionError
	default:
		panic(fmt.Sprintf("unknown planview action %q", a))
	}
}

// convertPlanNodeKind maps a planview node kind to the matching tree node kind.
func convertPlanNodeKind(k planview.NodeKind) tree.NodeKind {
	switch k {
	case planview.NodeGroup:
		return tree.NodeGroup
	case planview.NodeResource:
		return tree.NodeResource
	default:
		panic(fmt.Sprintf("unknown planview node kind %q", k))
	}
}
