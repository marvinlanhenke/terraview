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
	Changes    changeSet
	Children   []*Node
	Payload    any
}

func (n *Node) HasChildren() bool {
	return len(n.Children) > 0
}
