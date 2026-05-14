package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/marvinlanhenke/terraview/internal/components/tree"
	"github.com/marvinlanhenke/terraview/internal/terraform"
)

func main() {
	data, err := os.ReadFile("/home/mlanhenke/dev/projects/terraview/testdata/plans/mixed.json")
	if err != nil {
		panic(err)
	}

	plan, err := terraform.Parse(data)
	if err != nil {
		panic(err)
	}

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
		action, err := tree.GetAction(rc.Change.Actions)
		if err != nil {
			panic(err)
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

	data, err = json.MarshalIndent(root, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
