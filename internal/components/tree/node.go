package tree

import "errors"

type Action string

const (
	ActionCreate  Action = "create"
	ActionUpdate  Action = "update"
	ActionDelete  Action = "delete"
	ActionReplace Action = "replace"
	ActionNoOp    Action = "no-op"
	ActionError   Action = "error"
)

type Node struct {
	Id       string
	Label    string
	Action   Action
	Depth    int
	Expanded bool

	Parent   *Node
	Children []*Node

	Payload any
}

func (n *Node) hasChildren() bool {
	return len(n.Children) > 0
}

func GetAction(a []string) (Action, error) {
	if len(a) <= 0 {
		return ActionError, errors.New("failed to determine ActionType. No input actions provided.")
	}

	if len(a) >= 2 {
		switch {
		case string(ActionCreate) == a[0] && string(ActionDelete) == a[1]:
			return ActionReplace, nil
		case string(ActionDelete) == a[0] && string(ActionCreate) == a[1]:
			return ActionReplace, nil
		}
	}

	switch {
	case string(ActionCreate) == a[0]:
		return ActionCreate, nil
	case string(ActionUpdate) == a[0]:
		return ActionUpdate, nil
	case string(ActionDelete) == a[0]:
		return ActionDelete, nil
	case string(ActionReplace) == a[0]:
		return ActionReplace, nil
	case string(ActionNoOp) == a[0]:
		return ActionNoOp, nil
	default:
		return ActionError, nil
	}
}
