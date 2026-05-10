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

type NodeKind int

const (
	NodeGroup NodeKind = iota
	NodeResource
	NodeOutput
	NodeError
)

type Node struct {
	ID       string
	Label    string
	Kind     NodeKind
	Action   Action
	Depth    int
	Expanded bool

	Parent   *Node
	Children []*Node

	Payload any
}

func (n *Node) HasChildren() bool {
	return len(n.Children) > 0
}
