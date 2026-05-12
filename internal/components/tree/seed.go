package tree

import (
	"fmt"
	"math/rand"
)

func GetNestedRoot(maxDepth, maxSiblings int) *node {
	root := &node{
		id:     "root",
		label:  "root",
		kind:   nodeGroup,
		action: actionNoOp,
		depth:  0,
	}

	actions := []action{
		actionCreate,
		actionUpdate,
		actionDelete,
		actionReplace,
		actionNoOp,
		actionNoOp,
	}

	minDepth := 2
	minSiblings := 3

	depth := rand.Intn(max(maxDepth, minDepth)-minDepth+1) + minDepth
	siblings := rand.Intn(max(maxSiblings, minSiblings)-minSiblings+1) + minSiblings

	root.children = getChildren(root, "child", 1, depth, siblings, actions)

	return root
}

func getChildren(parent *node, label string, depth, maxDepth, maxSiblings int, actions []action) []*node {
	if depth > maxDepth {
		return nil
	}

	children := make([]*node, maxSiblings)

	for i := range maxSiblings {
		id := fmt.Sprintf("%s-%d-%d", label, depth, i)
		action := actions[rand.Intn(len(actions))]

		child := &node{
			id:     id,
			label:  id,
			kind:   nodeResource,
			action: action,
			depth:  depth,
			parent: parent,
		}

		child.children = getChildren(child, label, depth+1, maxDepth, maxSiblings, actions)

		children[i] = child
	}

	return children
}
