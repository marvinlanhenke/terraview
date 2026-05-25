package terraform

import "encoding/json"

// Parse unmarshals terraform show -json output into a Plan.
func Parse(input []byte) (Plan, error) {
	var plan Plan

	err := json.Unmarshal(input, &plan)

	return plan, err
}
