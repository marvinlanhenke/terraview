package details

import (
	"sort"

	"github.com/marvinlanhenke/terraview/internal/ui"
)

type change struct {
	path   string
	before any
	after  any
}

func flattenChanges(ch ui.ChangeSet) []change {
	var rows []change

	flattenChangeMap("", normalizeMap(ch.Before), normalizeMap(ch.After), &rows)

	return rows
}

func flattenChangeMap(prefix string, before, after map[string]any, rows *[]change) {
	for _, key := range sortedUnionKeys(before, after) {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}

		appendValueChange(path, before[key], after[key], rows)
	}
}

func appendValueChange(path string, before, after any, rows *[]change) {
	beforeMap, beforeIsMap := asMap(before)
	afterMap, afterIsMap := asMap(after)

	if beforeIsMap || afterIsMap {
		start := len(*rows)

		flattenChangeMap(path, normalizeMap(beforeMap), normalizeMap(afterMap), rows)

		if len(*rows) == start {
			*rows = append(*rows, change{
				path:   path,
				before: before,
				after:  after,
			})
		}

		return
	}

	*rows = append(*rows, change{
		path:   path,
		before: before,
		after:  after,
	})
}

func asMap(v any) (map[string]any, bool) {
	m, ok := v.(map[string]any)
	return m, ok
}

func normalizeMap(m map[string]any) map[string]any {
	if m == nil {
		return map[string]any{}
	}

	return m
}

func sortedUnionKeys(before, after map[string]any) []string {
	keys := make(map[string]struct{}, len(before)+len(after))

	for k := range before {
		keys[k] = struct{}{}
	}

	for k := range after {
		keys[k] = struct{}{}
	}

	sorted := make([]string, 0, len(keys))
	for k := range keys {
		sorted = append(sorted, k)
	}

	sort.Strings(sorted)

	return sorted
}
