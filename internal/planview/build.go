// Package planview builds the action-grouped tree used by the Terraview UI.
package planview

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"github.com/marvinlanhenke/terraview/internal/terraform"
)

// actionIndex defines the stable ordering of top-level action groups in the UI.
var actionIndex = map[Action]int{
	ActionCreate:  0,
	ActionUpdate:  1,
	ActionDelete:  2,
	ActionReplace: 3,
	ActionNoOp:    4,
	ActionError:   5,
}

// FromTerraform converts a Terraform plan into an action-grouped node tree.
// Error diagnostics are included as ActionError nodes.
func FromTerraform(tfplan terraform.Plan) (*Node, error) {
	root := createRoot(tfplan)

	totalChanges, childCounter := 0, make(map[Action]int)

	if err := parseResourceChanges(root, tfplan.ResourceChanges, childCounter, &totalChanges); err != nil {
		return nil, err
	}

	if err := parseDiagnostics(root, tfplan.Diagnostics, childCounter, &totalChanges); err != nil {
		return nil, err
	}

	addLabelCount(root, childCounter, totalChanges)

	return root, nil
}

func createRoot(tfplan terraform.Plan) *Node {
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

	return root
}

func parseResourceChanges(root *Node, changes []terraform.ResourceChange, childCounter map[Action]int, total *int) error {
	for _, c := range changes {
		action, err := parseAction(c.Change.Actions)

		if err != nil {
			return err
		}

		idx, exists := actionIndex[action]

		if !exists {
			return errors.New("failed to lookup node group index by action type")
		}

		changes := compareChanges(c.Change.Before, c.Change.After)

		child := &Node{
			Id:      c.Address,
			Label:   c.Address,
			Kind:    NodeResource,
			Action:  action,
			changes: changes,
			Payload: c,
		}

		*total++
		childCounter[action]++

		root.Children[idx].Children = append(root.Children[idx].Children, child)
	}

	return nil
}

func parseDiagnostics(root *Node, diagnostics []terraform.Diagnostic, childCounter map[Action]int, total *int) error {
	for i, d := range diagnostics {
		if strings.ToLower(d.Severity) != "error" {
			continue
		}

		id := sha256.Sum256([]byte(d.Summary))
		summary := d.Summary
		action := ActionError

		idx, exists := actionIndex[action]

		if !exists {
			return errors.New("failed to lookup node group index by action type")
		}

		child := &Node{
			Id:      fmt.Sprintf("%s-%d", id, i),
			Label:   summary,
			Kind:    NodeResource,
			Action:  action,
			Payload: d,
		}

		*total++
		childCounter[action]++

		root.Children[idx].Children = append(root.Children[idx].Children, child)
	}

	return nil
}

func addLabelCount(root *Node, childCounter map[Action]int, totalChanges int) {
	for _, node := range root.Children {
		numerator, ok := childCounter[node.Action]

		if !ok {
			continue
		}

		denominator := totalChanges

		node.LabelCount = fmt.Sprintf("(%d/%d)", numerator, denominator)
	}
}
