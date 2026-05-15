package mapper

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/marvinlanhenke/terraview/internal/components/tree"
	"github.com/marvinlanhenke/terraview/internal/terraform"
)

var actionIndex = map[tree.Action]int{
	tree.ActionCreate:  0,
	tree.ActionUpdate:  1,
	tree.ActionDelete:  2,
	tree.ActionReplace: 3,
	tree.ActionNoOp:    4,
	tree.ActionError:   5,
}

func BuildTree(plan terraform.Plan) (*tree.Node, error) {
	root := &tree.Node{
		Id:       "root",
		Label:    fmt.Sprintf("Terraform plan %s", plan.TerraformVersion),
		Kind:     tree.NodeGroup,
		Action:   tree.ActionNoOp,
		Depth:    0,
		Expanded: true,
	}

	nodeGroups := make([]*tree.Node, len(actionIndex))
	for action, idx := range actionIndex {
		nodeGroups[idx] = &tree.Node{
			Id:       fmt.Sprintf("%s-node-group", string(action)),
			Label:    fmt.Sprintf("%s", strings.ToUpper(string(action[0]))+string(action[1:])),
			Kind:     tree.NodeGroup,
			Action:   action,
			Depth:    0,
			Expanded: false,
		}
	}

	root.Children = nodeGroups

	totalChanges := len(plan.ResourceChanges)

	childCounter := make(map[tree.Action]int)

	for _, rc := range plan.ResourceChanges {
		action, err := actionFromString(rc.Change.Actions)

		if err != nil {
			return nil, err
		}

		idx, exists := actionIndex[action]

		if !exists {
			return nil, errors.New("failed to lookup node group index by action type")
		}

		changes := compareChanges(rc.Change.Before, rc.Change.After)

		child := &tree.Node{
			Id:       rc.Address,
			Label:    rc.Address,
			Kind:     tree.NodeResource,
			Action:   action,
			Changes:  changes,
			Depth:    0,
			Expanded: false,
			Parent:   root.Children[idx],
			Payload:  rc,
		}

		childCounter[action]++

		root.Children[idx].Children = append(root.Children[idx].Children, child)
	}

	for _, node := range root.Children {
		numerator, exists := childCounter[node.Action]

		if !exists {
			continue
		}

		denominator := totalChanges

		node.LabelCount = fmt.Sprintf("(%d/%d)", numerator, denominator)
	}

	return root, nil
}

func compareChanges(before, after map[string]any) tree.Changes {
	result := tree.Changes{
		Before: map[string]any{},
		After:  map[string]any{},
	}

	for key, beforeVal := range before {
		afterVal, exists := after[key]

		// Removed field
		if !exists {
			result.Before[key] = beforeVal
			result.After[key] = nil
			continue
		}

		// Nested object
		beforeMap, beforeIsMap := beforeVal.(map[string]any)
		afterMap, afterIsMap := afterVal.(map[string]any)

		if beforeIsMap && afterIsMap {
			nested := compareChanges(beforeMap, afterMap)

			if len(nested.Before) > 0 || len(nested.After) > 0 {
				result.Before[key] = nested.Before
				result.After[key] = nested.After
			}

			continue
		}

		// Value changed
		if !reflect.DeepEqual(beforeVal, afterVal) {
			result.Before[key] = beforeVal
			result.After[key] = afterVal
		}
	}

	// Added fields
	for key, afterVal := range after {
		if _, exists := before[key]; !exists {
			result.Before[key] = nil
			result.After[key] = afterVal
		}
	}

	return result
}

func actionFromString(action []string) (tree.Action, error) {
	if len(action) <= 0 {
		return tree.ActionError, errors.New("failed to determine ActionType. No input actions provided.")
	}

	if len(action) >= 2 {
		switch {
		case string(tree.ActionCreate) == action[0] && string(tree.ActionDelete) == action[1]:
			return tree.ActionReplace, nil
		case string(tree.ActionDelete) == action[0] && string(tree.ActionCreate) == action[1]:
			return tree.ActionReplace, nil
		}
	}

	switch {
	case string(tree.ActionCreate) == action[0]:
		return tree.ActionCreate, nil
	case string(tree.ActionUpdate) == action[0]:
		return tree.ActionUpdate, nil
	case string(tree.ActionDelete) == action[0]:
		return tree.ActionDelete, nil
	case string(tree.ActionReplace) == action[0]:
		return tree.ActionReplace, nil
	case string(tree.ActionNoOp) == action[0]:
		return tree.ActionNoOp, nil
	default:
		return tree.ActionError, nil
	}
}
