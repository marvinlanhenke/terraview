package details

import (
	"sort"

	"github.com/marvinlanhenke/terraview/internal/planview"
)

type change struct {
	path   string
	before any
	after  any
}

func flattenChanges(n *planview.Node) []change {
	lines := make([]change, 0)

	flattenChangeMap("", n.Changes.Before, n.Changes.After, &lines)

	return lines
}

func flattenChangeMap(prefix string, before, after map[string]any, lines *[]change) {
	keys := make(map[string]struct{})
	for k := range before {
		keys[k] = struct{}{}
	}
	for k := range after {
		keys[k] = struct{}{}
	}

	sortedKeys := make([]string, 0, len(keys))
	for k := range keys {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}

		beforeVal := before[key]
		afterVal := after[key]

		beforeMap, beforeIsMap := beforeVal.(map[string]any)
		afterMap, afterIsMap := afterVal.(map[string]any)

		if beforeIsMap && afterIsMap {
			flattenChangeMap(path, beforeMap, afterMap, lines)
		}

		*lines = append(*lines, change{
			path:   path,
			before: beforeVal,
			after:  afterVal,
		})
	}
}
