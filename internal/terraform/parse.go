package terraform

import "encoding/json"

func Parse(input []byte) (Plan, error) {
	var plan Plan

	err := json.Unmarshal(input, &plan)

	return plan, err
}
