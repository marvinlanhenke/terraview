package planview

type NodeKind int

const (
	NodeGroup NodeKind = iota
	NodeResource
)

type Node struct {
	Id         string
	Label      string
	LabelCount string
	Kind       NodeKind
	Action     Action
	Children   []*Node
	Payload    any
	changes    changeSet
}

func (n *Node) HasChildren() bool {
	return len(n.Children) > 0
}

func (n *Node) IsResource() bool {
	return n != nil && n.Kind == NodeResource
}

func (n *Node) ChangeSetBefore() map[string]any {
	return n.changes.changeSetBefore()
}

func (n *Node) ChangeSetAfter() map[string]any {
	return n.changes.changeSetAfter()
}
