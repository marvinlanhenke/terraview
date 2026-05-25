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

func (m matcher) active() bool {
	return m.raw != ""
}

func (m matcher) matchNode(n *Node) bool {
	if n == nil {
		return false
	}

	if !m.active() {
		return true
	}

	return m.matchString(n.Id) ||
		m.matchString(n.Label) ||
		m.matchString(string(n.Action)) ||
		m.matchString(n.searchPayload)
}

func (m matcher) matchString(v string) bool {
	if !m.active() {
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

func prepareSearchPayloads(n *Node) {
	if n == nil {
		return
	}

	n.searchPayload = searchablePayload(n.Payload)

	for _, child := range n.Children {
		prepareSearchPayloads(child)
	}
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
