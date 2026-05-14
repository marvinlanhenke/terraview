package mapper

import (
	"errors"
	"fmt"

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

	for i, rc := range plan.ResourceChanges {
		action, err := actionFromString(rc.Change.Actions)
		if err != nil {
			return nil, err
		}

		children[i] = &tree.Node{
			Id:       rc.Address,
			Label:    rc.Address,
			Action:   action,
			Depth:    1,
			Expanded: false,
			Payload:  rc,
		}
	}

	root.Children = children

	return root, nil
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
