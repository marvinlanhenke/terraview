package planview

// NodeKind identifies whether a node groups children or represents a single entry.
type NodeKind int

const (
	// NodeGroup is a grouping node, such as the root or an action bucket.
	NodeGroup NodeKind = iota
	// NodeResource is a detail node for a resource change or diagnostic.
	NodeResource
)

// Node is a tree node rendered by the Terraview view.
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

// HasChildren reports whether n has child nodes.
func (n *Node) HasChildren() bool {
	return len(n.Children) > 0
}

// IsResource reports whether n is a detail node rather than a group node.
func (n *Node) IsResource() bool {
	return n != nil && n.Kind == NodeResource
}

// ChangeSetBefore returns a cloned map of changed attributes before the change.
func (n *Node) ChangeSetBefore() map[string]any {
	return n.changes.changeSetBefore()
}

// ChangeSetAfter returns a cloned map of changed attributes after the change.
func (n *Node) ChangeSetAfter() map[string]any {
	return n.changes.changeSetAfter()
}
