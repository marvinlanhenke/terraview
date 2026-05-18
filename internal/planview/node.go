package planview

import "maps"

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
	Changes    ChangeSet
	Children   []*Node
	Payload    any
}

func (n *Node) HasChildren() bool {
	return len(n.Children) > 0
}

func (n *Node) Diff() ChangeSet {
	if n == nil {
		return ChangeSet{}
	}

	return ChangeSet{
		Before: maps.Clone(n.Changes.Before),
		After:  maps.Clone(n.Changes.After),
	}
}

func (n *Node) IsResource() bool {
	return n != nil && n.Kind == NodeResource
}
