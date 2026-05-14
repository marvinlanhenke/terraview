package mapper

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/marvinlanhenke/terraview/internal/components/tree"
	"github.com/marvinlanhenke/terraview/internal/terraform"
)

func BuildTree(plan terraform.Plan) (*tree.Node, error) {
	root := &tree.Node{
		Id:       "root",
		Label:    fmt.Sprintf("Terraform plan %s", plan.TerraformVersion),
		Action:   tree.ActionNoOp,
		Depth:    0,
		Expanded: true,
		Payload:  plan,
	}

	children := make([]*tree.Node, len(plan.ResourceChanges))

	// TODO: add children into buckets for each action and allocate under Root:1=>N:NodeGroup
	for i, rc := range plan.ResourceChanges {
		action, err := actionFromString(rc.Change.Actions)
		if err != nil {
			return nil, err
		}

		changes := compareChanges(rc.Change.Before, rc.Change.After)

		children[i] = &tree.Node{
			Id:       rc.Address,
			Label:    rc.Address,
			Action:   action,
			Changes:  changes,
			Depth:    0,
			Expanded: false,
			Payload:  rc,
		}
	}

	root.Children = children

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
