package details

import "github.com/marvinlanhenke/terraview/internal/components/tree"

type summary struct {
	header string
}

func newSummary(n *tree.Node) *summary {
	_ = n
	header := "Changed Attributes"

	return &summary{
		header: header,
	}
}
