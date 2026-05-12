package tree

type nodeKind int

const (
	nodeGroup nodeKind = iota
	nodeResource
	nodeOutput
	nodeError
)

type node struct {
	id       string
	label    string
	kind     nodeKind
	action   action
	depth    int
	expanded bool

	parent   *node
	children []*node

	payload any
}

func (n *node) hasChildren() bool {
	return len(n.children) > 0
}
