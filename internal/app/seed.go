package app

import "github.com/marvinlanhenke/terraview/internal/components/tree"

var (
	child1 = &tree.Node{
		ID:       "child-1",
		Label:    "child-1",
		Kind:     tree.NodeResource,
		Action:   tree.ActionCreate,
		Depth:    1,
		Expanded: false,
	}
	child2 = &tree.Node{
		ID:       "child-2",
		Label:    "child-2",
		Kind:     tree.NodeResource,
		Action:   tree.ActionDelete,
		Depth:    1,
		Expanded: false,
	}
	child2Sub1 = &tree.Node{
		ID:       "child-2-sub-1",
		Label:    "child-2-sub-1",
		Kind:     tree.NodeResource,
		Action:   tree.ActionUpdate,
		Depth:    2,
		Expanded: false,
	}
	root = &tree.Node{
		ID:       "example-root",
		Label:    "example-root",
		Kind:     tree.NodeGroup,
		Action:   tree.ActionNoOp,
		Depth:    0,
		Expanded: false,
	}
)

func getRoot() *tree.Node {

	child2.Children = []*tree.Node{child2Sub1}
	child2Sub1.Parent = child2
	root.Children = []*tree.Node{child1, child2}

	return root
}
