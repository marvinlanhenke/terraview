package planview

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/marvinlanhenke/terraview/internal/terraform"
)

var actionIndex = map[Action]int{
	ActionCreate:  0,
	ActionUpdate:  1,
	ActionDelete:  2,
	ActionReplace: 3,
	ActionNoOp:    4,
	ActionError:   5,
}

func BuildTree(tfplan terraform.Plan) (*Node, error) {
	root := &Node{
		Id:       "root",
		Label:    fmt.Sprintf("Terraform plan %s", tfplan.TerraformVersion),
		Kind:     NodeGroup,
		Action:   ActionNoOp,
		Depth:    0,
		Expanded: true,
	}

	nodeGroups := make([]*Node, len(actionIndex))

	for action, idx := range actionIndex {
		nodeGroups[idx] = &Node{
			Id:       fmt.Sprintf("%s-node-group", string(action)),
			Label:    fmt.Sprintf("%s", strings.ToUpper(string(action[0]))+string(action[1:])),
			Kind:     NodeGroup,
			Action:   action,
			Depth:    0,
			Expanded: false,
		}
	}

	root.Children = nodeGroups

	totalChanges := len(tfplan.ResourceChanges)

	childCounter := make(map[Action]int)

	for _, rc := range tfplan.ResourceChanges {
		action, err := actionFromString(rc.Change.Actions)

		if err != nil {
			return nil, err
		}

		idx, exists := actionIndex[action]

		if !exists {
			return nil, errors.New("failed to lookup node group index by action type")
		}

		changes := compareChanges(rc.Change.Before, rc.Change.After)

		child := &Node{
			Id:       rc.Address,
			Label:    rc.Address,
			Kind:     NodeResource,
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

func compareChanges(before, after map[string]any) Changes {
	result := Changes{
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

func actionFromString(action []string) (Action, error) {
	if len(action) <= 0 {
		return ActionError, errors.New("failed to determine ActionType. No input actions provided.")
	}

	if len(action) >= 2 {
		switch {
		case string(ActionCreate) == action[0] && string(ActionDelete) == action[1]:
			return ActionReplace, nil
		case string(ActionDelete) == action[0] && string(ActionCreate) == action[1]:
			return ActionReplace, nil
		}
	}

	switch {
	case string(ActionCreate) == action[0]:
		return ActionCreate, nil
	case string(ActionUpdate) == action[0]:
		return ActionUpdate, nil
	case string(ActionDelete) == action[0]:
		return ActionDelete, nil
	case string(ActionReplace) == action[0]:
		return ActionReplace, nil
	case string(ActionNoOp) == action[0]:
		return ActionNoOp, nil
	default:
		return ActionError, nil
	}
}
