package tree

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
