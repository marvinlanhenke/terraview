package app

import (
	"fmt"
	"math/rand"

	"github.com/marvinlanhenke/terraview/internal/components/tree"
)

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

func getNestedRoot(maxDepth, maxSiblings int) *tree.Node {
	root := &tree.Node{
		ID:     "root",
		Label:  "root",
		Kind:   tree.NodeGroup,
		Action: tree.ActionNoOp,
		Depth:  0,
	}

	actions := []tree.Action{
		tree.ActionCreate,
		tree.ActionUpdate,
		tree.ActionDelete,
		tree.ActionReplace,
		tree.ActionNoOp,
		tree.ActionNoOp,
	}

	minDepth := 2
	minSiblings := 3

	depth := rand.Intn(max(maxDepth, minDepth)-minDepth+1) + minDepth
	siblings := rand.Intn(max(maxSiblings, minSiblings)-minSiblings+1) + minSiblings

	root.Children = getChildren(root, "child", 1, depth, siblings, actions)

	return root
}

func getChildren(parent *tree.Node, label string, depth, maxDepth, maxSiblings int, actions []tree.Action) []*tree.Node {
	if depth > maxDepth {
		return nil
	}

	children := make([]*tree.Node, maxSiblings)

	for i := range maxSiblings {
		id := fmt.Sprintf("%s-%d-%d", label, depth, i)
		action := actions[rand.Intn(len(actions))]

		child := &tree.Node{
			ID:     id,
			Label:  id,
			Kind:   tree.NodeResource,
			Action: action,
			Depth:  depth,
			Parent: parent,
		}

		child.Children = getChildren(child, label, depth+1, maxDepth, maxSiblings, actions)

		children[i] = child
	}

	return children
}
