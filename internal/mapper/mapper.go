package mapper

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/marvinlanhenke/terraview/internal/plan"
	"github.com/marvinlanhenke/terraview/internal/terraform"
)

var actionIndex = map[plan.Action]int{
	plan.ActionCreate:  0,
	plan.ActionUpdate:  1,
	plan.ActionDelete:  2,
	plan.ActionReplace: 3,
	plan.ActionNoOp:    4,
	plan.ActionError:   5,
}

func BuildTree(p terraform.Plan) (*plan.Node, error) {
	root := &plan.Node{
		Id:       "root",
		Label:    fmt.Sprintf("Terraform plan %s", p.TerraformVersion),
		Kind:     plan.NodeGroup,
		Action:   plan.ActionNoOp,
		Depth:    0,
		Expanded: true,
	}

	nodeGroups := make([]*plan.Node, len(actionIndex))
	for action, idx := range actionIndex {
		nodeGroups[idx] = &plan.Node{
			Id:       fmt.Sprintf("%s-node-group", string(action)),
			Label:    fmt.Sprintf("%s", strings.ToUpper(string(action[0]))+string(action[1:])),
			Kind:     plan.NodeGroup,
			Action:   action,
			Depth:    0,
			Expanded: false,
		}
	}

	root.Children = nodeGroups

	totalChanges := len(p.ResourceChanges)

	childCounter := make(map[plan.Action]int)

	for _, rc := range p.ResourceChanges {
		action, err := actionFromString(rc.Change.Actions)

		if err != nil {
			return nil, err
		}

		idx, exists := actionIndex[action]

		if !exists {
			return nil, errors.New("failed to lookup node group index by action type")
		}

		changes := compareChanges(rc.Change.Before, rc.Change.After)

		child := &plan.Node{
			Id:       rc.Address,
			Label:    rc.Address,
			Kind:     plan.NodeResource,
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

func compareChanges(before, after map[string]any) plan.Changes {
	result := plan.Changes{
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

func actionFromString(action []string) (plan.Action, error) {
	if len(action) <= 0 {
		return plan.ActionError, errors.New("failed to determine ActionType. No input actions provided.")
	}

	if len(action) >= 2 {
		switch {
		case string(plan.ActionCreate) == action[0] && string(plan.ActionDelete) == action[1]:
			return plan.ActionReplace, nil
		case string(plan.ActionDelete) == action[0] && string(plan.ActionCreate) == action[1]:
			return plan.ActionReplace, nil
		}
	}

	switch {
	case string(plan.ActionCreate) == action[0]:
		return plan.ActionCreate, nil
	case string(plan.ActionUpdate) == action[0]:
		return plan.ActionUpdate, nil
	case string(plan.ActionDelete) == action[0]:
		return plan.ActionDelete, nil
	case string(plan.ActionReplace) == action[0]:
		return plan.ActionReplace, nil
	case string(plan.ActionNoOp) == action[0]:
		return plan.ActionNoOp, nil
	default:
		return plan.ActionError, nil
	}
}
