package tree

type NodeKind int

const (
	NodeGroup NodeKind = iota
	NodeResource
	NodeOutput
	NodeError
)

type Node struct {
	Id       string
	Label    string
	Kind     NodeKind
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
