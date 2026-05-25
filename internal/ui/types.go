// Package ui defines shared UI types and theme primitives.
package ui

// Action identifies a Terraform plan action shown by UI components.
type Action string

const (
	// ActionCreate marks a resource that will be created.
	ActionCreate  Action = "create"
	// ActionUpdate marks a resource that will be updated in place.
	ActionUpdate  Action = "update"
	// ActionDelete marks a resource that will be deleted.
	ActionDelete  Action = "delete"
	// ActionReplace marks a resource that will be replaced.
	ActionReplace Action = "replace"
	// ActionNoOp marks a resource with no planned changes.
	ActionNoOp    Action = "no-op"
	// ActionError marks a resource with plan or diff errors.
	ActionError   Action = "error"
)

// ChangeSet contains before and after values for a resource diff.
type ChangeSet struct {
	// Before contains values from the prior resource state.
	Before map[string]any
	// After contains values from the planned resource state.
	After  map[string]any
}
