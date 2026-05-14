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
	NodeGroup int = iota
	NodeResource
)

type Changes struct {
	Before map[string]any
	After  map[string]any
}

// TODO:
// Add NodeKind to struct
// Use this to construct a root (static) with one NodeGroup per Action
// Example:
// root.Children = [{<Create>}, {<Update>}, ...]
// actual NodeResource Nodes will be allocated as childrens for each NodeGroup based on Action.
// If a NodeGroup has no children, it won't be rendered in the treeView.
type Node struct {
	Id       string
	Label    string
	Action   Action
	Changes  Changes
	Depth    int
	Expanded bool

	Parent   *Node
	Children []*Node

	Payload any
}

func (n *Node) hasChildren() bool {
	return len(n.Children) > 0
}
