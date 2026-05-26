package details

import (
	"sort"

	"github.com/marvinlanhenke/terraview/internal/ui"
)

// change describes one flattened attribute change rendered in the details pane.
type change struct {
	path   string
	before any
	after  any
}

// flattenChanges converts a ChangeSet into sorted detail rows.
func flattenChanges(ch ui.ChangeSet) []change {
	var rows []change

	flattenChangeMap("", normalizeMap(ch.Before), normalizeMap(ch.After), &rows)

	return rows
}

// flattenChangeMap appends changes for every key in before or after.
func flattenChangeMap(prefix string, before, after map[string]any, rows *[]change) bool {
	keys := sortedUnionKeys(before, after)

	for _, key := range keys {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}

		appendValueChange(path, before[key], after[key], rows)
	}

	return len(keys) > 0
}

// appendValueChange appends a change row, recursively expanding map values first.
func appendValueChange(path string, before, after any, rows *[]change) {
	beforeMap, beforeIsMap := asMap(before)
	afterMap, afterIsMap := asMap(after)

	// map <-> map: should expand children and suppress the parent row
	// nil/map <-> map/nil: should expand children
	// map <-> scalar: should not expand; keep the parent row to preserve scalar side
	shouldExpand := (beforeIsMap && afterIsMap) ||
		(beforeIsMap && after == nil) ||
		(afterIsMap && before == nil)

	if shouldExpand {
		ok := flattenChangeMap(path, normalizeMap(beforeMap), normalizeMap(afterMap), rows)

		// If both sides are empty we add the row as a leaf
		if !ok {
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

// asMap returns v as a string-keyed map when possible.
func asMap(v any) (map[string]any, bool) {
	m, ok := v.(map[string]any)
	return m, ok
}

// normalizeMap returns an empty map when m is nil.
func normalizeMap(m map[string]any) map[string]any {
	if m == nil {
		return map[string]any{}
	}

	return m
}

// sortedUnionKeys returns the sorted set of keys present in before or after.
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
