package planview

import "reflect"

type Changes struct {
	Before map[string]any
	After  map[string]any
}

func compareChanges(before, after map[string]any) Changes {
	result := Changes{
		Before: map[string]any{},
		After:  map[string]any{},
	}

	for key, beforeVal := range before {
		afterVal, exists := after[key]

		// Removed field
		if !exists {
			result.Before[key] = beforeVal
			result.After[key] = nil
			continue
		}

		// Nested object
		beforeMap, beforeIsMap := beforeVal.(map[string]any)
		afterMap, afterIsMap := afterVal.(map[string]any)

		if beforeIsMap && afterIsMap {
			nested := compareChanges(beforeMap, afterMap)

			if len(nested.Before) > 0 || len(nested.After) > 0 {
				result.Before[key] = nested.Before
				result.After[key] = nested.After
			}

			continue
		}

		// Value changed
		if !reflect.DeepEqual(beforeVal, afterVal) {
			result.Before[key] = beforeVal
			result.After[key] = afterVal
		}
	}

	// Added fields
	for key, afterVal := range after {
		if _, exists := before[key]; !exists {
			result.Before[key] = nil
			result.After[key] = afterVal
		}
	}

	return result
}
