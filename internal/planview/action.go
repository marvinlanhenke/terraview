package planview

import (
	"errors"
	"fmt"
)

type Action string

const (
	ActionCreate  Action = "create"
	ActionUpdate  Action = "update"
	ActionDelete  Action = "delete"
	ActionReplace Action = "replace"
	ActionNoOp    Action = "no-op"
	ActionError   Action = "error"
)

var validActions = map[string]Action{
	"create":  ActionCreate,
	"update":  ActionUpdate,
	"delete":  ActionDelete,
	"replace": ActionReplace,
	"no-op":   ActionNoOp,
	"error":   ActionError,
}

func parseAction(actions []string) (Action, error) {
	if len(actions) == 0 {
		return ActionError, errors.New("failed to determine action: no input actions provided")
	}

	if len(actions) >= 2 {
		a, b := actions[0], actions[1]

		if (a == "create" && b == "delete") ||
			(a == "delete" && b == "create") {
			return ActionReplace, nil
		}
	}

	action, ok := validActions[actions[0]]
	if !ok {
		return ActionError, fmt.Errorf("unknown action: %q", actions[0])
	}

	return action, nil
}
