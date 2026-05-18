package planview

import (
	"errors"
	"fmt"
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

func FromTerraform(tfplan terraform.Plan) (*Node, error) {
	root := &Node{
		Id:     "root",
		Label:  fmt.Sprintf("Terraform plan %s", tfplan.TerraformVersion),
		Kind:   NodeGroup,
		Action: ActionNoOp,
	}

	nodeGroups := make([]*Node, len(actionIndex))

	for action, idx := range actionIndex {
		nodeGroups[idx] = &Node{
			Id:     fmt.Sprintf("%s-node-group", string(action)),
			Label:  fmt.Sprintf("%s", strings.ToUpper(string(action[0]))+string(action[1:])),
			Kind:   NodeGroup,
			Action: action,
		}
	}

	root.Children = nodeGroups

	totalChanges := len(tfplan.ResourceChanges)

	childCounter := make(map[Action]int)

	for _, rc := range tfplan.ResourceChanges {
		action, err := parseAction(rc.Change.Actions)

		if err != nil {
			return nil, err
		}

		idx, exists := actionIndex[action]

		if !exists {
			return nil, errors.New("failed to lookup node group index by action type")
		}

		changes := compareChanges(rc.Change.Before, rc.Change.After)

		child := &Node{
			Id:      rc.Address,
			Label:   rc.Address,
			Kind:    NodeResource,
			Action:  action,
			Changes: changes,
			Payload: rc,
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
