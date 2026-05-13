package tree

import (
	"fmt"
	"math/rand"
)

func GetNestedRoot(maxDepth, maxSiblings int) *Node {
	root := &Node{
		Id:     "root",
		Label:  "root",
		Kind:   NodeGroup,
		Action: ActionNoOp,
		Depth:  0,
	}

	actions := []Action{
		ActionCreate,
		ActionUpdate,
		ActionDelete,
		ActionReplace,
		ActionNoOp,
		ActionNoOp,
	}

	minDepth := 2
	minSiblings := 3

	depth := rand.Intn(max(maxDepth, minDepth)-minDepth+1) + minDepth
	siblings := rand.Intn(max(maxSiblings, minSiblings)-minSiblings+1) + minSiblings

	root.Children = getChildren(root, "child", 1, depth, siblings, actions)

	return root
}

func getChildren(parent *Node, label string, depth, maxDepth, maxSiblings int, actions []Action) []*Node {
	if depth > maxDepth {
		return nil
	}

	children := make([]*Node, maxSiblings)

	for i := range maxSiblings {
		id := fmt.Sprintf("%s-%d-%d", label, depth, i)
		action := actions[rand.Intn(len(actions))]

		child := &Node{
			Id:     id,
			Label:  id,
			Kind:   NodeResource,
			Action: action,
			Depth:  depth,
			Parent: parent,
		}

		child.Children = getChildren(child, label, depth+1, maxDepth, maxSiblings, actions)

		children[i] = child
	}

	return children
}
