package planview

import "reflect"

type changeSet struct {
	before map[string]any
	after  map[string]any
}

func (c changeSet) changeSetBefore() map[string]any { return cloneMap(c.before) }

func (c changeSet) changeSetAfter() map[string]any { return cloneMap(c.after) }

// compareChanges returns only the attributes that differ between before and after.
// Nested maps are compared recursively, and missing values are represented as nil.
func compareChanges(before, after map[string]any) changeSet {
	result := changeSet{
		before: map[string]any{},
		after:  map[string]any{},
	}

	for key, beforeVal := range before {
		afterVal, exists := after[key]

		// Removed field
		if !exists {
			result.before[key] = beforeVal
			result.after[key] = nil
			continue
		}

		// Nested object
		beforeMap, beforeIsMap := beforeVal.(map[string]any)
		afterMap, afterIsMap := afterVal.(map[string]any)

		if beforeIsMap && afterIsMap {
			nested := compareChanges(beforeMap, afterMap)

			if len(nested.before) > 0 || len(nested.after) > 0 {
				result.before[key] = nested.before
				result.after[key] = nested.after
			}

			continue
		}

		// Value changed
		if !reflect.DeepEqual(beforeVal, afterVal) {
			result.before[key] = beforeVal
			result.after[key] = afterVal
		}
	}

	// Added fields
	for key, afterVal := range after {
		if _, exists := before[key]; !exists {
			result.before[key] = nil
			result.after[key] = afterVal
		}
	}

	return result
}

func cloneValue(v any) any {
	switch t := v.(type) {
	case map[string]any:
		return cloneMap(t)
	case []any:
		return cloneSlice(t)
	default:
		return t
	}
}

func cloneMap(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}

	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = cloneValue(v)
	}

	return dst
}

func cloneSlice(src []any) []any {
	if src == nil {
		return nil
	}

	dst := make([]any, len(src))
	for i, v := range src {
		dst[i] = cloneValue(v)
	}

	return dst
}
