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
)

type Changes struct {
	Before map[string]any
	After  map[string]any
}

type Node struct {
	Id         string
	Label      string
	LabelCount string
	Kind       NodeKind
	Action     Action
	Changes    Changes
	Depth      int
	Expanded   bool

	Parent   *Node
	Children []*Node

	Payload any
}

func (n *Node) hasChildren() bool {
	return len(n.Children) > 0
}
