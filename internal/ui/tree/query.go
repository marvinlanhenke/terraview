package tree

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type matcher struct {
	raw     string
	literal string
	re      *regexp.Regexp
}

func newMatcher(query string) matcher {
	raw := strings.TrimSpace(query)

	m := matcher{
		raw:     raw,
		literal: strings.ToLower(raw),
	}

	if !isRegexQuery(raw) {
		return m
	}

	re, err := regexp.Compile("(?i)" + unwrapRegex(raw))
	if err == nil {
		m.re = re
	}

	return m
}

func (m matcher) Active() bool {
	return m.raw != ""
}

func (m matcher) MatchNode(n *Node) bool {
	if n == nil {
		return false
	}

	if !m.Active() {
		return true
	}

	return m.MatchString(n.Id) ||
		m.MatchString(n.Label) ||
		m.MatchString(string(n.Action)) ||
		m.MatchString(searchablePayload(n.Payload))
}

func (m matcher) MatchString(v string) bool {
	if !m.Active() {
		return true
	}

	if m.re != nil {
		return m.re.MatchString(v)
	}

	return strings.Contains(strings.ToLower(v), m.literal)
}

func isRegexQuery(query string) bool {
	if len(query) < 3 {
		// Require at least /x/, so "/" and "//" stay plain text.
		return false
	}

	if query[0] != '/' || query[len(query)-1] != '/' {
		return false
	}

	return strings.TrimSpace(unwrapRegex(query)) != ""
}

func unwrapRegex(query string) string {
	if len(query) < 3 {
		return ""
	}
	return query[1 : len(query)-1]
}

func searchablePayload(v any) string {
	if v == nil {
		return "null"
	}

	switch t := v.(type) {
	case string:
		return t
	default:
		b, err := json.Marshal(t)
		if err != nil {
			return fmt.Sprintf("%v", t)
		}
		return "\n" + string(b)
	}
}
